package main

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"

	"github.com/cloudwego/eino-examples/internal/logs"
)

var chatModel model.ChatModel
var bgCtx context.Context

func main() {
	ctx := context.Background()

	bgCtx = ctx

	naviToTool, err := utils.InferTool("navi_to", "带领游客去某个地方, eg: address...", NaviToFunc)
	if err != nil {
		logs.Errorf("InferTool failed, err=%v", err)
		return
	}

	discToTool, err := utils.InferTool("disc_to", "简单介绍, eg: artwork...", DiscToFunc)
	if err != nil {
		logs.Errorf("InferTool failed, err=%v", err)
		return
	}

	discMoreTool, err := utils.InferTool("disc_more", "详细介绍某个东西的具体情况, eg: artwork...", DiscMoreTool)
	if err != nil {
		logs.Errorf("InferTool failed, err=%v", err)
		return
	}

	tongjiToTool, err := utils.InferTool("tongji", "获取展厅统计数据, eg: sceneName...", TongjiToolFunc)
	if err != nil {
		logs.Errorf("InferTool failed, err=%v", err)
		return
	}

	// 初始化 tools
	todoTools := []tool.BaseTool{
		naviToTool, // 使用 InferTool 方式
		discToTool, discMoreTool,
		tongjiToTool,
	}

	chatModel, err = ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434", // Ollama 服务地址
		Model:   "qwen2.5:7b",             // 模型名称
	})
	if err != nil {
		logs.Errorf("NewChatModel failed, err=%v", err)
		return
	}

	// 获取工具信息, 用于绑定到 ChatModel
	toolInfos := make([]*schema.ToolInfo, 0, len(todoTools))
	var info *schema.ToolInfo
	for _, todoTool := range todoTools {
		info, err = todoTool.Info(ctx)
		if err != nil {
			logs.Infof("get ToolInfo failed, err=%v", err)
			return
		}
		toolInfos = append(toolInfos, info)
	}

	// 将 tools 绑定到 ChatModel
	err = chatModel.BindTools(toolInfos)
	if err != nil {
		logs.Errorf("BindTools failed, err=%v", err)
		return
	}

	// 创建 tools 节点
	todoToolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
		Tools: todoTools,
	})
	if err != nil {
		logs.Errorf("NewToolNode failed, err=%v", err)
		return
	}

	// 构建完整的处理链
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(chatModel, compose.WithNodeName("chat_model")).
		AppendToolsNode(todoToolsNode, compose.WithNodeName("tools"))

	// 编译并运行 chain
	agent, err := chain.Compile(ctx)
	if err != nil {
		logs.Errorf("chain.Compile failed, err=%v", err)
		return
	}

	// 运行示例
	resp, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.System,
			Content: "你是一个展馆导游，用中文回答所有问题",
		},
		{
			Role:    schema.User,
			Content: "带我去场馆东南门，简单介绍展品2，详细介绍鲁迅具体情况， 告诉我当前展厅访问统计数据", //"带我去 展品1"
		},
	})
	if err != nil {
		logs.Errorf("agent.Invoke failed, err=%v", err)
		return
	}

	// 输出结果
	for idx, msg := range resp {
		fmt.Println(msg.Name, msg.String())
		logs.Infof("消息序号: %d", idx)
		logs.Infof("消息产生: %s", msg.Role)
		logs.Infof("消息内容: %s", msg.Content)
	}

	resp = append(resp, &schema.Message{
		Role:    schema.System,
		Content: "你是一个展馆导游，用中文回答所有问题",
	})
	msg, err := chatModel.Generate(bgCtx, resp)

	if err != nil {
		return
	}

	fmt.Println("stream chunk: ", msg.String())

}

type NaviToParams struct {
	Address string `json:"address" jsonschema:"description=要去的地方"`
}

func NaviToFunc(_ context.Context, params *NaviToParams) (string, error) {
	logs.Infof("开始执行工具 navi_to 参数: %+v", params)

	return "启动导航去：" + params.Address, nil
}

type DiscToParams struct {
	Artwork string `json:"artwork" jsonschema:"description=要介绍的东西"`
}

func DiscToFunc(_ context.Context, params *DiscToParams) (string, error) {
	logs.Infof("开始执行工具 disc_to 参数: %+v", params)

	return "开始介绍：" + params.Artwork, nil
}

func DiscMoreTool(ctx context.Context, params *DiscToParams) (string, error) {
	logs.Infof("开始执行工具 disc_more 参数: %+v", params)

	msg, err := chatModel.Generate(ctx, []*schema.Message{
		schema.UserMessage("介绍一下" + params.Artwork + "?"),
	})

	if err != nil {
		return "", err
	}

	if len(msg.Content) == 0 {
		fmt.Println("stream chunk: ", msg.String())
	}

	return params.Artwork + "相关信息是:" + msg.Content, nil
}

type TongjiParams struct {
	SceneName string `json:"sceneName" jsonschema:"description=展厅名称"`
}

func TongjiToolFunc(_ context.Context, params *TongjiParams) (string, error) {
	logs.Infof("开始执行工具 tongji 参数: %+v", params)

	return "展厅[" + params.SceneName + "]访问数：" + "1001", nil
}
