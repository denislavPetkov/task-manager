package nlp

var (
	CreateCommands = []string{
		"add", "create", "new",
	}

	EditCommands = []string{
		"edit", "modify", "change", "update",
	}

	MarkAsCommands = []string{
		"set", "mark", "finish",
	}

	DoneState = []string{
		"completed", "done",
	}

	NotDoneState = []string{
		"not completed", "not done", "in progress",
	}

	PossibleState = append(DoneState, NotDoneState...)

	DeleteCommands = []string{
		"delete", "remove", "cancel",
	}

	FilterCommands = []string{
		"filter", "get",
	}

	taskSamples = []string{
		"{Command} task {TaskTitle}",
		"{Command} task {TaskTitle}",
		"{Command} task {TaskTitle}",
		"{Command} task {TaskTitle}",
		"{Command} task {TaskTitle}",
		"{Command} task {TaskTitle}",
		"{Command} task {TaskTitle}",
		"{Command} task {TaskTitle} to {State}",
		"{Command} task {TaskTitle} as {State}",
		"{Command} task {TaskTitle}",
		"{Command} tasks by {Filter} {FilterValue}",
		"{Command} tasks with {Filter} {FilterValue}",
	}
)
