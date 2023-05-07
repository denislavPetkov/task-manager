package main

import "github.com/denislavpetkov/task-manager/pkg/controller"

func main() {
	err := controller.NewController().Start()
	println(err.Error())
}
