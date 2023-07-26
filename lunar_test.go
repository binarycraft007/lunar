package lunar

import (
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
			t.Log("original gregorian date:", date.Format("2006-01-01"))
			var lunarDate *string
			if lunarDate, err = l.TimeToLunar(date); err != nil {
				t.Fatalf("to lunar: %v", err)
			}
			t.Log("converted lunar date:", *lunarDate)

			date = date.AddDate(0, 0, 1)
		}
	}
}
