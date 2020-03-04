// ctime provides time utils that work with posix format specifiers
// like used in strftime(3)
package ctime

import (
	"errors"
	"fmt"
	"strconv"
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
	time                tm
}

func Parse(format, value string, loc *time.Location) (t time.Time, err error) {
	tm := tm{year: 1, month: 1, day: 1, loc: loc}
	state := parseState{format: []rune(format), value: []rune(value), time: tm}

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Parsing date failed: %s", r))
		}
	}()

	state.parse()
	return state.time.toTime(), nil
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

	state.indexS = end
	return state.value[begin:end]
}

func (state *parseState) expect(expected string) {
	got := string(state.read(len(expected)))
	if got != expected {
		panic(fmt.Sprintf("Didn't expect %s", got))
	}
}

func (state *parseState) readInt(chars int) int {
	s := string(state.read(chars))
	i, err := strconv.Atoi(s)
	if err != nil {
		panic("Parse int error")
	}
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
			c := state.read(1)[0]
			if c != f {
				panic(fmt.Sprintf("Unexpected char %c, expected: %c", c, f))
			}
		}
	}
}

func (state *parseState) parseFormat() {
	switch f := state.readFormat(); f {
	case '%':
		state.expect("%")
	case 'Y':
		state.time.year = state.readInt(4)
	case 'm':
		state.time.month = state.readInt(2)
	case 'd':
		state.time.day = state.readInt(2)
	case 'D':
		state.time.month = state.readInt(2)
		state.expect("/")
		state.time.day = state.readInt(2)
		state.expect("/")
		year := state.readInt(2)
		if year >= 69 {
			state.time.year = 1900 + year
		} else {
			state.time.year = 2000 + year
		}
	case 'H':
		state.time.hour = state.readInt(2)
	case 'M':
		state.time.minute = state.readInt(2)
	case 'S':
		state.time.second = state.readInt(2)
	default:
		panic(fmt.Sprintf("Unsupported format specifier %c. Patches are welcome.", f))
	}
}
