package main

import (
	"context"
	"fmt"

	"github.com/octokerbs/chronocode/internal/service"
)

func main() {
	ctx := context.Background()
	application := service.NewApplication(ctx)
	fmt.Println(application)
}
