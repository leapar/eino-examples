package testgraph

import (
	"context"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
)

type DocumentTransformerImpl struct {
	config *DocumentTransformerConfig
}

type DocumentTransformerConfig struct {
}

// newDocumentTransformer component initialization function of node 'CustomDocumentTransformer5' in graph 'demo'
func newDocumentTransformer(ctx context.Context) (tfr document.Transformer, err error) {
	// TODO Modify component configuration here.
	config := &DocumentTransformerConfig{}
	tfr = &DocumentTransformerImpl{config: config}
	return tfr, nil
}

func (impl *DocumentTransformerImpl) Transform(ctx context.Context, src []*schema.Document, opts ...document.TransformerOption) ([]*schema.Document, error) {
	panic("implement me")
}
