package dsl

import (
	"fmt"
	"math"
	"sort"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func (d *dslConfig) dslHumanizeTime(L *lua.LState) int {
	tt := L.CheckInt64(1)
	then := time.Unix(tt, 0)
	L.Push(lua.LString(humanTime(then)))
	return 1
}

// Seconds-based time units
const (
	humanDay      = 24 * time.Hour
	humanWeek     = 7 * humanDay
	humanMonth    = 30 * humanDay
	humanYear     = 12 * humanMonth
	humanLongTime = 37 * humanYear
)

// Time formats a time into a relative string.
//
// Time(someT) -> "3 weeks ago"
func humanTime(then time.Time) string {
	return humanRelTime(then, time.Now(), "ago", "from now")
}

// A RelTimeMagnitude struct contains a relative time point at which
// the relative format of time will switch to a new format string.  A
// slice of these in ascending order by their "D" field is passed to
// CustomRelTime to format durations.
//
// The Format field is a string that may contain a "%s" which will be
// replaced with the appropriate signed label (e.g. "ago" or "from
// now") and a "%d" that will be replaced by the quantity.
//
// The DivBy field is the amount of time the time difference must be
// divided by in order to display correctly.
//
// e.g. if D is 2*time.Minute and you want to display "%d minutes %s"
// DivBy should be time.Minute so whatever the duration is will be
// expressed in minutes.
type humanRelTimeMagnitude struct {
	D      time.Duration
	Format string
	DivBy  time.Duration
}

var defaultMagnitudes = []humanRelTimeMagnitude{
	{time.Second, "now", time.Second},
	{2 * time.Second, "1 second %s", 1},
	{time.Minute, "%d seconds %s", time.Second},
	{2 * time.Minute, "1 minute %s", 1},
	{time.Hour, "%d minutes %s", time.Minute},
	{2 * time.Hour, "1 hour %s", 1},
	{humanDay, "%d hours %s", time.Hour},
	{2 * humanDay, "1 day %s", 1},
	{humanWeek, "%d days %s", humanDay},
	{2 * humanWeek, "1 week %s", 1},
	{humanMonth, "%d weeks %s", humanWeek},
	{2 * humanMonth, "1 month %s", 1},
	{humanYear, "%d months %s", humanMonth},
	{18 * humanMonth, "1 year %s", 1},
	{2 * humanYear, "2 years %s", 1},
	{humanLongTime, "%d years %s", humanYear},
	{math.MaxInt64, "a long while %s", 1},
}

// RelTime formats a time into a relative string.
//
// It takes two times and two labels.  In addition to the generic time
// delta string (e.g. 5 minutes), the labels are used applied so that
// the label corresponding to the smaller time is applied.
//
// RelTime(timeInPast, timeInFuture, "earlier", "later") -> "3 weeks earlier"
func humanRelTime(a, b time.Time, albl, blbl string) string {
	return humanCustomRelTime(a, b, albl, blbl, defaultMagnitudes)
}

// CustomRelTime formats a time into a relative string.
//
// It takes two times two labels and a table of relative time formats.
// In addition to the generic time delta string (e.g. 5 minutes), the
// labels are used applied so that the label corresponding to the
// smaller time is applied.
func humanCustomRelTime(a, b time.Time, albl, blbl string, magnitudes []humanRelTimeMagnitude) string {
	lbl := albl
	diff := b.Sub(a)

	if a.After(b) {
		lbl = blbl
		diff = a.Sub(b)
	}

	n := sort.Search(len(magnitudes), func(i int) bool {
		return magnitudes[i].D > diff
	})

	if n >= len(magnitudes) {
		n = len(magnitudes) - 1
	}
	mag := magnitudes[n]
	args := []interface{}{}
	escaped := false
	for _, ch := range mag.Format {
		if escaped {
			switch ch {
			case 's':
				args = append(args, lbl)
			case 'd':
				args = append(args, diff/mag.DivBy)
			}
			escaped = false
		} else {
			escaped = ch == '%'
		}
	}
	return fmt.Sprintf(mag.Format, args...)
}
