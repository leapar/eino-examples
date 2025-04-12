package testgraph

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func Builddemo(ctx context.Context) (r compose.Runnable[any, any], err error) {
	const ChatModel1 = "ChatModel1"
	g := compose.NewGraph[any, any]()
	chatModel1KeyOfChatModel, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatModelNode(ChatModel1, chatModel1KeyOfChatModel)
	_ = g.AddEdge(compose.START, ChatModel1)
	_ = g.AddEdge(ChatModel1, compose.END)
	r, err = g.Compile(ctx, compose.WithGraphName("demo"))
	if err != nil {
		return nil, err
	}
	return r, err
}
