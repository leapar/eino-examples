package testgraph

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func Builddemo(ctx context.Context) (r compose.Runnable[any, any], err error) {
	const (
		ChatModel1                 = "ChatModel1"
		chromem                    = "chromem"
		ollama                     = "ollama"
		CustomDocumentTransformer5 = "CustomDocumentTransformer5"
		CustomChatTemplate6        = "CustomChatTemplate6"
		Lambda9                    = "Lambda9"
		Graph3                     = "Graph3"
	)
	g := compose.NewGraph[any, any]()
	chatModel1KeyOfChatModel, err := newChatModel(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatModelNode(ChatModel1, chatModel1KeyOfChatModel)
	chromemKeyOfRetriever, err := newRetriever(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddRetrieverNode(chromem, chromemKeyOfRetriever)
	ollamaKeyOfEmbedding, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddEmbeddingNode(ollama, ollamaKeyOfEmbedding)
	customDocumentTransformer5KeyOfDocumentTransformer, err := newDocumentTransformer(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddDocumentTransformerNode(CustomDocumentTransformer5, customDocumentTransformer5KeyOfDocumentTransformer)
	customChatTemplate6KeyOfChatTemplate, err := newChatTemplate(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddChatTemplateNode(CustomChatTemplate6, customChatTemplate6KeyOfChatTemplate)
	_ = g.AddLambdaNode(Lambda9, compose.InvokableLambda(newLambda))
	graph3KeyOftest, err := buildtest(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddGraphNode(Graph3, graph3KeyOftest,
		compose.WithGraphCompileOptions(
			compose.WithGraphName("test")))
	_ = g.AddEdge(compose.START, ChatModel1)
	_ = g.AddEdge(compose.START, ollama)
	_ = g.AddEdge(compose.START, Lambda9)
	_ = g.AddEdge(compose.START, Graph3)
	_ = g.AddEdge(ChatModel1, compose.END)
	_ = g.AddEdge(Graph3, compose.END)
	_ = g.AddEdge(CustomChatTemplate6, ChatModel1)
	_ = g.AddEdge(ollama, chromem)
	_ = g.AddEdge(chromem, Lambda9)
	_ = g.AddEdge(Lambda9, CustomDocumentTransformer5)
	_ = g.AddEdge(CustomDocumentTransformer5, CustomChatTemplate6)
	_ = g.AddEdge(CustomDocumentTransformer5, Graph3)
	r, err = g.Compile(ctx, compose.WithGraphName("demo"))
	if err != nil {
		return nil, err
	}
	return r, err
}

func buildtest(ctx context.Context) (ag compose.AnyGraph, err error) {
	g := compose.NewGraph[any, any]()
	return g, err
}
