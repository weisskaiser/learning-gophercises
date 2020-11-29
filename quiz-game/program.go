package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {

	fileName := flag.String("fileName", "problems.csv", "Use to specify the questions csv file (format = question,answer)")
	timeout := flag.Int("timeout", 30, "timeout to exit the quiz")

	flag.Parse()

	file, openErr := os.Open(*fileName)
	defer file.Close()

	if openErr != nil {
		exit(openErr.Error())
	}

	csvReader := csv.NewReader(file)

	var userInput string
	var correct, wrong int

	csvReader.FieldsPerRecord = 2

	go func() {
		time.Sleep(time.Second * time.Duration(*timeout))
		exit("\nTime is over Anakin")
	}()

	for {

		record, readErr := csvReader.Read()

		if readErr == io.EOF {
			break
		}

		fmt.Printf("%s? ", record[0])
		fmt.Scanln(&userInput)

		if strings.TrimSpace(userInput) == strings.TrimSpace(record[1]) {
			correct++
		} else {
			wrong++
		}
	}

	fmt.Println("correct", correct)
	fmt.Println("wrong", wrong)
	fmt.Println("total", correct+wrong)

}
