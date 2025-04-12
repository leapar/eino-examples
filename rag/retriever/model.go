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

package retriever

import (
	"context"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
)

// newChatModel component initialization function of node 'ChatModel1' in graph 'demo'
func newChatModel(ctx context.Context) (cm model.ChatModel, err error) {
	return newSiliconflowChatModel(ctx)
}

func newOllamaChatModel(ctx context.Context) (cm model.ChatModel, err error) {
	// TODO Modify component configuration here.
	config := &ollama.ChatModelConfig{
		BaseURL: "http://127.0.0.1:11434",
		Timeout: time.Duration(5 * time.Minute),
		Model:   "qwen2.5:7b"}
	cm, err = ollama.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}

func newSiliconflowChatModel(ctx context.Context) (cm model.ChatModel, err error) {
	key := "sk-blyesrzaohzcmruwwxngupcpguywdrszdlpxcnctkenjcgqv"
	modelName := "Qwen/Qwen2.5-7B-Instruct"
	baseURL := "https://api.siliconflow.cn/v1"
	cm, err = openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
		APIKey:  key,
	})
	if err != nil {
		return nil, err
	}
	return cm, nil
}
