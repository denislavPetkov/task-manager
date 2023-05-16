package main

import (
	"fmt"
	"os"

	"github.com/denislavpetkov/task-manager/pkg/controller"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}

	zap.ReplaceGlobals(logger)

	err = controller.NewController().Start()
	if err != nil {
		logger.Error(fmt.Sprintf("Controller failed to start, err: %v", err))
		os.Exit(1)
	}
}
