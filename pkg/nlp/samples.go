package nlp

var (
	CreateCommands = []string{
		"add", "add a",
		"create", "create a",
		"new",
	}

	EditCommands = []string{
		"edit", "edit the",
		"modify", "modify the",
		"change", "change the",
		"update", "update the",
	}

	DeleteCommands = []string{
		"delete", "delete the",
		"remove", "remove the",
		"cancel", "cancel the",
	}

	taskSamples = []string{
		"please {Command} new {TaskTitle}",
		"please {Command} task {TaskTitle}",
		"please {Command} the {TaskTitle}",
	}
)
