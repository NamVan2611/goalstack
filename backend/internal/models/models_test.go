package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestCreateGoalRequestAcceptsDateOnlyStartDate(t *testing.T) {
	var req CreateGoalRequest

	err := json.Unmarshal([]byte(`{
		"title": "Launch sprint",
		"startDate": "2026-06-09",
		"totalDuration": 7,
		"durationType": "days"
	}`), &req)
	if err != nil {
		t.Fatalf("expected date-only startDate to parse: %v", err)
	}

	if got := req.StartDate.Format("2006-01-02"); got != "2026-06-09" {
		t.Fatalf("expected parsed date 2026-06-09, got %s", got)
	}
}

func TestCreateGoalRequestAcceptsRFC3339StartDate(t *testing.T) {
	var req CreateGoalRequest

	err := json.Unmarshal([]byte(`{
		"title": "Launch sprint",
		"startDate": "2026-06-09T00:00:00Z",
		"totalDuration": 7,
		"durationType": "days"
	}`), &req)
	if err != nil {
		t.Fatalf("expected RFC3339 startDate to parse: %v", err)
	}

	want := time.Date(2026, 6, 9, 0, 0, 0, 0, time.UTC)
	if !req.StartDate.Equal(want) {
		t.Fatalf("expected parsed date %s, got %s", want, req.StartDate.Time)
	}
}
