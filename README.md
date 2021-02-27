```go
package main

import "github.com/salah93/go-cron"
import "os/exec"

func main() {
	job := cron.NewJob()
	item := &cron.Item{
		Command: exec.Command("touch", "/tmp/x.txt"),
		Comment: "testing the water",
		Time: &cron.ItemTime{
			Minute:     "*",
			Hour:       "*",
			DayOfMonth: "*",
			Month:      "*",
			WeekDay:    "*",
		},
	}
	job.AddItem(item)
	job.Save()
}
```
