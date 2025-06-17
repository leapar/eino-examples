/*
 * Copyright 2025 3dman.cn
 */

package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	payAgent := initAgent(NewPayGraph())
	hello := initAgent(NewHelloGraph())

	h, _ := newHost(ctx)
	hostMA, err := host.NewMultiAgent(ctx, &host.MultiAgentConfig{
		Host: *h,
		Specialists: []*host.Specialist{
			hello,
			payAgent,
		},
	})
	if err != nil {
		panic(err)
	}
	cb := &logCallback{}
	for { // 多轮对话，除非用户输入了 "exit"，否则一直循环
		println("\n\nYou: ") // 提示轮到用户输入了

		var message string
		scanner := bufio.NewScanner(os.Stdin) // 获取用户在命令行的输入
		for scanner.Scan() {
			message += scanner.Text()
			break
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}

		if message == "exit" {
			return
		}

		msg := &schema.Message{
			Role:    schema.User,
			Content: message,
		}

		out, err := hostMA.Stream(ctx, []*schema.Message{msg}, host.WithAgentCallbacks(cb))
		if err != nil {
			panic(err)
		}
		for {
			msg, err := out.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
			}
			fmt.Println(msg)
			println("Answer:", msg.Content)
		}

		out.Close()
	}

}

func newHost(ctx context.Context) (*host.Host, error) {
	chatModel := NewOpenaiLLM(ctx)
	return &host.Host{
		ToolCallingModel: chatModel,
		SystemPrompt:     "分析与的意图,进行合适的工具使用",
	}, nil
}

func initAgent(agent *host.Specialist, err error) *host.Specialist {
	if err != nil {
		panic(err)
	}
	return agent
}

type logCallback struct{}

func (l *logCallback) OnHandOff(ctx context.Context, info *host.HandOffInfo) context.Context {
	println("\nHandOff to", info.ToAgentName, "with argument", info.Argument)
	return ctx
}
