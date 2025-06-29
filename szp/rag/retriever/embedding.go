package retriever

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/embedding/ollama"
	"github.com/cloudwego/eino-ext/components/embedding/openai"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/joho/godotenv"
)

func newEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	return newQWenEmbedding(ctx)
}

func newOllamaEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	config := &ollama.EmbeddingConfig{}
	eb, err = ollama.NewEmbedder(ctx, config)
	if err != nil {
		return nil, err
	}
	return eb, nil
}

func newQWenEmbedding(ctx context.Context) (eb embedding.Embedder, err error) {
	err = godotenv.Load()
	if err != nil {
		return nil, err
	}
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	config := &openai.EmbeddingConfig{
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		APIKey:  apiKey,
		Model:   "text-embedding-v3",
	}
	eb, err = openai.NewEmbedder(ctx, config)
	if err != nil {
		return nil, err
	}
	return eb, nil
}
