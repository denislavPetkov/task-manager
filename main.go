package main

import "github.com/denislavpetkov/task-manager/pkg/controller"

// "github.com/notes-project/server/pkg/controller"

func main() {
	_ = controller.NewController().Start()
}
