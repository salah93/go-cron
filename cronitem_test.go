package cron

import "testing"

func TestParseMacroAnnually(t *testing.T) {
	itemTime := parseMacro("@annually")
	if itemTime == nil {
		t.Fatal("expected annual parsing, got nil")
	}
	if itemTime.Minute != "0" ||
		itemTime.Hour != "0" ||
		itemTime.DayOfMonth != "1" ||
		itemTime.Month != "1" ||
		itemTime.WeekDay != "*" {
		t.Fatal("does not match annual time")
	}
	yearlyItemTime := parseMacro("@yearly")
	if *yearlyItemTime != *itemTime {
		t.Fatal("yearly does not match annual")
	}
}

func TestParseTimeEveryMinute(t *testing.T) {
	itemTime := parseTime("* * * * *")
	if itemTime == nil {
		t.Fatal("expected annual parsing, got nil")
	}
	if itemTime.Minute != "*" ||
		itemTime.Hour != "*" ||
		itemTime.DayOfMonth != "*" ||
		itemTime.Month != "*" ||
		itemTime.WeekDay != "*" {
		t.Fatal("does not match set time")
	}
}

func TestParseTimeCustom(t *testing.T) {
	itemTime := parseTime("10 * * 10 *")
	if itemTime == nil {
		t.Fatal("expected annual parsing, got nil")
	}
	if itemTime.Minute != "10" ||
		itemTime.Hour != "*" ||
		itemTime.DayOfMonth != "*" ||
		itemTime.Month != "10" ||
		itemTime.WeekDay != "*" {
		t.Fatal("does not match set time")
	}
}

func TestParseCommandEnv(t *testing.T) {
	command := []string{"ENV=staging", "PATH=$PATH:/opt/bin", "/usr/local/bin/python3", "/opt/app/server.py", "--port", "7000"}
	cmd, comment := parseCommand(command)
	if comment != "" {
		t.Fatal("comment not parsed")
	}
	if len(cmd.Env) != 2 ||
		cmd.Env[0] != "ENV=staging" ||
		cmd.Env[1] != "PATH=$PATH:/opt/bin" {
		t.Fatal("env not parsed correctly")
	}
	if len(cmd.Args) != 4 ||
		cmd.Args[0] != "/usr/local/bin/python3" ||
		cmd.Args[1] != "/opt/app/server.py" ||
		cmd.Args[2] != "--port" ||
		cmd.Args[3] != "7000" {
		t.Fatal("could not parse arguments")
	}
}

func TestParseCommandEnvComment(t *testing.T) {
	command := []string{"ENV=staging", "PATH=$PATH:/opt/bin", "/usr/local/bin/python3", "/opt/app/server.py", "--port", "7000", "#", "run", "server"}
	cmd, comment := parseCommand(command)
	if comment != "run server" {
		t.Fatalf("comment not parsed: %s", comment)
	}
	if len(cmd.Env) != 2 ||
		cmd.Env[0] != "ENV=staging" ||
		cmd.Env[1] != "PATH=$PATH:/opt/bin" {
		t.Fatal("env not parsed correctly")
	}
	if len(cmd.Args) != 4 ||
		cmd.Args[0] != "/usr/local/bin/python3" ||
		cmd.Args[1] != "/opt/app/server.py" ||
		cmd.Args[2] != "--port" ||
		cmd.Args[3] != "7000" {
		t.Fatal("could not parse arguments")
	}
}
