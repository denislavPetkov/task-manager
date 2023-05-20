package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/denislavpetkov/task-manager/pkg/constants"
	"github.com/denislavpetkov/task-manager/pkg/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	csrf "github.com/utrack/gin-csrf"
)

const (
	tasksHtml    = "tasks.html"
	newTaskHtml  = "newTask.html"
	editTaskHtml = "editTask.html"

	taskKey = "task"
)

func (c *controller) getUserFromSession(gc *gin.Context) string {
	return sessions.Default(gc).Get(constants.SessionUserKey).(string)
}

func (c *controller) getTasks(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)

	currentUser := c.getUserFromSession(gc)

	tasks, err := c.taskDb.GetTasks(currentUser)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get tasks from db, error: %v", err))
		gc.HTML(http.StatusInternalServerError, loginHtml, gin.H{
			constants.CsrfKey: csrfToken,
			errorKey:          serverErrorErrMsg,
		})
		return
	}

	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to marshal tasks, error: %v", err))
		gc.HTML(http.StatusInternalServerError, loginHtml, gin.H{
			constants.CsrfKey: csrfToken,
			errorKey:          serverErrorErrMsg,
		})
		return
	}

	gc.HTML(http.StatusOK, tasksHtml, gin.H{
		constants.CsrfKey: csrfToken,
		"tasks":           string(tasksJson),
	})
}

func (c *controller) getNewTask(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)
	gc.HTML(http.StatusOK, newTaskHtml, gin.H{constants.CsrfKey: csrfToken})
}

func (c *controller) postNewTask(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)

	title := gc.PostForm("title")
	if strings.TrimSpace(title) == "" {
		logger.Info("Empty task title")
		gc.HTML(http.StatusBadRequest, newTaskHtml, gin.H{
			constants.CsrfKey: csrfToken,
			errorKey:          "Empty title",
		})
		return
	}

	if len(title) > 30 {
		logger.Info("Tittle exceeds 30 characters")
		gc.HTML(http.StatusBadRequest, newTaskHtml, gin.H{
			constants.CsrfKey: csrfToken,
			errorKey:          "Tittle exceeds 30 characters",
		})
		return
	}

	description := gc.PostForm("description")
	if strings.TrimSpace(description) == "" {
		logger.Info("Empty task description")
		gc.HTML(http.StatusBadRequest, newTaskHtml, gin.H{
			constants.CsrfKey: csrfToken,
			errorKey:          "Empty description",
		})
		return
	}

	deadline := gc.PostForm("deadline")
	category := gc.PostForm("category")
	tags := gc.PostFormArray("tags[]")
	notificationDeadline := gc.PostForm("notificationDeadline")

	notificationDeadlineDuration, err := time.ParseDuration(notificationDeadline)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to parse notification deadline, error: %v", err))

		gc.HTML(http.StatusInternalServerError, newTaskHtml, gin.H{
			constants.CsrfKey: csrfToken,
			errorKey:          serverErrorErrMsg,
		})

		return
	}

	title = strings.ReplaceAll(strings.Trim(title, " "), " ", "-")

	task := model.Task{
		Title:                title,
		Description:          description,
		Category:             category,
		Tags:                 tags,
		Deadline:             deadline,
		Completed:            false,
		NotificationDeadline: notificationDeadlineDuration,
	}

	currentUser := c.getUserFromSession(gc)

	err = c.taskDb.AddTask(currentUser, task)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			logger.Info(fmt.Sprintf("Task with title %s already exists", task.Title))

			gc.HTML(http.StatusBadRequest, newTaskHtml, gin.H{
				constants.CsrfKey: csrfToken,
				errorKey:          "Task with that title already exists",
				taskKey:           task,
			})

			return
		}

		logger.Error(fmt.Sprintf("Failed to add a task to db, error: %v", err))

		gc.HTML(http.StatusInternalServerError, newTaskHtml, gin.H{
			constants.CsrfKey: csrfToken,
			errorKey:          serverErrorErrMsg,
			taskKey:           task,
		})

		return
	}

	logger.Info("Added a new task to db")

	gc.HTML(http.StatusOK, newTaskHtml, gin.H{
		constants.CsrfKey: csrfToken,
		successKey:        "Added a task successful!",
		taskKey:           task,
	})
}

