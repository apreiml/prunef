// ctime provides time utils that work with posix format specifiers
// like used in strftime(3)
package ctime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type tm struct {
	year, month, day, hour, minute, second int
	loc                                    *time.Location
}

type parseState struct {
	format              []rune
	value               []rune
	indexFormat, indexS int
	readS               int
	time                tm
}

func Parse(format, value string, loc *time.Location) (t time.Time, err error) {
	tm := tm{year: 1, month: 1, day: 1, loc: loc}
	state := parseState{format: []rune(format), value: []rune(value), time: tm}

	defer func() {
		if r := recover(); r != nil {
			indexMarker := state.sprintIndexMarker()
			msg := fmt.Sprintf("%s\n%s", r, indexMarker)
			err = errors.New(msg)
		}
	}()

	state.parse()
	return state.time.toTime(), nil
}

func (state parseState) sprintIndexMarker() string {
	prefix := "at: "
	indexMarkerOffset := state.indexS + len(prefix)
	return fmt.Sprintf("\t%s%s\n\t%s^", prefix, string(state.value),
		strings.Repeat(" ", indexMarkerOffset))
}

func (t *tm) toTime() time.Time {
	return time.Date(t.year, time.Month(t.month), t.day, t.hour,
		t.minute, t.second, 0, t.loc)
}

func (state *parseState) isEndOfFormat() bool {
	return state.indexFormat >= len(state.format)
}

func (state *parseState) readFormat() rune {
	f := state.format[state.indexFormat]
	state.indexFormat += 1
	return f
}

func (state *parseState) read(chars int) []rune {
	begin := state.indexS
	end := begin + chars

	state.readS = end
	return state.value[begin:end]
}

func (state *parseState) next() {
	state.indexS = state.readS
}

func (state *parseState) expectNext(expected string) {
	got := string(state.read(len(expected)))
	if got != expected {
		panic(fmt.Sprintf("Expected %s", expected))
	}
	state.next()
}

func (state *parseState) readInt(length int) int {
	s := string(state.read(length))
	i, err := strconv.Atoi(s)
	if err != nil {
		panic("Expected integer")
	}
	return i
}


func (state *parseState) nextInt(length int) int {
	i := state.readInt(length)
	state.next()
	return i
}

func (state *parseState) nextIntBetween(length, min, max int) int {
	i := state.readInt(length)
	if (i < min || i > max) {
		panic(fmt.Sprintf("Expected integer between %d and %d", min, max))
	}
	state.next()
	return i
}

func (state *parseState) parse() {
	for {
		if state.isEndOfFormat() {
			break
		}
		f := state.readFormat()

		if f == '%' {
			state.parseFormat()
		} else {
			state.expectNext(string(f))
		}
	}
}

func (state *parseState) parseFormat() {
	switch f := state.readFormat(); f {
	case '%':
		state.expectNext("%")
	case 'Y':
		state.time.year = state.nextIntBetween(4, 0000, 9999)
	case 'm':
		state.time.month = state.nextIntBetween(2, 1, 12)
	case 'd':
		state.time.day = state.nextIntBetween(2, 1, 31)
	case 'D':
		state.time.month = state.nextIntBetween(2, 1, 12)
		state.expectNext("/")
		state.time.day = state.nextIntBetween(2, 1, 31)
		state.expectNext("/")
		year := state.nextIntBetween(2, 00, 99)
		if year >= 69 {
			state.time.year = 1900 + year
		} else {
			state.time.year = 2000 + year
		}
	case 'H':
		state.time.hour = state.nextIntBetween(2, 0, 24)
	case 'M':
		state.time.minute = state.nextIntBetween(2, 0, 59)
	case 'S':
		state.time.second = state.nextIntBetween(2, 0, 59)
	default:
		panic(fmt.Sprintf("Unsupported format specifier %%%c. Patches are welcome.", f))
	}
}
