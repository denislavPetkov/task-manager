package controller

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/denislavpetkov/task-manager/pkg/nlp"
	"github.com/gin-gonic/gin"
	openai "github.com/rakyll/openai-go"
	"github.com/rakyll/openai-go/audio"
	"golang.org/x/exp/slices"
)

const (
	OPENAI_API_KEY = "OPENAI_API_KEY"

	whisperLanguage    = "en"
	whisperModel       = "whisper-1"
	whisperAudioFormat = "wav"
)

func (c *controller) postAudio(gc *gin.Context) {

	ctx := context.Background()

	s := openai.NewSession(c.openaiApiKey)
	client := audio.NewClient(s, whisperModel)
	resp, err := client.CreateTranscription(ctx, &audio.CreateTranscriptionParams{
		Language:    whisperLanguage,
		Audio:       gc.Request.Body,
		AudioFormat: whisperAudioFormat,
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to transcribe, error: %v", err))
		gc.HTML(http.StatusInternalServerError, loginHtml, gin.H{errorKey: serverErrorErrMsg})
		return
	}

	voiceCommand := strings.ToLower(resp.Text)
	trimmed := trimTaskTitle(voiceCommand)
	command := nlp.GetCommand(trimmed)

	logger.Info(fmt.Sprintf("Raw voice command: %s", voiceCommand))
	logger.Info(fmt.Sprintf("Trimmed command: %s", trimmed))
	logger.Info(fmt.Sprintf("Processed voice command: %v", command))

	switch {
	case slices.Contains(nlp.CreateCommands, command.Command):
		gc.JSON(http.StatusFound, gin.H{"redirect": "/tasks/new"})

	case slices.Contains(nlp.EditCommands, command.Command):
		taskTitle := trimTaskTitle(command.TaskTitle)

		gc.JSON(http.StatusFound, gin.H{
			"redirect": fmt.Sprintf("/tasks/edit/%s", taskTitle),
			"edit":     taskTitle,
		})

	case slices.Contains(nlp.DeleteCommands, command.Command):
		taskTitle := trimTaskTitle(command.TaskTitle)

		gc.JSON(http.StatusFound, gin.H{
			"redirect": fmt.Sprintf("/tasks/delete/%s", taskTitle),
			"delete":   taskTitle,
		})
	}
}

func trimTaskTitle(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, ".", ""), ",", "")
}
