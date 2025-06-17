/*
 * Copyright 2025 3dman.cn
 */

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/cloudwego/eino-examples/internal/logs"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func main() {
	// 使用 SSE 初始化 MCP client
	ctx := context.Background()
	cli, _ := client.NewSSEMCPClient("http://localhost:21727/sse")
	cli.Start(ctx)
	defer cli.Close()

	// 发送 init request
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "current-time",
		Version: "1.0.0",
	}
	cli.Initialize(ctx, initRequest)

	// 查询 MCP Server 提供的 tools
	tools, _ := mcpp.GetTools(ctx, &mcpp.Config{Cli: cli})

	// for _, tool := range tools {
	// 	info, err := tool.Info(ctx)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(info.Name)
	// 	fmt.Println(info.Desc)
	// }

	// 将 MCP Tools 与 Eino 绑定
	// llm, _ := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
	// 	BaseURL: os.Getenv("OPENAI_API_URL"),
	// 	Model:   os.Getenv("MODEL_ID"),
	// 	APIKey:  os.Getenv("OPENAI_API_KEY"),
	// })

	config := &ollama.ChatModelConfig{
		BaseURL: "http://127.0.0.1:11434",
		Timeout: time.Duration(time.Second * 60),
		Model:   "qwen2.5:7b",
	}
	llm, err := ollama.NewChatModel(ctx, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	ragent, _ := react.NewAgent(ctx, &react.AgentConfig{
		Model:       llm,
		ToolsConfig: compose.ToolsNodeConfig{Tools: tools},
	})

	sr, err := ragent.Stream(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "所有展厅列表",
			//Content: "获取展厅“艺术展馆”的详细信息",
			//Content: "鲁迅详细介绍",
		},
	})
	if err != nil {
		logs.Errorf("failed to stream: %v", err)
		return
	}

	defer sr.Close() // remember to close the stream

	logs.Infof("\n\n===== start streaming =====\n\n")

	for {
		msg, err := sr.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// finish
				break
			}
			// error
			logs.Infof("failed to recv: %v", err)
			return
		}

		// 打字机打印
		logs.Tokenf("%v", msg.Content)
	}

	logs.Infof("\n\n===== finished =====\n")
}
