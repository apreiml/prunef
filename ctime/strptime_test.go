package ctime

import (
    "testing"
    "time"
)
func TestParseFail(t *testing.T) {
    cases := []struct {
        format, in string
    }{
        {"asdf", "jkld"},
        {"%", "%"},
    }

    for _, c := range cases {
        _, err := Parse(c.format, c.in, time.UTC)
        if err == nil {
            t.Errorf("Strptime(%s, %s) should fail", c.format, c.in)
        }
    }
}

func TestSimple(t *testing.T) {
    cases := []struct {
        format, in string
        expected string
    }{
        {"nodate", "nodate", "0001-01-01T00:00:00Z"},
        {"%%", "%", "0001-01-01T00:00:00Z"},
        {"%Y", "2020", "2020-01-01T00:00:00Z"},
    }

    for _, c := range cases {
        result, err := Parse(c.format, c.in, time.UTC)
        if err != nil {
            t.Errorf("failed Parsing: %s", err)
        }

        expectedTime, _ := time.Parse(time.RFC3339, c.expected)
        if expectedTime != result {
            t.Errorf("Strptime(%s, %s) == %s want %s", c.format, c.in, result, c.expected)
        }
    }
}
