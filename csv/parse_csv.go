// CSV parsing example.
// Data and example problem taken from CodeForAmerica challenge.

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"
)

const (
	longDate  = "2006-01-02 00:00:00"
	shortDate = "2006-01-02"
	openLabel = "[OPEN]"
)

type violation struct {
	violationID   string
	inspectionID  string
	category      string
	opened        time.Time
	closed        time.Time
	isClosed      bool
	violationType string
}

func (v *violation) String() string {
	closedFmt := openLabel
	if v.isClosed {
		closedFmt = v.closed.Format(shortDate)
	}

	return fmt.Sprintf("(%s) %s - %s %q",
		v.violationID,
		v.opened.Format(shortDate),
		closedFmt,
		v.violationType,
	)
}

type summary struct {
	count    int
	earliest *violation
	latest   *violation
}

func (s summary) String() string {
	return fmt.Sprintf("[%d] e: %s, l: %s", s.count, s.earliest, s.latest)
}

func main() {
	var csvFile string
	flag.StringVar(&csvFile, "f", "", "CSV file")
	flag.Parse()

	var in io.Reader = os.Stdin

	if csvFile != "" {
		f, err := os.Open(csvFile)
		if err != nil {
			log.Fatalf("Error opening %q - %s", csvFile, err)
		}
		defer f.Close()

		in = f
	}

	skipped := 0
	byCategory := make(map[string]summary)

	csvSource := csv.NewReader(in)
	for {
		record, err := csvSource.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error parsing CSV - %s", err)
		}

		v, err := extractViolation(record)
		if err != nil {
			skipped++
			continue
		}

		sum := byCategory[v.category]
		sum.count++

		if sum.earliest == nil || v.opened.Before(sum.earliest.opened) {
			sum.earliest = v
		}
		if sum.latest == nil || v.opened.After(sum.latest.opened) {
			sum.latest = v
		}

		byCategory[v.category] = sum
	}

	cats := []string{}
	for c := range byCategory {
		cats = append(cats, c)
	}

	sort.Strings(cats)

	for _, cat := range cats {
		sum := byCategory[cat]
		fmt.Printf("%s (%d)\n", cat, sum.count)
		fmt.Printf("  - %s\n", sum.earliest)
		fmt.Printf("  + %s\n\n", sum.latest)
	}
}

func extractViolation(fields []string) (*violation, error) {
	v := &violation{}
	var err error

	// 0 violation_id
	// 1 inspection_id
	// 2 violation_category
	// 3 violation_date
	// 4 violation_date_closed
	// 5 violation_type

	v.violationID = fields[0]
	v.inspectionID = fields[1]
	v.category = fields[2]

	v.opened, err = time.Parse(longDate, fields[3])
	if err != nil {
		return v, err
	}

	v.isClosed = false
	if fields[4] != "" {
		v.closed, err = time.Parse(longDate, fields[4])
		if err != nil {
			return v, err
		}
		v.isClosed = true
	}

	v.violationType = fields[5]

	return v, nil
}
