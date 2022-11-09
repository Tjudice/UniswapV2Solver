package data

import (
	"context"
	"fmt"

	"gfx.cafe/open/arango"
	"git.tuxpa.in/a/zlog/log"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-velocypack"
)

type BatchError struct {
	Errors     driver.ErrorSlice
	Collection string
}

func (b BatchError) Error() string {
	errStr := "BatchError for collection " + b.Collection + "\n"
	for i, e := range b.Errors {
		if e == nil {
			continue
		}
		errStr = errStr + fmt.Sprintf("Index %d: %s\n", i, e.Error())
	}
	return errStr
}

type BatchHandler[T arango.Document] struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	collection *arango.Collection[T]
	batch      []*packedDoc[T]
	batchSize  int
	numDocs    int
}

type packedDoc[T arango.Document] struct {
	V T
	K string
}

func (p *packedDoc[T]) MarshalVPack() (velocypack.Slice, error) {
	val, err := velocypack.Marshal(p.V)
	if err != nil {
		return nil, err
	}
	key, err := velocypack.Marshal(map[string]string{"_key": p.V.Key()})
	if err != nil {
		return nil, err
	}
	return velocypack.Merge(val, key)
}

func NewBatchHandler[T arango.Document](ctx context.Context, collection *arango.Collection[T], batchSize int) *BatchHandler[T] {
	ctxCancel, cancelFunc := context.WithCancel(ctx)
	return &BatchHandler[T]{
		ctx:        ctxCancel,
		cancelFunc: cancelFunc,
		collection: collection,
		batch:      make([]*packedDoc[T], batchSize),
		batchSize:  batchSize,
		numDocs:    0,
	}
}

func (b *BatchHandler[T]) Write(vals ...T) error {
	for _, val := range vals {
		b.batch[b.numDocs] = &packedDoc[T]{V: val, K: val.Key()}
		b.numDocs = b.numDocs + 1
		if b.numDocs >= b.batchSize {
			err := b.Flush()
			if err != nil {
				// log.Err(err).Msg("Error while writing batch to database")
				return err
			}
		}
	}
	return nil
}

func (b *BatchHandler[T]) Flush() error {
	if b.numDocs == 0 {
		return nil
	}
	docs := b.batch[:b.numDocs]
	_, eSlice, err := b.collection.C().CreateDocuments(b.ctx, docs)
	if err != nil {
		return err
	}
	b.numDocs = 0
	if eSlice.FirstNonNil() == nil {
		log.Debug().Int("add", len(docs)).Int("replace", 0).Str("collection", b.collection.Name()).Msg("BatchHandler")
		return nil
	}
	errDocs := make([]T, len(eSlice))
	errKeys := make([]string, len(eSlice))
	errs := 0
	for _, e := range eSlice {
		if e == nil {
			continue
		}
		errDocs[errs] = docs[errs].V
		errKeys[errs] = docs[errs].K
		errs = errs + 1
	}
	_, eSlice, err = b.collection.C().ReplaceDocuments(b.ctx, errKeys[:errs], errDocs[:errs])
	if err != nil {
		return err
	}
	if eSlice.FirstNonNil() != nil {
		return BatchError{Errors: eSlice, Collection: b.collection.Name()}
	}

	log.Debug().Int("add", len(docs)-errs).Int("replace", errs).Str("collection", b.collection.Name()).Msg("BatchHandler")

	return err
}

func (b *BatchHandler[T]) Close() error {
	err := b.Flush()
	b.batch = nil
	b.cancelFunc()
	if err != nil {
		log.Err(err).Str("collection", b.collection.Name()).Msg("err closing batch handler")
	}
	return err
}
