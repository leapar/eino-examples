# RAG 知识

## 文件内容切割

### 采用 Loader 进行文档加载

#### 功能描述

- **目标**：将多种格式的文档加载到系统中，为后续处理提供原始数据。
- **支持的文件类型**：
  - **PDF 文件**：通过解析 PDF 文档，提取文本内容。
  - **DOC/DOCX 文件**：解析 Microsoft Word 文档，提取文本和格式信息。
  - **Markdown 文件**：解析 Markdown 格式的文档，提取文本内容。
  - **HTML 文件**：解析 HTML 文件，提取网页文本内容。
  - **纯文本文件**：直接加载 `.txt` 文件内容。
  - **网络数据**：通过网络爬虫工具爬取网页内容，提取文本数据。

#### 示例工具

- **Python 库**：
  - **PDF 文件**：使用 `PyMuPDF` 或 `PyPDF2`。
  - **DOC/DOCX 文件**：使用 `python-docx`。
  - **HTML 文件**：使用 `BeautifulSoup`。
  - **网络爬虫**：使用 `Scrapy` 或 `requests` + `BeautifulSoup`。
- **代码示例**：

  ```go
    // 初始化 loader (以file loader为例)
    loader, err := file.NewFileLoader(ctx, &file.FileLoaderConfig{
        // 配置参数
        UseNameAsID: true,
    })
    if err != nil {
        log.Fatalf("file.NewFileLoader failed, err=%v", err)
    }

    // 加载文档
    filePath := "../../testdata/test.md"
    docs, err := loader.Load(ctx, document.Source{
        URI: filePath,
    })
    if err != nil {
        log.Fatalf("loader.Load failed, err=%v", err)
    }

  ```

### 文件分块

#### 功能描述

- **目标**：将加载的文档内容分割成多个小块（Chunk），以便后续处理。
- **分块策略**：
  - **按字符数分块**：根据预设的字符数（例如 512 个字符）将文本分割成多个块。
  - **按段落分块**：以段落为单位进行分块，每个段落作为一个独立的块。
  - **按语义分块**：通过自然语言处理技术，根据语义边界进行分块，确保每个块在语义上相对完整。
- **代码示例**：
  ```go
    func split() {
        f, err := os.Open("./testdata/test_pdf.pdf")

        p, err := NewPDFParser(ctx, nil)

        docs, err := p.Parse(ctx, f, WithToPages(true), parser.WithExtraMeta(map[string]any{"test": "test"}))
    }
  ```

## 文本块向量化

### 生成向量数据 Embedding

#### 功能描述

- **目标**：将文本块转换为向量表示，以便能够在向量空间中进行相似性计算。
- **实现方式**：通过调用外部的嵌入模型（Embedding Model）来完成，例如 OpenAI 的 Embedding API 或 Ollama 提供的模型。

### 参数

- **`model`**：指定嵌入模型的名称，例如 `nomic-embed-text:latest` 或 `all-minilm`。
- **`input`**：要生成向量数据的内容文本，可以是一个字符串或字符串列表。
- **高级参数**：
  - **`truncate`**：是否根据上下文长度裁剪内容。当输入文本长度超过模型的最大处理长度时，可以选择截断。
  - **`options`**：传入大模型本身的额外参数，例如 `temperature`（用于控制生成的随机性）。
  - **`keep_alive`**：默认超时时间，通常设置为 5 分钟。

### 示例

#### 请求（单个输入）

```shell
curl http://localhost:11434/api/embed -d '{
  "model": "all-minilm",
  "input": "Why is the sky blue?"
}'
```

#### 返回

```json
{
  "model": "all-minilm",
  "embeddings": [[0.010071029, -0.0017594862, 0.05007221, 0.04692972, 0.054916814, 0.008599704, 0.105441414, -0.025878139, 0.12958129, 0.031952348]],
  "total_duration": 14143917,
  "load_duration": 1019500,
  "prompt_eval_count": 8
}
```

#### 请求（多个输入）

