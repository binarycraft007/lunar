package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"golang.org/x/text/encoding/traditionalchinese"
)

//go:embed lunar/data.tar.gz
var b []byte
var lookupMap map[string][]byte
var numberAlias = [...]string{
	"零", "一", "二", "三", "四",
	"五", "六", "七", "八", "九",
}

func yearAlias(year int) string {
	s := fmt.Sprintf("%d", year)
	for i, replace := range numberAlias {
		s = strings.Replace(s, fmt.Sprintf("%d", i), replace, -1)
	}
	return s
}

func main() {
	var err error
	if lookupMap, err = lookupTablesToMap(b); err != nil {
		panic(err)
	}
	timeIn, err := time.Parse("2006-01-02", "1901-01-20")
	if err != nil {
		panic(err)
	}
	day, err := findLunarDay(timeIn)
	if err != nil {
		panic(err)
	}
	fmt.Println(*day)
}

func lookupTablesToMap(b []byte) (map[string][]byte, error) {
	// Create a gzip reader
	gzipReader, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer gzipReader.Close()

	// Create a tar reader
	tarReader := tar.NewReader(gzipReader)

	// Create a map to store the file contents
	files := make(map[string][]byte)

	// Iterate over the files in the tar archive
	for {
		header, err := tarReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Read the file contents into a buffer
		buffer := new(bytes.Buffer)
		_, err = io.Copy(buffer, tarReader)
		if err != nil {
			return nil, err
		}

		// Add the file contents to the map
		files[path.Base(header.Name)] = buffer.Bytes()
	}
	return files, nil
}

func findLunarDay(timeIn time.Time) (*string, error) {
	fileName := fmt.Sprintf("%d.txt", timeIn.Year())
	data, ok := lookupMap[fileName]
	if !ok {
		return nil, errors.New("Year not found")
	}

	// Split the text into lines
	lines := strings.Split(string(data), "\n")

	var months []string
	var year string
	var month string
	var day string
	layout := "2006年01月02日"
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
					if months[len(months)-1] == "十二月" {
						year = yearAlias(timeIn.Year() - 1)
					} else {
						year = yearAlias(timeIn.Year())
					}
					month = monthConvert(months[len(months)-1])
					break
				} else {
					continue
				}
			}

			if len(day) > 0 && len(months) > 0 {
				if months[0] == "一月" || months[0] == "十二月" {
					year = yearAlias(timeIn.Year() - 1)
				} else {
					year = yearAlias(timeIn.Year())
				}
				month = lastMonthConvert(months[len(months)-1])
				break
			}
		}
	}

	result := year + "年" + month + day
	return &result, nil
}

func monthConvert(month string) string {
	switch month {
	case "一月":
		return "正月"
	case "十一月":
		return "冬月"
	case "十二月":
		return "腊月"
	default:
		return month
	}
}

func lastMonthConvert(month string) string {
	switch month {
	case "一月":
		return "腊月"
	case "二月":
		return "一月"
	case "三月":
		return "二月"
	case "四月":
		return "三月"
	case "五月":
		return "四月"
	case "六月":
		return "五月"
	case "七月":
		return "六月"
	case "八月":
		return "七月"
	case "九月":
		return "八月"
	case "十月":
		return "九月"
	case "十一月":
		return "十月"
	case "十二月":
		return "冬月"
	default:
		return ""
	}
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

	err = ioutil.WriteFile(path.Join("lunar", fmt.Sprintf("%d.txt", year)), bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
