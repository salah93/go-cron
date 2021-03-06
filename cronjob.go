package cron

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

//Job - representing cronjob file
type Job struct {
	Items []*Item
	// Env - environment variables set in cronfile
	// each entry is of the form key=value
	Env      []string
	Comments []string
}

//AddItem - adds a cron entry to the cronjob
func (j *Job) AddItem(i *Item) {
	j.Items = append(j.Items, i)
}

// RemoveItemsByComment - filter out cron entries by their comment
func (j *Job) RemoveItemsByComment(comment string) {
	var items []*Item
	for _, item := range j.Items {
		if !strings.Contains(item.Comment, comment) {
			items = append(items, item)
		}
	}
	j.Items = items
}

// Save - save cronjob to user's cron file
func (j *Job) Save() error {
	tempfile, err := ioutil.TempFile("", "go-cron-*")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tempfile.Name())
	defer tempfile.Close()

	if j.Env != nil {
		tempfile.WriteString(fmt.Sprintf("%s\n", strings.Join(j.Env, "\n")))
	}
	if j.Comments != nil {
		tempfile.WriteString(fmt.Sprintf("%s\n", strings.Join(j.Comments, "\n")))
	}
	for _, item := range j.Items {
		if item.Raw != "" {
			tempfile.WriteString(fmt.Sprintf("%s\n", item.Raw))
		} else {
			// time
			tempfile.WriteString(fmt.Sprintf("%s %s %s %s %s ", item.Time.Minute, item.Time.Hour, item.Time.DayOfMonth, item.Time.Month, item.Time.WeekDay))
			// env
			tempfile.WriteString(fmt.Sprintf("%s ", strings.Join(item.Command.Env, " ")))
			// command
			tempfile.WriteString(strings.Join(item.Command.Args, " "))

			// redirect
			stdout, stderr := getStdOutErrFileNames(item.Command)
			tempfile.WriteString(fmt.Sprintf(" > %s 2>%s", stdout, stderr))

			// comments
			if item.Comment != "" {
				tempfile.WriteString(fmt.Sprintf(" # %s", item.Comment))
			}
			tempfile.WriteString("\n")
		}
	}
	tempfile.Sync()

	cmd := exec.Command(CronCmd, tempfile.Name())
	return cmd.Run()
}

func getStdOutErrFileNames(cmd *exec.Cmd) (string, string) {
	stdout := os.DevNull
	stderr := os.DevNull
	if cmd.Stdout != nil {
		f, ok := cmd.Stdout.(*os.File)
		if ok {
			stdout = f.Name()
		}
	}
	if cmd.Stderr == cmd.Stdout {
		stderr = "&1"
	} else if cmd.Stderr != nil {
		f, ok := cmd.Stderr.(*os.File)
		if ok {
			stderr = f.Name()
		}
	}
	return stdout, stderr
}

//NewJob - get a Job object, grabs user's current cron entries to start
func NewJob() *Job {
	oldCronJobs, err := exec.Command(CronCmd, "-l").Output()
	job := new(Job)
	if err == nil {
		cronEntriesSplitUp := strings.Split(string(oldCronJobs), "\n")
		for _, entry := range cronEntriesSplitUp {
			envMatch := EnvPattern.FindString(entry)
			if envMatch != "" && envMatch == entry {
				job.Env = append(job.Env, entry)
			} else if strings.HasPrefix(entry, "#") {
				job.Comments = append(job.Comments, entry)
			} else {
				item := NewItem(entry)
				if item != nil {
					job.Items = append(job.Items, item)
				}
			}
		}
	}
	return job
}
