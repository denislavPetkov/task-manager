package model

import "time"

type Task struct {
	Title                string        `json:"title" binding:"required"`
	Deadline             string        `json:"deadline" binding:"required"`
	Description          string        `json:"description" binding:"required"`
	Category             string        `json:"category" binding:"-"`
	Tags                 []string      `json:"tags" binding:"-"`
	Completed            bool          `json:"completed" binding:"-"`
	NotificationDeadline time.Duration `json:"notificationDeadline" binding:"required"`
}
