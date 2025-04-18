package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino-examples/rag/indexing"
	"github.com/cloudwego/eino-examples/rag/retriever"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()
	//index(ctx)
	rag(ctx)
}

func rag(ctx context.Context) {

	// Call RunAgent with the input
	sr, err := runAgent(ctx, "介绍 eino agent")
	if err != nil {
		fmt.Printf("Error from RunAgent: %v\n", err)
		return
	}

	// Print the response
	fmt.Print("🤖 : ")
	for {
		msg, err := sr.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error receiving message: %v\n", err)
			break
		}
		fmt.Print(msg.Content)
	}
	fmt.Println()
	fmt.Println()

}

func index(ctx context.Context) {
	err := indexMarkdownFiles(ctx, "./eino-docs")
	if err != nil {
		panic(err)
	}

	fmt.Println("index success")
}

func indexMarkdownFiles(ctx context.Context, dir string) error {
	runner, err := indexing.BuildKnowledgeIndexing(ctx)
	if err != nil {
		return fmt.Errorf("build index graph failed: %w", err)
	}

	// 遍历 dir 下的所有 markdown 文件
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk dir failed: %w", err)
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".md") {
			fmt.Printf("[skip] not a markdown file: %s\n", path)
			return nil
		}

		fmt.Printf("[start] indexing file: %s\n", path)

		ids, err := runner.Invoke(ctx, document.Source{URI: path})
		if err != nil {
			return fmt.Errorf("invoke index graph failed: %w", err)
		}

		fmt.Printf("[done] indexing file: %s, len of parts: %d\n", path, len(ids))

		return nil
	})

	return err
}

func runAgent(ctx context.Context, query string) (*schema.StreamReader[*schema.Message], error) {

	runner, err := retriever.BuildEinoAgent(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to build agent graph: %w", err)
	}

	fmt.Println("问：", query)

	userMessage := &retriever.UserMessage{
		Query: query,
		//	History: make([]*schema.Message, 0),
	}

	sr, err := runner.Stream(ctx, userMessage)
	if err != nil {
		return nil, fmt.Errorf("failed to stream: %w", err)
	}

	srs := sr.Copy(2)

	go func() {
		// for save to memory
		fullMsgs := make([]*schema.Message, 0)

		defer func() {
			// close stream if you used it
			srs[1].Close()
			fmt.Println(query)
			fullMsg, err := schema.ConcatMessages(fullMsgs)
			if err != nil {
				fmt.Println("error concatenating messages: ", err.Error())
			} else {
				fmt.Println(fullMsg)
			}
		}()

	outer:
		for {
			select {
			case <-ctx.Done():
				fmt.Println("context done", ctx.Err())
				return
			default:
				chunk, err := srs[1].Recv()
				if err != nil {
					if errors.Is(err, io.EOF) {
						break outer
					}
				}

				fullMsgs = append(fullMsgs, chunk)
			}
		}
	}()

	return srs[0], nil
}
