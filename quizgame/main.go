package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	// set up flags
	fileFlag := flag.String("f", "problems.csv", "override the filename")
	timeFlag := flag.Int64("t", 30, "override the quiz timeout")
	flag.Parse()
	// read file
	records, err := ReadFile(fileFlag)
	var tc int
	if err != nil {
		return
	}
	// start timer
	timer := time.NewTimer(time.Duration(*timeFlag) * time.Second)
	qc := make(chan string)
	// start quiz
	go StartQuiz(records, &tc, qc)
	// wait for either the quiz to complete or the timer to run out
	select {
	case <-qc:
		fmt.Println("quiz ended")
	case <-timer.C:
		fmt.Println("time out")
	}
	fmt.Println("Total no of questions:", len(records))
	fmt.Println("Total no of correct answers:", tc)
}

func ReadFile(filename *string) ([][]string, error) {
	b, err := os.ReadFile(*filename)
	if err != nil {
		_ = fmt.Errorf("error reading file %v", err)
	}
	sReader := strings.NewReader(string(b))
	r := csv.NewReader(sReader)
	records, err := r.ReadAll()
	if err != nil {
		_ = fmt.Errorf("error in reading file %v", err)
		return nil, err
	}
	return records, nil
}

func StartQuiz(records [][]string, tc *int, qc chan string) {
	rand.Shuffle(len(records), func(i, j int) {
		records[i], records[j] = records[j], records[i]
	})
	//var tc, tq int
	for _, record := range records {
		// get question and answer from record
		q := record[0 : len(record)-1]
		a := record[len(record)-1]
		// prompt the question and get the answer
		var ua string
		fmt.Println(strings.Join(q, ","))
		_, err := fmt.Scanf("%s", &ua)
		if err != nil {
			_ = fmt.Errorf("error in getting input %v", err)
			return
		}
		if strings.EqualFold(ua, a) {
			*tc++
		}
	}
	qc <- "done"
}
