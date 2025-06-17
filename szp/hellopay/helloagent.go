package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	"github.com/cloudwego/eino/schema"
)

func NewHelloGraph() (*host.Specialist, error) {
	ctx := context.Background()
	chatModel := NewOpenaiLLM(ctx)
	chain := compose.NewChain[[]*schema.Message, *schema.Message]()
	chain.AppendLambda(compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) ([]*schema.Message, error) {
		systemMsg := &schema.Message{
			Role:    schema.System,
			Content: "你负责与用户友好的交流",
		}
		return append([]*schema.Message{systemMsg}, input...), nil
	})).
		AppendChatModel(chatModel).
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input *schema.Message) (*schema.Message, error) {
			fmt.Println(`HelloGraph`, input.Content)
			return &schema.Message{
				Role:    schema.Assistant,
				Content: "你好啊: " + input.Content,
			}, nil
		}))

	r, err := chain.Compile(ctx)
	if err != nil {
		return nil, err
	}

	return &host.Specialist{
		AgentMeta: host.AgentMeta{
			Name:        "hello agent",
			IntendedUse: "负责与用户友好的交流,当用户的意图不清晰的时候使用",
		},
		Invokable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (*schema.Message, error) {
			return r.Invoke(ctx, input, agent.GetComposeOptions(opts...)...)
		},
		Streamable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (*schema.StreamReader[*schema.Message], error) {
			return r.Stream(ctx, input, agent.GetComposeOptions(opts...)...)
		},
	}, nil
}
