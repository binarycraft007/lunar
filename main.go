package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"golang.org/x/text/encoding/traditionalchinese"
)

func main() {
	timeIn, err := time.Parse("2006-01-02", "1996-07-15")
	if err != nil {
		panic(err)
	}
	day, err := findLunarDay(timeIn)
	if err != nil {
		panic(err)
	}
	fmt.Println(*day)
}

func findLunarDay(timeIn time.Time) (*string, error) {
	fileName := fmt.Sprintf("%d.txt", timeIn.Year())
	data, err := ioutil.ReadFile(path.Join("lunar", fileName))
	if err != nil {
		return nil, err
	}

	// Split the text into lines
	lines := strings.Split(string(data), "\n")

	// Loop through the lines and split each line into columns
	for i, line := range lines {
		// Skip the first three lines
		if i < 3 {
			continue
		}

		// Split the line into fields based on whitespace
		fields := strings.Fields(line)

		// Print the columns
		if len(fields) >= 3 {
			layout := "2006年01月02日"

			// Parse the date string into a time.Time value
			timeParsed, err := time.Parse(layout, fields[0])
			if err != nil {
				return nil, err
			}

			if timeIn.Year() == timeParsed.Year() &&
				timeIn.Month() == timeParsed.Month() &&
				timeIn.Day() == timeParsed.Day() {
				return &fields[1], nil
			}
		}
	}
	return nil, errors.New("Not found")
}

func downloadConvert(year int) error {
	url := fmt.Sprintf("https://www.hko.gov.hk/tc/gts/time/calendar/text/files/T%dc.txt", year)

	// Download the file from the URL
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the contents of the file into a buffer
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Decode the Big5 encoded text
	decoder := traditionalchinese.Big5.NewDecoder()
	bytes, err := decoder.Bytes(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%d.txt", year), bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
