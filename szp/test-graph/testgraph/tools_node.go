package testgraph

import (
	"context"

	"github.com/cloudwego/eino-ext/components/tool/bingsearch"
	"github.com/cloudwego/eino/components/tool"
)

func newTool(ctx context.Context) (bt tool.BaseTool, err error) {
	// TODO Modify component configuration here.
	config := &bingsearch.Config{}
	bt, err = bingsearch.NewTool(ctx, config)
	if err != nil {
		return nil, err
	}
	return bt, nil
}
