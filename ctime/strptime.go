package ctime

import (
    "time"
    "fmt"
    "strconv"
    "errors"
    "log"
)

type parseState struct {
    format []rune
    value []rune
    indexFormat, indexS int
    year, month, day, hour, minute, second int
}


func Parse(format, value string, loc *time.Location) (t time.Time, err error) {
    state := parseState{format: []rune(format), value: []rune(value), year: 1, month: 1, day: 1}

    defer func () {
        if r := recover(); r != nil {
            fmt.Println("Recoverd parsing")
            err = errors.New(fmt.Sprintf("Parsing failed: %s", r))
        }
    }()

    state.parse()

    return time.Date(state.year, time.Month(state.month), state.day, state.hour,
                     state.minute, state.second, 0, loc), nil
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
    log.Printf("read int of length %d", chars)
    i, err := strconv.Atoi(string(state.read(chars)))
    if err != nil {
        panic("Parse int error")
    }
    log.Printf("int read %d", i)
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
            log.Printf("Got format %c compare to %c", f, c)
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
        state.year = state.readInt(4)
    case 'm':
        state.month = state.readInt(2)
    case 'd':
        state.day = state.readInt(2)
    case 'D':
        state.month = state.readInt(2)
        state.expect("/")
        state.day = state.readInt(2)
        state.expect("/")
        state.year = state.readInt(4)
    case 'H':
        state.hour = state.readInt(2)
    case 'M':
        state.minute = state.readInt(2)
    case 'S':
        state.second = state.readInt(2)
    default:
        panic(fmt.Sprintf("Unsupported format specifier %c. Patches are welcome.", f))
    }
}

