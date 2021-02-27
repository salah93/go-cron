package cron

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const croncmd = "/usr/bin/crontab"

func (j *Job) AddItem(i *Item) {
	j.Items = append(j.Items, i)
}

func (j *Job) RemoveItemsByComment(comment string) {
	var items []*Item
	for _, item := range j.Items[:] {
		if strings.Contains(item.Comment, comment) {
			continue
		}
		items = append(items, item)
	}
	j.Items = items
}

func (j *Job) Save() {
	f, err := ioutil.TempFile("", "go-cron-*")
	if err != nil {
		panic(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	for _, env := range j.Env {
		f.WriteString(fmt.Sprintf("%s\n", env))
	}
	for _, item := range j.Items {
		if item.Raw != "" {
			f.WriteString(fmt.Sprintf("%s\n", item.Raw))
		} else {
			f.WriteString(fmt.Sprintf("%s %s %s %s %s ", item.Time.Minute, item.Time.Hour, item.Time.DayOfMonth, item.Time.Month, item.Time.WeekDay))
			f.WriteString(fmt.Sprintf("%s ", strings.Join(item.Command.Env, " ")))
			f.WriteString(strings.Join(item.Command.Args, " "))
			if item.Comment != "" {
				f.WriteString(fmt.Sprintf(" # %s", item.Comment))
			}
			f.WriteString("\n")
		}
	}
	f.Sync()

	cmd := exec.Command(croncmd, f.Name())
	err = cmd.Run()
}

func NewJob() *Job {
	oldCronJobs, err := exec.Command(croncmd, "-l").Output()
	//job := Job{Items: []*Item{}, Comments: []string{}, Env: []string{}}
	job := Job{}
	if err == nil {
		cronEntriesSplitUp := strings.Split(string(oldCronJobs), "\n")
		for _, entry := range cronEntriesSplitUp {
			env := EnvPattern.FindString(entry)
			if env != "" {
				job.Env = append(job.Env, env)
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
	return &job
}
