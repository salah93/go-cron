package cron

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//MacroPattern pattern to match cron strings
var MacroPattern = regexp.MustCompile("@(annually|yearly|monthly|weekly|daily|hourly)")

//TimePattern pattern to match time strings
var TimePattern = regexp.MustCompile("(\\d+|\\*) (\\d+|\\*) (\\d+|\\*) (\\d+|\\*) (\\d+|\\*)")

//EnvPattern pattern to match environment strings
var EnvPattern = regexp.MustCompile("(\\w+)=((?:\\w|[$:/])+)")

//  ((\w+=\w+ )*(\w+))")

//NewItem - parse a cron entry
func NewItem(cronEntry string) *Item {
	// var item Item
	// var comment string
	var itemTime *ItemTime
	cronEntrySplitUp := strings.Split(cronEntry, " ")
	if len(cronEntrySplitUp) == 0 {
		return nil
	}
	macro := MacroPattern.FindString(cronEntrySplitUp[0])
	if macro != "" {
		itemTime = parseMacro(macro)
		cronEntrySplitUp = cronEntrySplitUp[1:]
	} else {
		timeString := strings.Join(cronEntrySplitUp[:5], " ")
		itemTime = parseTime(timeString)
		cronEntrySplitUp = cronEntrySplitUp[5:]
	}
	cmd, comment := parseCommand(cronEntrySplitUp)
	return &Item{
		Command: cmd,
		Time:    itemTime,
		Comment: comment,
		Raw:     cronEntry,
	}
}

//parseMacro - transforms macro statement to itemtime object
func parseMacro(macro string) *ItemTime {
	if len(macro) == 0 {
		return nil
	} else if macro[0] == '@' {
		macro = macro[1:]
	}
	switch macro {
	case Annually:
		fallthrough
	case Yearly:
		// 0 0 1 1 *
		return &ItemTime{
			Minute:     "0",
			Hour:       "0",
			DayOfMonth: "1",
			Month:      strconv.Itoa(int(time.January)),
			WeekDay:    "*",
		}
	case Monthly:
		// 0 0 1 * *
		return &ItemTime{
			Minute:     "0",
			Hour:       "0",
			DayOfMonth: "1",
			Month:      "*",
			WeekDay:    "*",
		}
	case Weekly:
		// 0 0 * * 0
		return &ItemTime{
			Minute:     "0",
			Hour:       "0",
			DayOfMonth: "*",
			Month:      "*",
			WeekDay:    strconv.Itoa(int(time.Sunday)),
		}
	case Daily:
		// 0 0 * * *
		return &ItemTime{
			Minute:     "0",
			Hour:       "0",
			DayOfMonth: "*",
			Month:      "*",
			WeekDay:    "*",
		}
	case Hourly:
		// 0 * * * *
		return &ItemTime{
			Minute:     "0",
			Hour:       "*",
			DayOfMonth: "*",
			Month:      "*",
			WeekDay:    "*",
		}
	default:
		return nil
	}
}

/* parseTime
* * * * *
10 * 10 * *
*/
func parseTime(timeString string) *ItemTime {
	matches := TimePattern.FindStringSubmatch(timeString)
	if matches == nil {
		return nil
	}
	return &ItemTime{
		Minute:     matches[1],
		Hour:       matches[2],
		DayOfMonth: matches[3],
		Month:      matches[4],
		WeekDay:    matches[5],
	}

}

/* parseCommand
ENV=staging /usr/bin/python /opt/app/main.py
/usr/bin/bash /opt/app/script.sh
*/
func parseCommand(cronEntrySplitUp []string) (*exec.Cmd, string) {
	commandEnv := []string{}
	command := []string{}
	var comment string
	for index, entry := range cronEntrySplitUp {
		env := EnvPattern.FindString(entry)
		if env != "" {
			commandEnv = append(commandEnv, env)
		} else if strings.HasPrefix(entry, "#") {
			comment = strings.TrimLeft(strings.Join(cronEntrySplitUp[index:], " ")[1:], " ")
			break
		} else {
			command = append(command, entry)
		}
	}
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = commandEnv
	return cmd, comment
}
