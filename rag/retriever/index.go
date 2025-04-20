package retriever

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func BuildEinoAgent(ctx context.Context, query string) (r compose.Runnable[*UserMessage, *schema.Message], err error) {
	const (
		InputToQuery   = "InputToQuery"
		ChatTemplate   = "ChatTemplate"
		ReactAgent     = "ReactAgent"
		Retriever      = "Retriever"
		InputToHistory = "InputToHistory"
	)

	g := compose.NewGraph[*UserMessage, *schema.Message]()
	_ = g.AddLambdaNode(InputToQuery, compose.InvokableLambdaWithOption(newLambda), compose.WithNodeName("UserMessageToQuery"))
	chatTemplateKeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}

	_ = g.AddChatTemplateNode(ChatTemplate, chatTemplateKeyOfChatTemplate)
	reactAgentKeyOfLambda, err := newLambda1(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLambdaNode(ReactAgent, reactAgentKeyOfLambda, compose.WithNodeName("ReAct Agent"))
	retriever, err := newRetriever(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(Retriever, retriever, compose.WithOutputKey("documents"))
	_ = g.AddLambdaNode(InputToHistory, compose.InvokableLambdaWithOption(newLambda2), compose.WithNodeName("UserMessageToVariables"))
	_ = g.AddEdge(compose.START, InputToQuery)
	_ = g.AddEdge(compose.START, InputToHistory)
	_ = g.AddEdge(ReactAgent, compose.END)
	_ = g.AddEdge(InputToQuery, Retriever)
	_ = g.AddEdge(Retriever, ChatTemplate)

	_ = g.AddEdge(InputToHistory, ChatTemplate)
	_ = g.AddEdge(ChatTemplate, ReactAgent)
	r, err = g.Compile(ctx, compose.WithGraphName("EinoAgent"), compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
