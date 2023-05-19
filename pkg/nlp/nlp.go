package nlp

import (
	"fmt"

	"github.com/shixzie/nlp"
	"go.uber.org/zap"
)

var (
	nl *nlp.NL

	logger *zap.Logger
)

type Command struct {
	Command   string
	TaskTitle string
}

func init() {
	logger = zap.L().Named("nlp")

	nl = nlp.New()

	err := nl.RegisterModel(Command{}, taskSamples, nlp.WithTimeFormat("2006"))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to register nlp model, error: %v", err))
		panic(err)
	}

	err = nl.Learn()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to teach the nlp, error: %v", err))
		panic(err)
	}
}

func GetCommand(commandString string) Command {
	command := nl.P(commandString)
	if cm, ok := command.(*Command); ok {
		return *cm
	} else {
		logger.Info(fmt.Sprintf("Could not type assert command string %s to a command object", commandString))
		return Command{}
	}
}
