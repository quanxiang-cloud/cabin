package time

import (
	"testing"
)

func TestTolerant(t *testing.T) {
	var tests = []struct {
		in       string
		expected UTC
	}{
		{
			in:       "UTC-0",
			expected: UTC0,
		},
		{
			in:       "UTC-1",
			expected: UTC_1,
		},
		{
			in:       "UTC-2",
			expected: UTC_2,
		},
		{
			in:       "UTC-3",
			expected: UTC_3,
		},
		{
			in:       "UTC-4",
			expected: UTC_4,
		},
		{
			in:       "UTC-5",
			expected: UTC_5,
		},
		{
			in:       "UTC-6",
			expected: UTC_6,
		},
		{
			in:       "UTC-7",
			expected: UTC_7,
		},
		{
			in:       "UTC-8",
			expected: UTC_8,
		},
		{
			in:       "UTC-9",
			expected: UTC_9,
		},
		{
			in:       "UTC-10",
			expected: UTC_10,
		},
		{
			in:       "UTC-11",
			expected: UTC_11,
		},
		{
			in:       "UTC-12",
			expected: UTC_12,
		},
		{
			in:       "UTC+1",
			expected: UTC1,
		},
		{
			in:       "UTC+2",
			expected: UTC2,
		},
		{
			in:       "UTC+3",
			expected: UTC3,
		},
		{
			in:       "UTC+4",
			expected: UTC4,
		},
		{
			in:       "UTC+5",
			expected: UTC5,
		},
		{
			in:       "UTC+6",
			expected: UTC6,
		},
		{
			in:       "UTC+7",
			expected: UTC7,
		},
		{
			in:       "UTC+8",
			expected: UTC8,
		},
		{
			in:       "UTC+9",
			expected: UTC9,
		},
		{
			in:       "UTC+10",
			expected: UTC10,
		},
		{
			in:       "UTC+11",
			expected: UTC11,
		},
		{
			in:       "UTC+12",
			expected: UTC12,
		},
	}

	for _, tt := range tests {
		actual, err := Tolerant(tt.in)
		if err != nil {
			t.Errorf(err.Error())
		}
		if actual != tt.expected {
			t.Errorf("Fib(%s) = %d; expected %d", tt.in, actual, tt.expected)
		}
	}
}

func TestConversion(t *testing.T) {
	expected := NowUnix()

	timeStr := Format(expected)

	actual, err := Unix(timeStr)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if actual != expected {
		t.Errorf("Fib(%s) = %d; expected %d", timeStr, actual, expected)
	}
}

func TestRevers(t *testing.T) {
	expected := Now()

	interim, err := Revise(expected, UTC8)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	actual, err := Regular(interim, UTC8)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if actual != expected {
		t.Errorf("Fib(%s) = %s; expected %s", expected, actual, expected)
	}
}
