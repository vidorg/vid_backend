package vo

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type JsonDate time.Time

type JsonDateTime time.Time

var (
	location, _    = time.LoadLocation("Asia/Shanghai")
	dateFormat     = "2006-01-02"
	dateTimeFormat = "2006-01-02 15:04:05"
)

// json

func (d JsonDate) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", time.Time(d).Format(dateFormat))
	return []byte(str), nil
}

func (dt JsonDateTime) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", time.Time(dt).Format(dateTimeFormat))
	return []byte(str), nil
}

// parse

func (d JsonDate) String() string {
	return time.Time(d).Format(dateFormat)
}

func (dt JsonDateTime) String() string {
	return time.Time(dt).Format(dateTimeFormat)
}

func (d JsonDate) Parse(dateString string, defaultDate JsonDate) JsonDate {
	newD, err := time.ParseInLocation(dateFormat, dateString, location)
	if err != nil {
		return defaultDate
	} else {
		return JsonDate(newD)
	}
}

func (dt JsonDateTime) Parse(dateTimeString string, defaultDateTime JsonDateTime) JsonDateTime {
	newDt, err := time.ParseInLocation(dateTimeFormat, dateTimeString, location)
	if err != nil {
		return defaultDateTime
	} else {
		return JsonDateTime(newDt)
	}
}

// gorm !!

func (d *JsonDate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("wrong format value")
	}
	*d = JsonDate(t)
	return nil
}

func (d JsonDate) Value() (driver.Value, error) {
	return d.String(), nil
}

func (dt *JsonDateTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("wrong format value")
	}
	*dt = JsonDateTime(t)
	return nil
}

func (dt JsonDateTime) Value() (driver.Value, error) {
	return dt.String(), nil
}
