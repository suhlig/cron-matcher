package cron_matcher

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func Matches(expression string, t time.Time) (bool, error) {
	fields := strings.Fields(expression)

	if len(fields) != 5 {
		return false, fmt.Errorf("invalid number of fields in cron expression: expected 5, got %d", len(fields))
	}

	checks := []struct {
		field string
		value int
		min   int
		max   int
	}{
		{fields[0], t.Minute(), 0, 59},
		{fields[1], t.Hour(), 0, 23},
		{fields[2], t.Day(), 1, 31},
		{fields[3], int(t.Month()), 1, 12},
		{fields[4], int(t.Weekday()), 0, 6},
	}

	for i, check := range checks {
		match, err := matchField(check.field, check.value, check.min, check.max)

		if err != nil {
			return false, err
		}

		if !match {
			if i == 2 || i == 4 { // day of month and day of week are special
				continue
			}

			return false, nil
		}
	}

	return true, nil
}

func matchField(expr string, value, min, max int) (bool, error) {
	if expr == "*" {
		return true, nil
	}

	for _, part := range strings.Split(expr, ",") {
		if strings.Contains(part, "/") {
			if match, err := matchStep(part, value, min, max); err != nil {
				return false, err
			} else if match {
				return true, nil
			}
		} else if strings.Contains(part, "-") {
			if match, err := matchRange(part, value, min, max); err != nil {
				return false, err
			} else if match {
				return true, nil
			}
		} else {
			if num, err := strconv.Atoi(part); err != nil {
				return false, fmt.Errorf("invalid number: %s", part)
			} else if num < min || num > max {
				return false, fmt.Errorf("number %d outside allowed bounds %d-%d", num, min, max)
			} else if num == value {
				return true, nil
			}
		}
	}

	return false, nil
}

func matchRange(expr string, value, min, max int) (bool, error) {
	parts := strings.Split(expr, "-")

	if len(parts) != 2 {
		return false, fmt.Errorf("invalid range expression: %s", expr)
	}

	start, err := strconv.Atoi(parts[0])

	if err != nil {
		return false, fmt.Errorf("invalid start of range: %s", parts[0])
	}

	end, err := strconv.Atoi(parts[1])

	if err != nil {
		return false, fmt.Errorf("invalid end of range: %s", parts[1])
	}

	if start < min || end > max || start > end {
		return false, fmt.Errorf("range %d-%d outside allowed bounds %d-%d", start, end, min, max)
	}

	return value >= start && value <= end, nil
}

func matchStep(expr string, value, min, max int) (bool, error) {
	parts := strings.Split(expr, "/")

	if len(parts) != 2 {
		return false, fmt.Errorf("invalid step expression: %s", expr)
	}

	var start int
	var err error

	if parts[0] == "*" {
		start = min
	} else {
		start, err = strconv.Atoi(parts[0])

		if err != nil {
			return false, fmt.Errorf("invalid start of step: %s", parts[0])
		}
	}

	step, err := strconv.Atoi(parts[1])

	if err != nil {
		return false, fmt.Errorf("invalid step: %s", parts[1])
	}

	if step == 0 {
		return false, fmt.Errorf("step cannot be zero")
	}

	for i := start; i <= max; i += step {
		if i == value {
			return true, nil
		}
		if i > value {
			break
		}
	}

	return false, nil
}
