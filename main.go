package main

import (
	"fmt"
	"log"
	"os"
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

func parseExpressionWithStep(rng string) (string, int, error)
func parseExpression(rng string, min, max int) ([]int, error)
func filterByStep(minutes []int, step int) []int
