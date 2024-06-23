package tvdb

import (
	"testing"
)

func TestShiftDate(t *testing.T) {
	// Test case 1: today's date is 2024-06-21
	expected := "2024-06-23"
	result := shiftDate("2024-06-21")
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test case 2: today's date is 2024-06-22
	expected = "2024-06-24"
	result = shiftDate("2024-06-22")
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test case 3: today's date is 2024-06-23
	expected = "2024-06-25"
	result = shiftDate("2024-06-23")
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
