package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "CSV file in the format 'question,answer'")
	timeLimit := flag.Int("limit", 30, "time limit for the quiz in seconds")

	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file %v\n", *csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse the provided CSV file.\n", *csvFilename))
	}

	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)

		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C: // wait for the timer channel
			fmt.Println("Your time is up!")
			fmt.Printf("You scored %d out of %d\n", correct, len(problems))
			return
		case answer := <-answerCh: // wait for the user input channel to verify the answer
			if answer == p.a {
				correct++
			}
		}
	}
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

func exit(messge string) {
	fmt.Println(messge)
	os.Exit(1)
}
