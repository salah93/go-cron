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
