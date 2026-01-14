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
		vals := make([]int, 0, max-min+1)
	}
	return nil, nil
}

func filterByStep(minutes []int, step int) []int
