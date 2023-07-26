package lunar

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestFirstDayLookup(t *testing.T) {
	l, err := NewLunar()
	if err != nil {
		t.Fatalf("init lunar: %v", err)
	}

	for name, _ := range l.LookupMap {
		year := strings.TrimSuffix(name, ".txt")
		dateStr := fmt.Sprintf("%s-01-01", year)
		timeIn, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			t.Fatalf("parse time: %v", err)
		}
		if _, err = l.TimeToLunar(timeIn); err != nil {
			t.Fatalf("to lunar: %v", err)
		}
	}
}