```shell
curl http://localhost:11434/api/embed -d '{
  "model": "all-minilm",
  "input": ["Why is the sky blue?", "Why is the grass green?"]
}'
```

#### 返回

```json
{
  "model": "all-minilm",
  "embeddings": [
    [0.010071029, -0.0017594862, 0.05007221, 0.04692972, 0.054916814, 0.008599704, 0.105441414, -0.025878139, 0.12958129, 0.031952348],
    [-0.0098027075, 0.06042469, 0.025257962, -0.006364387, 0.07272725, 0.017194884, 0.09032035, -0.051705178, 0.09951512, 0.09072481]
  ]
}
```

## 向量存储 Indexing

### 功能描述

- **目标**：将生成的向量数据及其对应的内容存储到向量数据库中，以便后续的高效检索。
- **支持的向量数据库**：
  - **Milvus**：一个高性能的开源向量数据库，支持大规模向量检索。
  - **Volc VikingDB**：字节跳动提供的向量数据库解决方案。
  - **Redis**：通过 Redis 的向量扩展模块（如 RedisSearch）存储向量数据。
  - **LibSQL**：一个轻量级的 SQL 数据库，支持向量存储。
  - **Chromem-Go**：一个基于 Go 语言的向量数据库。
- **存储内容**：
  - 向量数据（Embedding）。
  - 原始文本内容。
  - 元数据（Metadata），例如文档的来源、创建时间等。

### 示例代码

```go
func (i *Indexer) Store(ctx context.Context, docs []*schema.Document, opts ...indexer.Option) (ids []string, err error) {
	defer func() {
		if err != nil {
			ctx = callbacks.OnError(ctx, err)
		}
	}()

	options := indexer.GetCommonOptions(&indexer.Options{
		Embedding: i.config.Embedding,
	}, opts...)

	ctx = callbacks.OnStart(ctx, &indexer.CallbackInput{Docs: docs})

	ids = make([]string, 0, len(docs))
	for _, sub := range chunk(docs, i.config.AddBatchSize) {
		documents, err := i.convertDocuments(ctx, sub, options)
		if err != nil {
			return nil, fmt.Errorf("convertDocuments failed: %w", err)
		}

		if err = i.collection.AddDocuments(ctx, documents, i.config.AddBatchSize); err != nil {
			return nil, fmt.Errorf("AddDocuments failed: %w", err)
		}

		ids = append(ids, iter(sub, func(t *schema.Document) string { return t.ID })...)
	}

	ctx = callbacks.OnEnd(ctx, &indexer.CallbackOutput{IDs: ids})

	return ids, nil
}
```

## 向量查询 Retriever

### 功能描述

- **目标**：当用户对大模型提问时，首先将用户输入的内容进行向量化，然后将用户输入的向量传入向量数据库进行相似性查询，向量数据库会返回一系列与用户输入最相似的记录。

### 请求参数

- **`embedding`**：用户输入内容的向量数据，为浮点数组。
- **`topK`**：最多返回的记录条数。
- **`scoreThreshold`**：记录相似程度的阈值，低于该阈值的记录将被过滤掉。

### 请求返回

- **`id`**：文档的唯一标识符。
- **`content`**：文档的内容。
- **`similarity`**：记录与用户输入的相似程度，通常是一个介于 0 到 1 之间的数值，越接近 1 表示越相似。

### 示例代码

