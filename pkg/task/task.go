package task

import (
	"time"
)

type Priority int

const (
	High Priority = iota
	Mid
	Low
	None
)

func LabelToPriority(str string) Priority {
	if str == "High" {
		return High
	} else if str == "Middle" {
		return Mid
	} else if str == "Low" {
		return Low
	}
	return None
}

type Task struct {
	ProviderName string
	TaskID       string
	Group        string
	Title        string
	Priority     Priority
	Deadline     *time.Time
	Url          string
}
