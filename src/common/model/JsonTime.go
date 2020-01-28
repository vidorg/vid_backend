package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type JsonDate time.Time

type JsonDateTime time.Time

var (
	CurrLocation, _ = time.LoadLocation("Asia/Shanghai")
	DateFormat      = "2006-01-02"
	DateTimeFormat  = "2006-01-02 15:04:05"
)

// json

func (d JsonDate) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", time.Time(d).Format(DateFormat))
	return []byte(str), nil
}

func (dt JsonDateTime) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("\"%s\"", time.Time(dt).Format(DateTimeFormat))
	return []byte(str), nil
}

// date -> str

func (d JsonDate) String() string {
	return time.Time(d).Format(DateFormat)
}

func (dt JsonDateTime) String() string {
	return time.Time(dt).Format(DateTimeFormat)
}

// str -> date

func (d JsonDate) Parse(dateString string) (JsonDate, error) {
	newD, err := time.ParseInLocation(DateFormat, dateString, CurrLocation)
	return JsonDate(newD), err
}

func (dt JsonDateTime) Parse(dateTimeString string) (JsonDateTime, error) {
	newDt, err := time.ParseInLocation(DateTimeFormat, dateTimeString, CurrLocation)
	return JsonDateTime(newDt), err
}

func (d JsonDate) ParseDefault(dateString string, defaultDate JsonDate) JsonDate {
	newD, err := time.ParseInLocation(DateFormat, dateString, CurrLocation)
	if err != nil {
		return JsonDate(newD)
	} else {
		return defaultDate
	}
}

func (dt JsonDateTime) ParseDefault(dateTimeString string, defaultDateTime JsonDateTime) JsonDateTime {
	newDt, err := time.ParseInLocation(DateTimeFormat, dateTimeString, CurrLocation)
	if err != nil {
		return JsonDateTime(newDt)
	} else {
		return defaultDateTime
	}
}

// gorm

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
