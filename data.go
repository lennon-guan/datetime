package datetime

import "time"

type DateTime time.Time

func Parse(datetimeStr, layoutStr string) (DateTime, error) {
	var dt DateTime
	layout, err := newLayout(layoutStr)
	if err != nil {
		return dt, err
	}
	return layout.Parse(datetimeStr)
}

func Now() DateTime {
	return DateTime(time.Now())
}

func New(year, month, day, hour, minute, second, nsec int) DateTime {
	return DateTime(time.Date(year, time.Month(month), day, hour, minute, second, nsec, defaultLocation))
}

func (dt DateTime) Format(layoutStr string) string {
	layout, err := newLayout(layoutStr)
	if err != nil {
		return ""
	}
	return layout.Format(dt)
}

func (dt DateTime) Time() time.Time {
	return time.Time(dt)
}

var defaultLocation *time.Location

func init() {
	defaultLocation = time.Now().Location()
}
