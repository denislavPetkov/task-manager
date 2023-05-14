package nlp

import (
	"fmt"

	"github.com/shixzie/nlp"
)

var (
	nl *nlp.NL
)

type Command struct {
	Command   string
	TaskTitle string
}

func init() {
	nl = nlp.New()

	err := nl.RegisterModel(Command{}, taskSamples, nlp.WithTimeFormat("2006"))
	if err != nil {
		panic(err)
	}

	err = nl.Learn()
	if err != nil {
		panic(err)
	}
}

func GetCommand(commandString string) Command {
	command := nl.P(commandString)
	if cm, ok := command.(*Command); ok {
		return *cm
	} else {
		fmt.Println("Failed")
		return Command{}
	}
}
