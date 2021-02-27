package cron

import (
	"os/exec"
)

//ItemTime - represents the time for a cron entry
type ItemTime struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	WeekDay    string
}

//Item - repesenting cron entry
type Item struct {
	Command *exec.Cmd
	Comment string
	Time    *ItemTime
	Raw     string
}

//Macros
const (
	Annually = "annually"
	Yearly   = "yearly"
	Monthly  = "monthly"
	Weekly   = "weekly"
	Daily    = "daily"
	Hourly   = "hourly"
)

// Job - representing cronjob file
type Job struct {
	Items []*Item
	// Env - environment variables set in cronfile
	// each entry is of the form key=value
	Env      []string
	Comments []string
}
