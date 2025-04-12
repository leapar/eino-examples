package retriever

import (
	"context"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/cloudwego/eino/components/embedding"
)

func newEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	// TODO Modify component configuration here.
	config := &ollama.EmbeddingConfig{}
	eb, err = ollama.NewEmbedder(ctx, config)
	if err != nil {
		return nil, err
	}
	return eb, nil
}
