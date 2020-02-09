package core

import (
	"log"
	"time"
)

type (
	clock struct {
		location *time.Location
	}
	// ClockOption は option です
	ClockOption struct {
		Location string
	}
	// Clock は時計です
	Clock interface {
		Now() time.Time
		DefaultLocation() *time.Location
		GetMonthStartDay(tm *time.Time) time.Time
		GetDay(tm *time.Time) time.Time
	}
)

// NewClock is create instance
func NewClock() Clock {
	const location = "Asia/Tokyo"
	return NewClockWithLocation(ClockOption{location})
}

// NewClockWithLocation is create instance with option
func NewClockWithLocation(option ClockOption) Clock {
	var loc, err = time.LoadLocation(option.Location)
	if err != nil {
		log.Fatal(err)
	}
	return &clock{location: loc}
}

func (t *clock) Now() time.Time {
	return time.Now().In(t.location)
}
func (t *clock) DefaultLocation() *time.Location {
	return t.location
}
func (t *clock) GetMonthStartDay(tm *time.Time) time.Time {
	if tm == nil {
		now := t.Now()
		tm = &now
	}
	return time.Date(tm.Year(), tm.Month(), 1, 0, 0, 0, 0, t.location)
}
func (t *clock) GetDay(tm *time.Time) time.Time {
	if tm == nil {
		now := t.Now()
		tm = &now
	}
	return time.Date(tm.Year(), tm.Month(), tm.Day(), 0, 0, 0, 0, t.location)
}
