package main

import (
	"testing"
	"time"
)

func TestTask1GetTime(t *testing.T) {
	gotTime, err := getTime()
	if gotTime.Round(time.Minute) != time.Now().Round(time.Minute) {
		t.Errorf("Time from NTP: %v. Current time: %v mismatch", gotTime.Round(time.Second), time.Now().Round(time.Second))
	}
	if err != nil {
		t.Errorf("Got error while getting time.")
	}
}
