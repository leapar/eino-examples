package main

import (
	"context"
	"fmt"
	"testgraph/testgraph"
)

func main() {
	ctx := context.Background()
	testgraph.Builddemo(ctx)

	fmt.Println("sss")
}
