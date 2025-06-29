/*
 * Copyright 2025 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package indexing

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
