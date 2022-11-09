package arango

import (
	"UniswapV2Solver/src/data"
	"context"
	"crypto/tls"
	"os"
	"strings"

	"gfx.cafe/open/arango"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

func ArgsFromEnv() (uri string, auth string, dbName string) {
	uri = os.Getenv("ARANGO_ADDR")
	auth = os.Getenv("ARANGO_USER") + ":" + os.Getenv("ARANGO_PASS")
	dbName = os.Getenv("ARANGO_DATABASE")
	return
}

func ArgsFromEnvDefaults(db_name string) (uri string, auth string, dbName string) {
	uri, auth, dbName = ArgsFromEnv()
	if uri == "" {
		uri = "http://localhost:8529"
	}

	if auth == "" || auth == ":" {
		auth = "root:test"
	}
	if dbName == "" {
		if db_name == "" {
			dbName = "notify"
		} else {
			dbName = db_name
		}
	}
	return
}

type DB struct {
	d driver.Database

	PairCreatedEvent *arango.Collection[*data.PairCreatedEvent]
	MintEvent        *arango.Collection[*data.MintEvent]
	BurnEvent        *arango.Collection[*data.BurnEvent]
	SwapEvent        *arango.Collection[*data.SwapEvent]
	SyncEvent        *arango.Collection[*data.SyncEvent]

	// Stages

	StageProgress *arango.Collection[*data.StageProgress]
}

func (d *DB) D() driver.Database {
	return d.d
}

func NewDatabase(host string, auth string, dbName string) (*DB, error) {
	ctx := context.Background()
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{host},
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		ContentType: driver.ContentTypeVelocypack,
	})
	if err != nil {
		return nil, err
	}
	config := driver.ClientConfig{
		Connection: conn,
	}
	splt := strings.Split(auth, ":")
	if len(splt) > 1 {
		config.Authentication = driver.BasicAuthentication(splt[0], splt[1])
	}
	client, err := driver.NewClient(config)
	if err != nil {
		return nil, err
	}
	has, err := client.DatabaseExists(ctx, dbName)
	if err != nil {
		return nil, err
	}
	if !has {
		_, err = client.CreateDatabase(ctx, dbName, nil)
		if err != nil {
			return nil, err
		}
	}
	d, err := client.Database(ctx, dbName)
	if err != nil {
		return nil, err
	}

	o := &DB{d: d}

	o.PairCreatedEvent, err = arango.CollectionFromDocument(ctx, o.d, &data.PairCreatedEvent{})
	if err != nil {
		return nil, err
	}
	o.MintEvent, err = arango.CollectionFromDocument(ctx, o.d, &data.MintEvent{})
	if err != nil {
		return nil, err
	}
	o.BurnEvent, err = arango.CollectionFromDocument(ctx, o.d, &data.BurnEvent{})
	if err != nil {
		return nil, err
	}
	o.SwapEvent, err = arango.CollectionFromDocument(ctx, o.d, &data.SwapEvent{})
	if err != nil {
		return nil, err
	}
	o.SyncEvent, err = arango.CollectionFromDocument(ctx, o.d, &data.SyncEvent{})
	if err != nil {
		return nil, err
	}

	// Stages
	o.StageProgress, err = arango.CollectionFromDocument(ctx, o.d, &data.StageProgress{})
	if err != nil {
		return nil, err
	}

	return o, nil
}
