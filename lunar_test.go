package lunar

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestAllLookup(t *testing.T) {
	l, err := NewLunar()
	if err != nil {
		t.Fatalf("init lunar: %v", err)
	}

	for name, _ := range l.LookupMap {
		yearStr := strings.TrimSuffix(name, ".txt")
		year, err := strconv.Atoi(yearStr)
		date := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		for date.Year() == year {
			t.Log("original gregorian date:", date.Format("2006-01-02"))
			var lunarDate *string
			if lunarDate, err = l.TimeToLunar(date); err != nil {
				t.Fatalf("to lunar: %v", err)
			}

			if len(*lunarDate) != 27 {
				t.Log("diagnostics dump:")
				t.Log("    number of bytes:", len(*lunarDate))
				t.Log("    converted lunar date:", *lunarDate)
				t.Fatalf("convert lunar: %s", date.Format("2006-01-02"))
			}

			timeParsed, err := parseLunarToTime(*lunarDate)
			if err != nil {
				t.Fatalf("parse lunar: %v", err)
			}

			// Calculate the time gap between the two times
			gap := date.Sub(*timeParsed)

			// Check if the gap is greater than a year
			if gap > time.Hour*24*365 {
				t.Fatalf(
					"gap bigger than one year, %s %s",
					date.Format("2006-01-02"),
					timeParsed.Format("2006-01-02"),
				)
			}

			t.Log("converted lunar date:", *lunarDate)

			date = date.AddDate(0, 0, 1)
		}
	}
}

func yearAliasToNum(s string) (int, error) {
	for i, replace := range numberAlias {
		s = strings.Replace(s, replace, fmt.Sprintf("%d", i), -1)
	}
	return strconv.Atoi(s)
}

func parseLunarToTime(l string) (*time.Time, error) {
	indexYear := strings.IndexAny(l, "年")
	indexMonth := strings.IndexAny(l, "月")

	yearStr := l[0:indexYear]
	monthStr := l[indexYear+3 : indexMonth+3]
	dayStr := l[indexMonth+3 : len(l)]

	year, err := yearAliasToNum(yearStr)
	if err != nil {
		return nil, err
	}

	var month int
	switch monthStr {
	case "正月":
		month = 1
	case "二月":
		month = 2
	case "三月":
		month = 3
	case "四月":
		month = 4
	case "五月":
		month = 5
	case "六月":
		month = 6
	case "七月":
		month = 7
	case "八月":
		month = 8
	case "九月":
		month = 9
	case "十月":
		month = 10
	case "冬月":
		month = 11
	case "腊月":
		month = 12
	default:
		month = 0
	}

	var dayStrNum string
	switch dayStr[:3] {
	case "初":
		dayStrNum = "0"
	case "十":
		dayStrNum = "1"
	case "二":
		dayStrNum = "2"
	case "廿":
		dayStrNum = "2"
	case "三":
		dayStrNum = "3"
	default:
		dayStrNum = ""
	}

	switch strings.TrimPrefix(dayStr, dayStr[:3]) {
	case "十":
		if dayStrNum == "0" {
			dayStrNum = "1" + "0"
		} else {
			dayStrNum = dayStrNum + "0"
		}
	case "一":
		dayStrNum = dayStrNum + "1"
	case "二":
		dayStrNum = dayStrNum + "2"
	case "三":
		dayStrNum = dayStrNum + "3"
	case "四":
		dayStrNum = dayStrNum + "4"
	case "五":
		dayStrNum = dayStrNum + "5"
	case "六":
		dayStrNum = dayStrNum + "6"
	case "七":
		dayStrNum = dayStrNum + "7"
	case "八":
		dayStrNum = dayStrNum + "8"
	case "九":
		dayStrNum = dayStrNum + "9"
	default:
		dayStrNum = ""
	}

	day, err := strconv.Atoi(dayStrNum)
	if err != nil {
		return nil, err
	}

	if month == 2 && day > 28 {
		// lunar february may have days greater than 28,
		// won't got parsed by time.Parse, use workaround
		day = 28
	}

	dateStr := fmt.Sprintf("%d-%02d-%02d", year, month, day)

	timeParsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	return &timeParsed, nil
}
