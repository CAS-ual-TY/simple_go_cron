package simple_cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CronEntry interface {
	Matches(value int) bool
	GetStart() int
	string() string
}

type CronEntryUniversal struct{}

func (c *CronEntryUniversal) GetStart() int {
	return 0
}

func (c *CronEntryUniversal) Matches(value int) bool {
	return true
}

func (c *CronEntryUniversal) string() string {
	return "*"
}

type CronEntryValue struct {
	Value int
}

func (c *CronEntryValue) Matches(value int) bool {
	return value == c.Value
}

func (c *CronEntryValue) GetStart() int {
	return c.Value
}

func (c *CronEntryValue) string() string {
	return string(rune(c.Value))
}

type CronEntryRange struct {
	Start int
	End   int
}

func (c *CronEntryRange) Matches(value int) bool {
	return value >= c.Start && value <= c.End
}

func (c *CronEntryRange) GetStart() int {
	return c.Start
}

func (c *CronEntryRange) string() string {
	return string(rune(c.Start)) + "-" + string(rune(c.End))
}

type CronEntryStep struct {
	Entry CronEntry
	Step  int
}

func (c *CronEntryStep) Matches(value int) bool {
	return (value-c.Entry.GetStart())%c.Step == 0
}

func (c *CronEntryStep) GetStart() int {
	return c.Entry.GetStart()
}

func (c *CronEntryStep) string() string {
	return c.Entry.string() + "/" + string(rune(c.Step))
}

type CronEntries struct {
	Entries []CronEntry
}

func (c *CronEntries) Matches(value int) bool {
	for _, v := range c.Entries {
		if v.Matches(value) {
			return true
		}
	}
	return false
}

func (c *CronEntries) string() string {
	s := ""
	for i, v := range c.Entries {
		if i > 0 {
			s += ","
		}
		s += v.string()
	}
	return s
}

type Cron struct {
	Seconds  CronEntries
	Minutes  CronEntries
	Hours    CronEntries
	Days     CronEntries
	Months   CronEntries
	Weekdays CronEntries
}

func (c *Cron) Matches(second, minute, hour, day, month, weekday int) bool {
	return c.Seconds.Matches(second) &&
		c.Minutes.Matches(minute) &&
		c.Hours.Matches(hour) &&
		c.Days.Matches(day) &&
		c.Months.Matches(month) &&
		c.Weekdays.Matches(weekday)
}

func (c *Cron) TimeMatches(t time.Time) bool {
	return c.Matches(t.Second(), t.Minute(), t.Hour(), t.Day(), int(t.Month()), int(t.Weekday()))
}

func (c *Cron) String() string {
	return c.Seconds.string() + " " +
		c.Minutes.string() + " " +
		c.Hours.string() + " " +
		c.Days.string() + " " +
		c.Months.string() + " " +
		c.Weekdays.string()
}

type ErrInvalidCronEntry struct {
	ErrorPart string
	ErrorType string
}

func (e *ErrInvalidCronEntry) Error() string {
	return fmt.Sprintf("Invalid %s: %s", e.ErrorType, e.ErrorPart)
}

func ParseCronEntry(s string) (CronEntry, error) {
	if s == "*" {
		return &CronEntryUniversal{}, nil
	} else if strings.ContainsAny(s, "/") {
		parts := strings.Split(s, "/")
		if len(parts) != 2 {
			return nil, &ErrInvalidCronEntry{s, "step"}
		}
		entry, err := ParseCronEntry(parts[1])
		if err != nil {
			return nil, err
		}
		step, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, &ErrInvalidCronEntry{s, "step"}
		}
		return &CronEntryStep{
			Entry: entry,
			Step:  step,
		}, nil
	} else if strings.ContainsAny(s, "-") {
		parts := strings.Split(s, "-")
		if len(parts) != 2 {
			return nil, &ErrInvalidCronEntry{s, "range"}
		}
		start, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, &ErrInvalidCronEntry{s, "range"}
		}
		end, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, &ErrInvalidCronEntry{s, "range"}
		}
		return &CronEntryRange{
			Start: start,
			End:   end,
		}, nil
	} else {
		value, err := strconv.Atoi(s)
		if err != nil {
			return nil, &ErrInvalidCronEntry{s, "value"}
		}
		return &CronEntryValue{Value: value}, nil
	}
}

func ParseCronEntries(s string) (CronEntries, error) {
	entries := strings.Split(s, ",")
	cronEntries := CronEntries{
		Entries: []CronEntry{},
	}
	for _, v := range entries {
		entry, err := ParseCronEntry(v)
		if err != nil {
			return cronEntries, err
		}
		cronEntries.Entries = append(cronEntries.Entries, entry)
	}
	for len(cronEntries.Entries) < 6 {
		cronEntries.Entries = append(cronEntries.Entries, &CronEntryUniversal{})
	}
	return cronEntries, nil
}

func ParseCron(s string) (*Cron, error) {
	entries := make([]CronEntries, 0, 6)

	if len(s) > 0 {
		parts := strings.Split(s, " ")

		if len(parts) > 6 {
			return nil, fmt.Errorf("Invalid cron: %s", s)
		}

		for _, part := range parts {
			entry, err := ParseCronEntries(part)
			if err != nil {
				return nil, err
			}
			entries = append(entries, entry)
		}
	}

	for len(entries) < 6 {
		entries = append(entries, CronEntries{Entries: []CronEntry{&CronEntryUniversal{}}})
	}

	return &Cron{
		Seconds:  entries[0],
		Minutes:  entries[1],
		Hours:    entries[2],
		Days:     entries[3],
		Months:   entries[4],
		Weekdays: entries[5],
	}, nil
}
