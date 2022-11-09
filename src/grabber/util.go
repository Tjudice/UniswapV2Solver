package grabber

import (
	"UniswapV2Solver/src/data/arango"
	"context"
	"fmt"

	builder "gfx.cafe/open/arango"
)

func GetLastBlockForStage(ctx context.Context, db *arango.DB, stage int) (int, error) {
	return builder.NewBuilder[int](db.D()).
		Raw(fmt.Sprintf(LastStageBlockQuery, stage)).
		ReturnOne(ctx, "ret")
}

const LastStageBlockQuery = `
	for doc in StageProgress
		filter doc.stage == %d
		sort doc.block desc
		limit 1
		let ret = doc.block
`