```go
func (r *Retriever) Retrieve(ctx context.Context, query string, opts ...retriever.Option) (docs []*schema.Document, err error) {
	defer func() {
		if err != nil {
			ctx = callbacks.OnError(ctx, err)
		}
	}()

	options := retriever.GetCommonOptions(&retriever.Options{
		TopK:           &r.config.TopK,
		ScoreThreshold: &r.config.ScoreThreshold,
		Embedding:      r.config.Embedding,
	}, opts...)

	ctx = callbacks.OnStart(ctx, &retriever.CallbackInput{
		Query:          query,
		TopK:           r.config.TopK,
		ScoreThreshold: options.ScoreThreshold,
	})

	dense, err := r.customEmbedding(ctx, query, options)
	if err != nil {
		return nil, err
	}

	queryEmbedding := make([]float32, len(dense))
	for k, v := range dense {
		queryEmbedding[k] = float32(v)
	}

	result, err := r.collection.QueryEmbedding(ctx, queryEmbedding, int(math.Min(float64(r.collection.Count()), float64(*options.TopK))), nil, nil)
	if err != nil {
		return nil, err
	}

	docs = make([]*schema.Document, 0, len(result))
	for _, data := range result {
		if options.ScoreThreshold != nil && *options.ScoreThreshold > 0 && float64(data.Similarity) < *options.ScoreThreshold {
			continue
		}

		//fmt.Println(data.Similarity, data.Content)

		doc, err := r.data2Document(data)
		if err != nil {
			return nil, err
		}

		docs = append(docs, doc)
	}

	ctx = callbacks.OnEnd(ctx, &retriever.CallbackOutput{Docs: docs})

	return docs, nil
}
```

## 大模型回答

### 功能描述

- **目标**：根据向量查询返回的文档内容，将其作为上下文输入到大模型中，生成对用户问题的回答。
- **实现方式**：调用 Ollama 或 OpenAI 等大模型的接口，结合用户的问题和检索到的上下文内容，生成一个完整且准确的回答。

### 示例代码

```go
func createOllamaChatModel(ctx context.Context) model.ChatModel {
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434", // Ollama 服务地址
		Model:   "llama2",                 // 模型名称
	})
	if err != nil {
		log.Fatalf("create ollama chat model failed: %v", err)
	}
	return chatModel
}

func generate(ctx context.Context, llm model.ChatModel, in []*schema.Message) *schema.Message {
	result, err := llm.Generate(ctx, in)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}
	return result
}
```

### 注意事项

- **上下文长度限制**：大模型通常对上下文长度有限制（例如 2048 个 token），需要确保检索到的上下文内容不超过该限制。
- **多文档融合**：如果检索到多个相关文档，可以将它们的内容进行融合，生成一个更完整的上下文。
- **回答优化**：根据实际需求，可以对大模型的回答进行进一步优化，例如通过后处理去除无关内容或调整语言风格。

## 性能优化与注意事项

### 性能优化

- **向量检索优化**：
  - **索引类型选择**：根据数据特点选择合适的索引类型（如 IVF、HNSW 等），以提高检索效率。
  - **参数调优**：调整向量数据库的参数（如 `nprobe`、`efSearch` 等），以平衡检索速度和准确性。
- **向量化性能优化**：
  - **批量处理**：对多个文本块进行批量向量化，减少接口调用次数，提高效率。
  - **并行处理**：使用多线程或多进程对文本块进行并行向量化处理。
- **大模型调用优化**：
  - **缓存机制**：对常见问题的回答进行缓存，减少重复调用大模型的次数。
  - **上下文压缩**：对检索到的上下文内容进行压缩，去除冗余信息，以适应大模型的上下文长度限制。

### 注意事项

- **数据一致性**：确保向量数据库中的数据与原始文档内容保持一致，避免数据丢失或错误。
- **安全性**：在调用外部接口（如 OpenAI API）时，注意保护 API 密钥，避免泄露。
- **容错机制**：在系统中加入容错机制，例如在向量检索失败时提供备用方案，或在大模型调用失败时返回错误提示。
- **可扩展性**：设计系统时考虑可扩展性，方便后续增加新的功能或扩展数据量。

# 大模型调用 Tool 技术原理

Eino 等框架会把用户定义的 tool 函数进行封装转换，然后连同用户输入一起传入 Ollama 或 OpenAI 的请求接口中。根据返回内容中的 `tool_calls`，框架会解析返回内容的 `tool_calls`，然后传入参数调用对应的 tool。

## 大模型问答请求

```plaintext
POST /api/chat
```

使用提供的模型生成聊天中的下一条消息。通过设置 `"stream": false` 来禁用流式传输。最终的响应对象将包括请求的统计信息和其他附加数据。

