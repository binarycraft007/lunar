package main

import (
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

	var months []string
	var month string
	var day string
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

			if strings.HasSuffix(fields[1], "月") {
				months = append(months, fields[1])
			}

			// Parse the date string into a time.Time value
			timeParsed, err := time.Parse(layout, fields[0])
			if err != nil {
				return nil, err
			}

			if timeIn.Year() == timeParsed.Year() &&
				timeIn.Month() == timeParsed.Month() &&
				timeIn.Day() == timeParsed.Day() {
				if strings.HasSuffix(fields[1], "月") {
					day = "初一"
				} else {
					day = fields[1]
				}

				if len(months) > 0 {
					month = months[len(months)-1]
					break
				} else {
					continue
				}
			}

			if len(day) > 0 && len(months) > 0 {
				// TODO fix this, should minus one
				month = monthConvert(months[len(months)-1])
				break
			}
		}
	}

	result := month + day
	return &result, nil
}

func monthConvert(month string) string {
	switch month {
	case "一月":
		return "正月"
	case "二月":
		return month
	case "三月":
		return month
	case "四月":
		return month
	case "五月":
		return month
	case "六月":
		return month
	case "七月":
		return month
	case "八月":
		return month
	case "九月":
		return month
	case "十月":
		return month
	case "十一月":
		return "冬月"
	case "十二月":
		return "腊月"
	}
	return month
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
