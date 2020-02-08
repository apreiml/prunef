package ctime

import (
    "testing"
    "time"
)

func TestSimple(t *testing.T) {
    now := time.Now()

    cases := []struct {
        format, in string
        t *time.Time
    }{
        {"nodate", "nodate", &now},
    }

    for _, c := range cases {
        tt := time.Time(now)
        Strptime(c.format, c.in, &tt)
        if *c.t != tt {
            t.Errorf("Strptime(%s, %s) == %s want %s", c.format, c.in, tt, c.t)
        }
    }
}
