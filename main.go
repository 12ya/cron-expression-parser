package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("Usage: go run . <cron expression>", os.Args[0])
	}

	cronExpression := os.Args[1]
	if cronExpression == "" {
		return fmt.Errorf("cron expression is empty")
	}

	fmt.Println(cronExpression, "-----expresssion-----")

	parts := strings.Fields(cronExpression)

	minutes, err := parse(parts[0], 0, 59)
	if err != nil {
		return fmt.Errorf("parse minutes: %s", err)
	}

	fmt.Println(minutes)

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

		rangeMinutes, err := parseExpression(rng, min, max)
		if err != nil {
			return nil, err
		}

		if step > 0 {
			rangeMinutes = filterByStep(rangeMinutes, step)
		}

		values = append(values, rangeMinutes...)
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
