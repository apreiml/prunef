package main

import (
    "bufio"
	"fmt"
    "flag"
	"git.sr.ht/apreiml/prunef/ctime"
	"os"
	"time"
)

var location = time.Local
var format = "thor.server.aktrophy.at-%Y-%m-%d_%H-%M-%S"

var config = struct {
    hourly, daily uint
    utc, printSlots bool
}{}

type slot struct {
    maxTime time.Time
    t *time.Time
    value string
}

type archive struct {
    slots []slot
}

func main() {
    flag.UintVar(&config.hourly, "keep-hourly", 0, "Number of hourly entries to keep.")
    flag.UintVar(&config.daily, "keep-daily", 0, "Number of daily entries to keep.")
    flag.BoolVar(&config.utc, "utc", false, "Parse dates as UTC.")
    flag.BoolVar(&config.printSlots, "print-slots", false, "Print slots and exit.")
    flag.Parse()

    args := flag.Args()
    if len(args) != 1 {
        flag.Usage()
        os.Exit(1)
    }

    if config.utc {
        location = time.UTC
    }

    // formatString := args[0]
    archive := initArchive()
    _ = archive

    if config.printSlots {
        archive.printSlots()
        os.Exit(0)
    }

    scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
        entry := scanner.Text()
        out, err := archive.swapIn(entry)
        if err != nil {
            panic(err)
        }

        if out != "" {
            fmt.Println(out)
        }
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

    archive.printSlots()

	// var val, _ = ctime.Parse(formatString, "thor.server.aktrophy.at-2019-11-02_19-00-03", location)
	//fmt.Println(val)

	//var t = time.Now()
//	fmt.Printf("%#v\n", t)
}

func initArchive() archive {
    var countSlots uint = 1 + config.hourly + config.daily
    var a = archive{slots: make([]slot, countSlots)}
    var t = time.Now().UTC()
    var index int = 1

    a.slots[0].maxTime = time.Time(t)

    t = t.Truncate(time.Hour)
    d, err := time.ParseDuration("1h")
    if err != nil {
        panic("invalid hourly duration")
    }
    for i := uint(0); i < config.hourly; i++ {
        a.slots[index].maxTime = time.Time(t)
        t = t.Add(-d)
        index += 1
    }

    t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
    for i := uint(0); i < config.daily; i++ {
        a.slots[index].maxTime = time.Time(t)
        t = t.AddDate(0, 0, -1)
        index += 1
    }

    return a
}

func (s slot) String() string {
    return fmt.Sprintf("slot: {maxTime: \"%s\", value: \"%s\"", s.maxTime, s.value)
}

func (a archive) printSlots() {
    for _, s := range(a.slots) {
        fmt.Printf("%s\n", s)
    }
}

func (a *archive) swapIn(entry string) (string, error) {
    t, err := ctime.Parse(format, entry, location)
    if err != nil {
        return "", err
    }

    // do not prune entries, that are made while running prunef
    if t.After(a.slots[0].maxTime) {
        return "", nil
    }

    var s *slot = nil
    var swapped = entry

    for i := 0; i < len(a.slots); i++ {
        s = &a.slots[i]
        if t.After(s.maxTime) {
            s = &a.slots[i - 1]
            if s.t == nil || t.After(*s.t) {
                swapped = s.value
                s.t = &t
                s.value = entry
                return swapped, nil
            }
        }
    }

    _ = t
    _ = entry

    return swapped, nil
}
