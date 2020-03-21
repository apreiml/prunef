package main

import (
	"bufio"
	"flag"
	"fmt"
	"git.sr.ht/~apreiml/prunef/ctime"
	"os"
	"path/filepath"
	"time"
)

var config = struct {
	secondly, minutely, hourly, daily, weekly, monthly, yearly uint
	utc, printSlots, invert                                    bool
	format                                                     string
	location                                                   *time.Location
}{}

type slot struct {
	maxTime time.Time
	t       *time.Time
	value   string
}

type archive struct {
	slots    []slot
	numSlots uint
}

func main() {
	parseArgs()
	archive := initArchive()

	if config.printSlots && len(config.format) == 0 {
		archive.printSlots()
		os.Exit(0)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		entry := scanner.Text()
		out, err := archive.swapIn(entry)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if out != "" && !config.invert {
			fmt.Println(out)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input.", err)
		os.Exit(1)
	}

	if config.printSlots {
		archive.printSlots()
	} else if config.invert {
		archive.printValues()
	}
}

func parseArgs() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] FORMAT\n", filepath.Base(os.Args[0]))
		fmt.Println(`
  --list-kept
        Show entries to keep instead of entries to prune.
  --keep-daily uint
        Number of daily entries to keep.
  --keep-hourly uint
        Number of hourly entries to keep.
  --keep-minutely uint
        Number of minutely entries to keep.
  --keep-monthly uint
        Number of monthly entries to keep.
  --keep-secondly uint
        Number of secondly entries to keep.
  --keep-weekly uint
        Number of weekly entries to keep.
  --keep-yearly uint
        Number of yearly entries to keep.
  --print-slots
        Print slots and exit.
  --utc
        Expect entry dates in UTC.

FORMAT can have following specifiers:

  %%    a literal %
  %Y    year
  %m    month (01..12)
  %d    day of month (01..)
  %H    hour (00..23)
  %M    minute (00..59)
  %S    second (00..59)
  %D    short for %m/%d/%y
`)
	}

	flag.UintVar(&config.secondly, "keep-secondly", 0, "Number of secondly entries to keep.")
	flag.UintVar(&config.minutely, "keep-minutely", 0, "Number of minutely entries to keep.")
	flag.UintVar(&config.hourly, "keep-hourly", 0, "Number of hourly entries to keep.")
	flag.UintVar(&config.daily, "keep-daily", 0, "Number of daily entries to keep.")
	flag.UintVar(&config.weekly, "keep-weekly", 0, "Number of weekly entries to keep.")
	flag.UintVar(&config.monthly, "keep-monthly", 0, "Number of monthly entries to keep.")
	flag.UintVar(&config.yearly, "keep-yearly", 0, "Number of yearly entries to keep.")
	flag.BoolVar(&config.utc, "utc", false, "Expect entry dates in UTC.")
	flag.BoolVar(&config.invert, "invert", false, "Show entries to keep instead of entries to prune.")
	flag.BoolVar(&config.invert, "list-kept", false, "Show entries to keep instead of entries to prune. Replaces --invert")
	flag.BoolVar(&config.printSlots, "print-slots", false, "Print slots and exit.")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1  {
		if !config.printSlots {
			flag.Usage()
			os.Exit(1)
		}
	} else {
		config.format = args[0]
		// if there are args and print-slots is enabled, use invert
		// to hide pruned output so that the slots can be printed
		// with assigned values at the end of the program
		if config.printSlots {
			config.invert = true
		}
	}

	if config.utc {
		config.location = time.UTC
	} else {
		config.location = time.Local
	}

}

type timeAdder func(time.Time) time.Time

func initArchive() archive {
	var countSlots uint = 1 + config.secondly + config.minutely + config.hourly +
		config.daily + config.weekly + config.monthly + config.yearly

	var a = archive{slots: make([]slot, countSlots)}
	var t = time.Now().UTC()

	a.slots[0].maxTime = time.Time(t)
	a.numSlots = 1

	t = t.Truncate(time.Second)
	t = a.makeSlots(config.secondly, t, makeDurationTimeAdder("1s"))

	t = t.Truncate(time.Minute)
	t = a.makeSlots(config.minutely, t, makeDurationTimeAdder("1m"))

	t = t.Truncate(time.Hour)
	t = a.makeSlots(config.hourly, t, makeDurationTimeAdder("1h"))

	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	t = a.makeSlots(config.daily, t, func(t time.Time) time.Time {
		return t.AddDate(0, 0, -1)
	})

	t = endOfPreviousWeek(t)
	t = a.makeSlots(config.weekly, t, func(t time.Time) time.Time {
		return t.AddDate(0, 0, -7)
	})

	t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	t = a.makeSlots(config.monthly, t, func(t time.Time) time.Time {
		return t.AddDate(0, -1, 0)
	})

	t = time.Date(t.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	t = a.makeSlots(config.yearly, t, func(t time.Time) time.Time {
		return t.AddDate(-1, 0, 0)
	})

	return a
}

func (a *archive) makeSlots(amount uint, startTime time.Time, adder timeAdder) time.Time {
	var t = startTime
	var index = a.numSlots
	for i := index; i < index+amount; i++ {
		a.slots[i].maxTime = time.Time(t)
		a.numSlots++
		t = adder(t)
	}
	return t
}

func makeDurationTimeAdder(duration string) timeAdder {
	d, err := time.ParseDuration(duration)
	if err != nil {
		panic("Implementation Fail: Invalid duration: " + duration)
	}

	return func(t time.Time) time.Time {
		return t.Add(-d)
	}
}

func endOfPreviousWeek(t time.Time) time.Time {
	return t.AddDate(0, 0, -int(t.Weekday()))
}

func (a archive) printSlots() {
	for _, s := range a.slots {
		fmt.Printf("%s", s.maxTime)
		if len(s.value) > 0 {
			fmt.Printf("\t%s\t%s", s.t.UTC(), s.value)
		}
		fmt.Println()
	}
}

func (a archive) printValues() {
	for _, s := range a.slots {
		if s.t != nil {
			fmt.Println(s.value)
		}
	}
}

// swapIn tries to find a slot for given entry. If found the entry of found
// slot will be returned to be pruned. Otherwise the given entry will be
// returned.
func (a *archive) swapIn(entry string) (string, error) {
	t, err := ctime.Parse(config.format, entry, config.location)
	if err != nil {
		return "", fmt.Errorf("Error parsing '%s': %s", entry, err.Error())
	}

	// do not prune entries, that are made while running prunef
	if t.After(a.slots[0].maxTime) {
		if config.invert {
			fmt.Println(entry)
		}
		return "", nil
	}

	var s, next *slot
	var swappedOut = entry

	// find slot and check if given entry is newer than the entry of found slot
	for i := uint(0); i < a.numSlots; i++ {
		s = &a.slots[i]
		if i < a.numSlots-1 {
			next = &a.slots[i+1]
		} else {
			next = nil
		}
		if next == nil || t.After(next.maxTime) {
			if s.t == nil || t.After(*s.t) {
				swappedOut = s.value
				s.t = &t
				s.value = entry
			}
			return swappedOut, nil
		}
	}

	return swappedOut, nil
}
