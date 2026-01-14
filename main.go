// 1-5/2,10-15/3 * * * * /bin/echo

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const totalFields = 6

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("Usage: go run . <cron expression>: %s", os.Args[0])
	}

	cronExpression := os.Args[1]
	if cronExpression == "" {
		return fmt.Errorf("cron expression is empty")
	}

	parts := strings.Fields(cronExpression)
	if len(parts) != totalFields {
		return fmt.Errorf("invalid cron expression: %s - must have exactly %d fields, e.g.: */15 0 1,15 * 1-5 /bin/echo", cronExpression, totalFields)
	}

	minutes, err := parse(parts[0], 0, 59)
	if err != nil {
		return fmt.Errorf("error parsing minutes: %s", err)
	}
	hours, err := parse(parts[1], 0, 23)
	if err != nil {
		return fmt.Errorf("error parsing hours: %s", err)
	}
	daysOfMonth, err := parse(parts[2], 1, 31)
	if err != nil {
		return fmt.Errorf("error parsing days of month: %s", err)
	}
	months, err := parse(parts[3], 1, 12)
	if err != nil {
		return fmt.Errorf("error parsing months: %s", err)
	}
	daysOfWeek, err := parse(parts[4], 0, 6)
	if err != nil {
		return fmt.Errorf("error parsing days of week: %s", err)
	}

	command := parts[5]

	draw(minutes, hours, daysOfMonth, months, daysOfWeek, command)
	return nil
}

func parse(period string, min, max int) ([]int, error) {
	if period == "" {
		return nil, nil
	}

	var (
		values []int
		rng    string
		step   int
		err    error
	)

	ranges := strings.Split(period, ",")
	for _, rng = range ranges {
		if strings.Contains(rng, "/") {
			rng, step, err = parseExpressionWithStep(rng)
			if err != nil {
				return nil, err
			}
		}

		rangeValues, err := parseExpression(rng, min, max)
		if err != nil {
			return nil, err
		}

		if step > 0 {
			rangeValues = filterByStep(rangeValues, step)
		}

		values = append(values, rangeValues...)
	}

	return values, nil
}

func parseExpressionWithStep(rng string) (string, int, error) {
	parts := strings.Split(rng, "/")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid range/step: %s", rng)
	}

	rangeStr := parts[0]
	stepStr := parts[1]

	step, err := strconv.Atoi(stepStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid step: %s", stepStr)
	}

	return rangeStr, step, nil
}

func parseExpression(rng string, min, max int) ([]int, error) {
	if rng == "*" {
		values := make([]int, 0, max-min+1)
		for i := min; i <= max; i++ {
			values = append(values, i)
		}
		return values, nil
	}

	if !strings.Contains(rng, "-") {
		val, err := strconv.Atoi(rng)
		if err != nil {
			return nil, err
		}
		if val < min || val > max {
			return nil, fmt.Errorf("value is not in range (%d-%d): %d", min, max, val)
		}
		return []int{val}, nil
	}

	parts := strings.Split(rng, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid range: %s", rng)
	}

	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid range(%s) star: %s", rng, parts[0])
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid range(%s) end: %s", rng, parts[1])
	}

	if start < min || start > max {
		return nil, fmt.Errorf("invalid range(%s) start: %d", rng, start)
	}

	if end < min || end > max {
		return nil, fmt.Errorf("invalid range(%s) end: %d", rng, end)
	}

	if start > end {
		return nil, fmt.Errorf("invalid range(%s): %d > %d, want start < end", rng, start, end)
	}

	values := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		values = append(values, i)
	}

	return values, nil
}

func filterByStep(values []int, step int) []int {
	var filtered []int
	for i, v := range values {
		if i%step == 0 {
			filtered = append(filtered, v)
		}
	}

	return filtered
}

func draw(minutes, hours, daysOfMonth, months, daysOfWeek []int, command string) {
	fmt.Println("minute\t", strings.Join(intsToStrs(minutes), " "))
	fmt.Println("hour\t", strings.Join(intsToStrs(hours), " "))
	fmt.Println("day of month\t", strings.Join(intsToStrs(minutes), " "))
	fmt.Println("month\t", strings.Join(intsToStrs(minutes), " "))
	fmt.Println("day of week\t", strings.Join(intsToStrs(minutes), " "))
	fmt.Println("command\t", strings.Join(intsToStrs(minutes), " "))
}

func intsToStrs(ints []int) []string {
	var result []string
	for _, i := range ints {
		result = append(result, strconv.Itoa(i))
	}

	return result
}
