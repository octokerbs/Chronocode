package main

import (
	"context"
	"fmt"

	"github.com/octokerbs/chronocode-backend/internal2/service"
)

func main() {
	ctx := context.Background()
	application := service.NewApplication(ctx)
	fmt.Println(application)
}
