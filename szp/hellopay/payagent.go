package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	"github.com/cloudwego/eino/schema"
)

func NewPayGraph() (*host.Specialist, error) {
	ctx := context.Background()
	chatModel := NewOpenaiLLM(ctx)
	graph := compose.NewGraph[[]*schema.Message, *schema.Message]()

	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(`
你负责支付相关问题的回答,当你回答不了就直接返回'不好意思我不知道这个问题'
当前的获取的知识信息:{rag_content}
`),
		schema.UserMessage("{query}"),
	)
	var err error

	//提示词节点
	if err = graph.AddChatTemplateNode("template", chatTpl); err != nil {
		return nil, err
	}

	//大模型节点
	if err = graph.AddChatModelNode("model", chatModel); err != nil {
		return nil, err
	}

	if err = graph.AddLambdaNode("search_rag", compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) (string, error) {
		fmt.Println("search_rag", input)
		return "查询的知识", nil

	}), compose.WithOutputKey("rag_content")); err != nil {
		return nil, err
	}

	if err = graph.AddLambdaNode("query_extractor", compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) (string, error) {
		return input[len(input)-1].Content, nil
	}), compose.WithOutputKey("query")); err != nil {
		return nil, err
	}

	//查询知识库
	graph.AddEdge("search_rag", "template")
	//整理请求信息 最后一句是用户的问题
	graph.AddEdge("query_extractor", "template")

	graph.AddEdge(compose.START, "search_rag")

	graph.AddEdge("template", "model")

	graph.AddEdge(compose.START, "query_extractor")

	graph.AddEdge("model", compose.END)

	r, err := graph.Compile(ctx,
		[]compose.GraphCompileOption{compose.WithGraphCompileCallbacks(&cb{}),
			compose.WithGraphName("top_level")}...)

	return &host.Specialist{
		AgentMeta: host.AgentMeta{
			Name:        "pay agent",
			IntendedUse: "当用户询问支付相关问题时使用",
		},
		Invokable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (*schema.Message, error) {
			return r.Invoke(ctx, input, agent.GetComposeOptions(opts...)...)
		},
		Streamable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (*schema.StreamReader[*schema.Message], error) {
			return r.Stream(ctx, input, agent.GetComposeOptions(opts...)...)
		},
	}, nil
}

type cb struct {
}

func (c *cb) OnFinish(ctx context.Context, info *compose.GraphInfo) {
	fmt.Println(info.Name)
	fmt.Println(info.Nodes)
}
