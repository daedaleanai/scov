package main

import (
	"testing"
)

func TestCoverage_Rating(t *testing.T) {
	cases := []struct {
		hits     int
		count    int
		expected CoverageRating
		str      string
	}{
		{0, 100, LowCoverage, "low"},
		{50, 100, LowCoverage, "low"},
		{74, 100, LowCoverage, "low"},
		{75, 100, MediumCoverage, "medium"},
		{90, 100, HighCoverage, "high"},
	}

	for i, v := range cases {
		cov := Coverage{v.hits, v.count}
		if out := cov.Rating(); out != v.expected {
			t.Errorf("Case %d: mismatch, expected %v, got %v", i, v.expected, out)
		}
		if out := cov.Rating().String(); out != v.str {
			t.Errorf("Case %d: mismatch, expected %v, got %v", i, v.str, out)
		}
	}
}
