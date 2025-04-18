package testgraph

import (
	"context"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
)

// newChatModel component initialization function of node 'ChatModel1' in graph 'demo'
func newChatModel(ctx context.Context) (cm model.ChatModel, err error) {
	// TODO Modify component configuration here.
	config := &ollama.ChatModelConfig{
		BaseURL: "http://127.0.0.1:11434",
		Timeout: time.Duration(5000),
		Model:   "qwen2.5:7b"}
	cm, err = ollama.NewChatModel(ctx, config)
	if err != nil {
		return nil, err
	}
	return cm, nil
}
