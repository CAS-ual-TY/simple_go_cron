package simple_cron

import (
	"fmt"
	"testing"
	"time"
)

func test(expected bool, cronS string, time time.Time, t *testing.T) {
	cron, e := ParseCron(cronS)
	if e != nil {
		t.Errorf("Assertion is false! %s =/= %s with error: %s", cronS, time.Format("2006-01-02 15:04:05"), e.Error())
	}
	if !(cron.TimeMatches(time) == expected) {
		t.Errorf("Assertion is false! %s =/= %s", cronS, time.Format("2006-01-02 15:04:05"))
	}
}

func TestNow(t *testing.T) {
	now := time.Now()

	test(true, "* * * * *", now, t)
	test(true, "*", now, t)
	test(true, "", now, t)
	test(true, fmt.Sprintf("%d %d %d %d %d", now.Second(), now.Minute(), now.Hour(), now.Day(), int(now.Month())), now, t)
	test(true, fmt.Sprintf("%d %d", now.Second(), now.Minute()), now, t)
	test(true, fmt.Sprintf("* * * * * %d", now.Weekday()), now, t)
}

func TestOther(t *testing.T) {
	date := time.Date(2010, 12, 30, 23, 59, 30, 0, time.UTC)
	test(true, fmt.Sprintf("%d %d %d %d %d", date.Second(), date.Minute(), date.Hour(), date.Day(), int(date.Month())), date, t)
	test(true, "29-30", date, t)
	test(true, "30-31", date, t)
	test(true, "10,30-31,20", date, t)
}