func (c *controller) deleteTask(gc *gin.Context) {
	title := gc.Param("title")

	currentUser := c.getUserFromSession(gc)

	err := c.taskDb.DeleteTask(currentUser, title)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to delete a task from db, error: %v", err))
		gc.HTML(http.StatusInternalServerError, newTaskHtml, gin.H{errorKey: serverErrorErrMsg})
		return
	}

	logger.Info("Deleted a task from db")

	gc.Status(http.StatusOK)
}

func (c *controller) getUpdateTask(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)
	gc.HTML(http.StatusOK, editTaskHtml, gin.H{constants.CsrfKey: csrfToken})
}

func (c *controller) postUpdateTask(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)

	title := gc.Param("title")

	description := gc.PostForm("description")
	if strings.TrimSpace(description) == "" {
		logger.Info("Empty task description")
		gc.HTML(http.StatusBadRequest, editTaskHtml, gin.H{
			constants.CsrfKey: csrfToken,
			errorKey:          "Empty description",
		})
		return
	}

	category := gc.PostForm("category")
	tags := gc.PostFormArray("tags[]")
	isCompletedString := gc.PostForm("completed")
	isCompleted, err := strconv.ParseBool(isCompletedString)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to parse 'completed' boolean, error: %v", err))
		gc.HTML(http.StatusInternalServerError, editTaskHtml, gin.H{
			constants.CsrfKey: csrfToken,
			errorKey:          serverErrorErrMsg,
		})
		return
	}

	deadline := gc.PostForm("deadline")

	notificationDeadline := gc.PostForm("notificationDeadline")

	notificationDeadlineDuration, err := time.ParseDuration(notificationDeadline)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to parse notification deadline, error: %v", err))

		gc.HTML(http.StatusInternalServerError, editTaskHtml, gin.H{
			constants.CsrfKey: csrfToken,
			errorKey:          serverErrorErrMsg,
		})

		return
	}

	currentUser := c.getUserFromSession(gc)

	task := model.Task{
		Title:                title,
		Description:          description,
		Category:             category,
		Tags:                 tags,
		Deadline:             deadline,
		Completed:            isCompleted,
		NotificationDeadline: notificationDeadlineDuration,
	}

	err = c.taskDb.UpdateTask(currentUser, title, task)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to update task in db, error: %v", err))
		gc.HTML(http.StatusInternalServerError, editTaskHtml, gin.H{
			constants.CsrfKey: csrfToken,
			errorKey:          serverErrorErrMsg,
		})
		return
	}

	logger.Info("Updated a task in db")

	gc.HTML(http.StatusOK, editTaskHtml, gin.H{
		constants.CsrfKey: csrfToken,
		successKey:        "Edited the task successful!",
		taskKey:           task,
	})
}

func (c *controller) postCompleteTask(gc *gin.Context) {
	title := gc.Param("title")

	completed := gc.PostForm("completed")

	currentUser := c.getUserFromSession(gc)

	task, err := c.taskDb.GetTask(currentUser, title)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get task db, error: %v", err))
		gc.Status(http.StatusInternalServerError)
		return
	}

	completedBool, err := strconv.ParseBool(completed)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to parse 'completed' boolean, error: %v", err))
		gc.Status(http.StatusInternalServerError)
		return
	}

	task.Completed = completedBool

	err = c.taskDb.UpdateTask(currentUser, title, task)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to update task in db, error: %v", err))
		gc.Status(http.StatusInternalServerError)
		return
	}

	logger.Info("Updated a task in db successful")

	gc.JSON(http.StatusOK, gin.H{
		"completed": completedBool,
	})
}
