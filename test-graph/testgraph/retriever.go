package testgraph

import (
	"context"

	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
)

type RetrieverImpl struct {
	config *RetrieverConfig
}

type RetrieverConfig struct {
}

// newRetriever component initialization function of node 'chromem' in graph 'demo'
func newRetriever(ctx context.Context) (rtr retriever.Retriever, err error) {
	// TODO Modify component configuration here.
	config := &RetrieverConfig{}
	rtr = &RetrieverImpl{config: config}
	return rtr, nil
}

func (impl *RetrieverImpl) Retrieve(ctx context.Context, input string, opts ...retriever.Option) ([]*schema.Document, error) {
	panic("implement me")
}
