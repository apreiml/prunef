package ctime

import (
    "time"
    "fmt"
)

type parseState struct {
    format []rune
    s []rune
    t *time.Time
}


func Strptime(s, format string, t *time.Time) error {
    state := parseState{[]rune(format), []rune(s), t}

    defer func () {
        if r := recover(); r != nil {
            fmt.Println("Recoverd parsing")
        }
    }()

    return state.parse()
}

func (state *parseState) parse() error {
    fmt.Println("parsing")
    return nil
}

