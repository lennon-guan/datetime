package datetime

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	dt := New(2017, 1, 4, 15, 24, 59, 0)
	if r := dt.Format("%Y"); r != "2017" {
		t.Error("format %Y error", r)
	}
	if r := dt.Format("%Y年"); r != "2017年" {
		t.Error("format %Y年 error", r)
	}
	if r := dt.Format("%Y-%m-%d"); r != "2017-01-04" {
		t.Error("format %Y-%m-%d error", r)
	}
	if r := dt.Format("%Y-%m-%d %H:%M:%S"); r != "2017-01-04 15:24:59" {
		t.Error("format %Y-%m-%d %H:%M:%S error", r)
	}
}

func TestParse(t *testing.T) {
	ds := "2017-01-04 15:24:59"
	tt, _ := time.Parse("2006-01-02 15:04:05 -0700", ds+" +0800")
	if dt, err := Parse(ds, "%Y-%m-%d %H:%M:%S"); err != nil {
		t.Error("parse error", err.Error())
	} else if dt.Time() != tt {
		t.Error("parse error", dt.Time().String(), tt.String())
	}
}
