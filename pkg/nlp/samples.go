package nlp

var (
	CreateCommands = []string{
		"add", "create", "new",
	}

	EditCommands = []string{
		"edit", "modify", "change", "update",
	}

	DeleteCommands = []string{
		"delete", "remove", "cancel",
	}

	taskSamples = []string{
		"{Command} task {TaskTitle}",
		"{Command} {TaskTitle} task",
		"{Command} a new task",
		"{Command} task",
	}
)
