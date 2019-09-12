package main

import (
	"testing"
)

func TestNewReport(t *testing.T) {
	report := NewReport("deadBEAF")
	if report == nil {
		t.Fatalf("failed to get a new report")
	}
	if report.Title != "deadBEAF" {
		t.Errorf("report constructed with the wrong title")
	}
}
