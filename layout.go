package datetime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	datetime struct {
		year   int
		month  int
		day    int
		hour   int
		minute int
		sec    int
		nsec   int
	}

	formatter interface {
		Format(time.Time) string
		Parse([]rune, *datetime) ([]rune, error)
	}

	layout struct {
		formatters []formatter
	}

	formatter_str string
	formatter_Y   int
	formatter_m   int
	formatter_d   int
	formatter_H   int
	formatter_M   int
	formatter_S   int
	formatter_s   int
)

func (str formatter_str) Format(t time.Time) string {
	return string(str)
}

func (str formatter_str) Parse(s []rune, dt *datetime) ([]rune, error) {
	mylen := len([]rune(string(str)))
	if string(str) != string(s[:mylen]) {
		return nil, fmt.Errorf("string not match, want %s given %s", string(str), string(s[:mylen]))
	}
	return s[mylen:], nil
}

func (formatter_Y) Format(t time.Time) string {
	return strconv.Itoa(t.Year())
}

func (formatter_Y) Parse(s []rune, dt *datetime) ([]rune, error) {
	if len(s) < 4 {
		return nil, errors.New("too short")
	}
	ys := string(s[:4])
	s = s[4:]
	year, err := strconv.Atoi(ys)
	if err != nil {
		return nil, err
	}
	dt.year = year
	return s, nil
}

func (formatter_m) Format(t time.Time) string {
	return fmt.Sprintf("%02d", int(t.Month()))
}

func (formatter_m) Parse(s []rune, dt *datetime) ([]rune, error) {
	if len(s) < 2 {
		return nil, errors.New("too short")
	}
	ms := string(s[:2])
	s = s[2:]
	month, err := strconv.Atoi(ms)
	if err != nil {
		return nil, err
	}
	dt.month = month
	return s, nil
}

func (formatter_d) Format(t time.Time) string {
	return fmt.Sprintf("%02d", t.Day())
}

func (formatter_d) Parse(s []rune, dt *datetime) ([]rune, error) {
	if len(s) < 2 {
		return nil, errors.New("too short")
	}
	ds := string(s[:2])
	s = s[2:]
	day, err := strconv.Atoi(ds)
	if err != nil {
		return nil, err
	}
	dt.day = day
	return s, nil
}

func (formatter_H) Format(t time.Time) string {
	return fmt.Sprintf("%02d", t.Hour())
}

func (formatter_H) Parse(s []rune, dt *datetime) ([]rune, error) {
	if len(s) < 2 {
		return nil, errors.New("too short")
	}
	hs := string(s[:2])
	s = s[2:]
	hour, err := strconv.Atoi(hs)
	if err != nil {
		return nil, err
	}
	dt.hour = hour
	return s, nil
}

func (formatter_M) Format(t time.Time) string {
	return fmt.Sprintf("%02d", t.Minute())
}

func (formatter_M) Parse(s []rune, dt *datetime) ([]rune, error) {
	if len(s) < 2 {
		return nil, errors.New("too short")
	}
	ms := string(s[:2])
	s = s[2:]
	minute, err := strconv.Atoi(ms)
	if err != nil {
		return nil, err
	}
	dt.minute = minute
	return s, nil
}

func (formatter_S) Format(t time.Time) string {
	return fmt.Sprintf("%02d", t.Second())
}

func (formatter_S) Parse(s []rune, dt *datetime) ([]rune, error) {
	if len(s) < 2 {
		return nil, errors.New("too short")
	}
	ss := string(s[:2])
	s = s[2:]
	sec, err := strconv.Atoi(ss)
	if err != nil {
		return nil, err
	}
	dt.sec = sec
	return s, nil
}

func (formatter_s) Format(t time.Time) string {
	return fmt.Sprintf("%03d", t.Nanosecond()/1000000)
}

func (formatter_s) Parse(s []rune, dt *datetime) ([]rune, error) {
	if len(s) < 2 {
		return nil, errors.New("too short")
	}
	ss := string(s[:2])
	s = s[2:]
	msec, err := strconv.Atoi(ss)
	if err != nil {
		return nil, err
	}
	dt.nsec = msec * 1000000
	return s, nil
}

func newLayout(layoutStr string) (*layout, error) {
	special := false
	f := &layout{formatters: []formatter{}}
	chars := []rune(layoutStr)
	normalSeq := make([]rune, 0, len(chars))
	for _, ch := range chars {
		if special {
			var newfmt formatter
			switch ch {
			case rune('Y'):
				newfmt = formatter_Y(0)
			case rune('m'):
				newfmt = formatter_m(0)
			case rune('d'):
				newfmt = formatter_d(0)
			case rune('H'):
				newfmt = formatter_H(0)
			case rune('M'):
				newfmt = formatter_M(0)
			case rune('S'):
				newfmt = formatter_S(0)
			case rune('%'):
				normalSeq = append(normalSeq, ch)
				special = false
			default:
				return nil, errors.New("Invalid layout char " + string(ch))
			}
			if special {
				if len(normalSeq) > 0 {
					f.formatters = append(f.formatters, formatter_str(string(normalSeq)))
					normalSeq = normalSeq[0:0]
				}
				f.formatters = append(f.formatters, newfmt)
				special = false
			}
		} else if ch == rune('%') {
			special = true
			continue
		} else {
			normalSeq = append(normalSeq, ch)
		}
	}
	if len(normalSeq) > 0 {
		f.formatters = append(f.formatters, formatter_str(string(normalSeq)))
	}
	return f, nil
}

func (l *layout) Parse(datetimeStr string) (DateTime, error) {
	s := []rune(datetimeStr)
	var dt datetime
	var err error
	var out DateTime
	for _, f := range l.formatters {
		if s, err = f.Parse(s, &dt); err != nil {
			return out, err
		}
	}
	out = DateTime(time.Date(dt.year, time.Month(dt.month), dt.day, dt.hour, dt.minute, dt.sec, dt.nsec, defaultLocation))
	return out, nil
}

func (l *layout) Format(dt DateTime) string {
	sects := make([]string, len(l.formatters))
	for i, f := range l.formatters {
		sects[i] = f.Format(time.Time(dt))
	}
	return strings.Join(sects, "")
}