### 请求参数

- **`model`**：大模型名称，指定用于生成回答的模型版本。
- **`messages`**：消息列表，包含用户和助手的对话历史。
- **`tools`**：如果模型支持工具调用，以 JSON 格式列出模型可以使用的工具列表。

大模型回答返回内容中可能包含 `tool_calls`：

- **`tool_calls`** (optional)：模型希望调用的工具列表，以 JSON 格式表示。

### 带 Tools 的请求示例

#### 请求

```shell
curl http://localhost:11434/api/chat -d '{
  "model": "llama3.2",
  "messages": [
    {
      "role": "user",
      "content": "What is the weather today in Paris?"
    }
  ],
  "stream": false,
  "tools": [
    {
      "type": "function",
      "function": {
        "name": "get_current_weather",
        "description": "Get the current weather for a location",
        "parameters": {
          "type": "object",
          "properties": {
            "location": {
              "type": "string",
              "description": "The location to get the weather for, e.g. San Francisco, CA"
            },
            "format": {
              "type": "string",
              "description": "The format to return the weather in, e.g. 'celsius' or 'fahrenheit'",
              "enum": ["celsius", "fahrenheit"]
            }
          },
          "required": ["location", "format"]
        }
      }
    }
  ]
}'
```

#### 返回

```json
{
  "model": "llama3.2",
  "created_at": "2024-07-22T20:33:28.123648Z",
  "message": {
    "role": "assistant",
    "content": "",
    "tool_calls": [
      {
        "function": {
          "name": "get_current_weather",
          "arguments": {
            "format": "celsius",
            "location": "Paris, FR"
          }
        }
      }
    ]
  },
  "done_reason": "stop",
  "done": true,
  "total_duration": 885095291,
  "load_duration": 3753500,
  "prompt_eval_count": 122,
  "prompt_eval_duration": 328493000,
  "eval_count": 33,
  "eval_duration": 552222000
}
```

### 技术原理详细说明

1. **工具封装与转换**：
   - 用户定义的 tool 函数被框架封装为可调用的接口。
   - 框架将这些工具转换为 JSON 格式，以便在请求中传输。

2. **请求构造**：
   - 用户输入和工具列表被组织成一个 JSON 请求体。
   - 请求体中包含模型名称、消息列表、工具列表等信息。

3. **模型处理**：
   - 大模型接收到请求后，根据工具定义和用户输入生成响应。
   - 如果模型需要调用工具，会在响应中包含 `tool_calls` 字段，指示需要调用的工具及其参数。

4. **工具调用**：
   - 框架解析响应中的 `tool_calls` 字段。
   - 根据 `tool_calls` 中的定义，框架调用相应的工具函数，并传入参数。

5. **最终响应**：
   - 工具函数的输出被整合到最终的回答中。
   - 最终的响应对象包含回答内容、统计信息和其他附加数据。

### 关键字段说明

- **`tool_calls`**：
  - **`function`**：表示工具的类型为函数调用。
  - **`name`**：工具函数的名称。
  - **`arguments`**：调用工具函数时需要传递的参数。

- **`messages`**：
  - **`role`**：消息的角色，例如 `"user"` 或 `"assistant"`。
  - **`content`**：消息的内容。

- **`tools`**：
  - **`type`**：工具的类型，例如 `"function"`。
  - **`function`**：工具的具体定义，包括名称、描述和参数。



### 下面是一个tool调用示例

```go
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

    // 初始化 tools
	todoTools := []tool.BaseTool{
		naviToTool, // 使用 InferTool 方式
		discToTool,
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
```



## 总结

RAG（Retrieval-Augmented Generation）技术通过结合向量检索和大语言模型的能力，能够高效地处理大规模文档数据，并生成准确且相关的回答。其核心流程包括文件加载、分块、向量化、向量存储、向量检索和大模型回答生成。通过合理配置和优化每个环节，可以显著提升系统的性能和用户体验。
