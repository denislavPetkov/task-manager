package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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
	openaiApiKey := os.Getenv(OPENAI_API_KEY)

	ctx := context.Background()

	s := openai.NewSession(openaiApiKey)
	client := audio.NewClient(s, whisperModel)
	resp, err := client.CreateTranscription(ctx, &audio.CreateTranscriptionParams{
		Language:    whisperLanguage,
		Audio:       gc.Request.Body,
		AudioFormat: whisperAudioFormat,
	})
	if err != nil {
		log.Fatalf("error transcribing: %v", err)
		gc.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Server error"})
		return
	}

	voiceCommand := strings.ToLower(resp.Text)
	log.Println(voiceCommand)

	voiceCommand = "new task"

	command := nlp.GetCommand(voiceCommand)
	fmt.Println(command)

	switch {
	case slices.Contains(nlp.CreateCommands, command.Command):
		gc.JSON(http.StatusFound, gin.H{"redirect": "/tasks/new"})

	case slices.Contains(nlp.EditCommands, command.Command):
		gc.JSON(http.StatusFound, gin.H{
			"redirect": fmt.Sprintf("/tasks/edit/%s", command.TaskTitle),
			"edit":     command.TaskTitle,
		})

	case slices.Contains(nlp.DeleteCommands, command.Command):
		gc.JSON(http.StatusFound, gin.H{
			"redirect": fmt.Sprintf("/tasks/delete/%s", command.TaskTitle),
		})

	case slices.Contains(nlp.MarkAsCommands, command.Command):

		gc.JSON(http.StatusFound, gin.H{
			"redirect": fmt.Sprintf("/tasks/delete/%s", command.TaskTitle),
		})
	}
}