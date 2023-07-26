package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/binarycraft007/lunar"
)

func main() {
	dateFlag := flag.String("date", "", "the date to use (format: yyyy-mm-dd)")

	// Parse the command-line arguments
	flag.Parse()

	if *dateFlag == "" {
		panic("date input required")
	}

	timeIn, err := time.Parse("2006-01-02", *dateFlag)
	if err != nil {
		panic(err)
	}

	l, err := lunar.NewLunar()
	if err != nil {
		panic(err)
	}

	lunarDate, err := l.TimeToLunar(timeIn)
	if err != nil {
		panic(err)
	}

	fmt.Println(*lunarDate)
}
