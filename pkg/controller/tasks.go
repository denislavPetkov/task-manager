package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/denislavpetkov/task-manager/pkg/constants"
	"github.com/denislavpetkov/task-manager/pkg/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	csrf "github.com/utrack/gin-csrf"
)

func (c *controller) getUserFromSession(gc *gin.Context) string {
	return sessions.Default(gc).Get(constants.SessionUserKey).(string)
}

func (c *controller) getTasks(gc *gin.Context) {
	currentUser := c.getUserFromSession(gc)

	tasks, err := c.taskDb.GetTasks(currentUser)
	if err != nil {
		gc.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Server error"})
		return
	}
	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		gc.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Server error"})
		return
	}

	csrfToken := csrf.GetToken(gc)

	gc.HTML(http.StatusOK, "tasks.html", gin.H{
		"Tasks":           string(tasksJson),
		constants.CsrfKey: csrfToken,
	})
}

func (c *controller) getNewTask(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)
	gc.HTML(http.StatusOK, "newTask.html", gin.H{constants.CsrfKey: csrfToken})
}

func (c *controller) postNewTask(gc *gin.Context) {
	title := gc.PostForm("title")
	if strings.TrimSpace(title) == "" {
		gc.HTML(http.StatusBadRequest, "newTask.html", gin.H{"error": "Empty title"})
		return
	}

	description := gc.PostForm("description")
	if strings.TrimSpace(description) == "" {
		gc.HTML(http.StatusBadRequest, "newTask.html", gin.H{"error": "Empty description"})
		return
	}

	deadline := gc.PostForm("deadline")
	category := gc.PostForm("category")
	tags := gc.PostFormArray("tags[]")

	task := model.Task{
		Title:       title,
		Description: description,
		Category:    category,
		Tags:        tags,
		Deadline:    deadline,
		Completed:   false,
	}

	currentUser := c.getUserFromSession(gc)

	err := c.taskDb.AddTask(currentUser, task)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			if err != nil {
				gc.HTML(http.StatusInternalServerError, "newTask.html", gin.H{
					"error": "Server error",
					"task":  task,
				})
				return
			}

			gc.HTML(http.StatusBadRequest, "newTask.html", gin.H{
				"error": "Task with that title already exists",
				"task":  task,
			})
			return
		}
		gc.HTML(http.StatusInternalServerError, "newTask.html", gin.H{
			"error": "Server error",
			"task":  task,
		})
		return
	}

	gc.HTML(http.StatusOK, "newTask.html", gin.H{
		"success": "Added a Task successful!",
		"task":    task,
	})
}

func (c *controller) deleteTask(gc *gin.Context) {
	title := gc.Param("title")

	currentUser := c.getUserFromSession(gc)

	err := c.taskDb.DeleteTask(currentUser, title)
	if err != nil {
		gc.HTML(http.StatusInternalServerError, "newTask.html", gin.H{"error": "Server error"})
		return
	}

	gc.HTML(http.StatusOK, "newTask.html", gin.H{})
}

func (c *controller) getUpdateTask(gc *gin.Context) {
	csrfToken := csrf.GetToken(gc)
	gc.HTML(http.StatusOK, "editTask.html", gin.H{constants.CsrfKey: csrfToken})
}

func (c *controller) postUpdateTask(gc *gin.Context) {
	title := gc.Param("title")

	description := gc.PostForm("description")
	if strings.TrimSpace(description) == "" {
		gc.HTML(http.StatusBadRequest, "editTask.html", gin.H{"error": "Empty description"})
		return
	}

	category := gc.PostForm("category")
	tags := gc.PostFormArray("tags[]")
	isCompletedString := gc.PostForm("completed")
	isCompleted, err := strconv.ParseBool(isCompletedString)
	if err != nil {
		gc.HTML(http.StatusInternalServerError, "editTask.html", gin.H{"error": "Server error"})
		return
	}

	deadline := gc.PostForm("deadline")

	currentUser := c.getUserFromSession(gc)

	task := model.Task{
		Title:       title,
		Description: description,
		Category:    category,
		Tags:        tags,
		Deadline:    deadline,
		Completed:   isCompleted,
	}

	err = c.taskDb.UpdateTask(currentUser, title, task)
	if err != nil {
		gc.HTML(http.StatusInternalServerError, "editTask.html", gin.H{"error": "Server error"})
		return
	}

	gc.HTML(http.StatusOK, "editTask.html", gin.H{
		"success": "Added a Task successful!",
		"task":    task,
	})
}

func (c *controller) postCompleteTask(gc *gin.Context) {
	title := gc.Param("title")

	completed := gc.PostForm("completed")

	currentUser := c.getUserFromSession(gc)

	task, err := c.taskDb.GetTask(currentUser, title)
	if err != nil {
		gc.Status(http.StatusInternalServerError)
		return
	}

	completedBool, err := strconv.ParseBool(completed)
	if err != nil {
		gc.Status(http.StatusInternalServerError)
		return
	}

	task.Completed = completedBool

	err = c.taskDb.UpdateTask(currentUser, title, task)
	if err != nil {
		gc.Status(http.StatusInternalServerError)
		return
	}

	gc.JSON(http.StatusOK, gin.H{
		"completed": completedBool,
	})
}
