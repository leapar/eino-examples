package testgraph

import (
	"context"

	"github.com/cloudwego/eino/components/embedding"
)

type EmbeddingImpl struct {
	config *EmbeddingConfig
}

type EmbeddingConfig struct {
}

// newEmbedding component initialization function of node 'ollama' in graph 'demo'
func newEmbedding(ctx context.Context) (emb embedding.Embedder, err error) {
	// TODO Modify component configuration here.
	config := &EmbeddingConfig{}
	emb = &EmbeddingImpl{config: config}
	return emb, nil
}

func (impl *EmbeddingImpl) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float64, error) {
	panic("implement me")
}
