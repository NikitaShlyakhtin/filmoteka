package validator

import (
	"regexp"
	"testing"
)

func TestNew(t *testing.T) {
	v := New()

	if v.Errors == nil {
		t.Error("Expected Errors map to be initialized, but got nil")
	}

	if len(v.Errors) != 0 {
		t.Errorf("Expected Errors map to be empty, but got %d errors", len(v.Errors))
	}
}

func TestValidator_Valid(t *testing.T) {
	v := New()

	if !v.Valid() {
		t.Error("Expected Valid() to return true for empty errors map")
	}

	v.Errors["field1"] = "error1"
	if v.Valid() {
		t.Error("Expected Valid() to return false for non-empty errors map")
	}
}

func TestValidator_AddError(t *testing.T) {
	v := New()

	v.AddError("field1", "error1")
	if len(v.Errors) != 1 {
		t.Errorf("Expected Errors map to have 1 error, but got %d errors", len(v.Errors))
	}
	if v.Errors["field1"] != "error1" {
		t.Errorf("Expected error message for field1 to be 'error1', but got '%s'", v.Errors["field1"])
	}

	v.AddError("field1", "error2")
	if len(v.Errors) != 1 {
		t.Errorf("Expected Errors map to still have 1 error, but got %d errors", len(v.Errors))
	}
	if v.Errors["field1"] != "error1" {
		t.Errorf("Expected error message for field1 to still be 'error1', but got '%s'", v.Errors["field1"])
	}

	v.AddError("field2", "error3")
	if len(v.Errors) != 2 {
		t.Errorf("Expected Errors map to have 2 errors, but got %d errors", len(v.Errors))
	}
	if v.Errors["field2"] != "error3" {
		t.Errorf("Expected error message for field2 to be 'error3', but got '%s'", v.Errors["field2"])
	}
}

func TestValidator_Check(t *testing.T) {
	v := New()

	// Test case 1: Check with ok=true
	v.Check(true, "field1", "error1")
	if len(v.Errors) != 0 {
		t.Errorf("Expected Errors map to be empty, but got %d errors", len(v.Errors))
	}

	// Test case 2: Check with ok=false
	v.Check(false, "field2", "error2")
	if len(v.Errors) != 1 {
		t.Errorf("Expected Errors map to have 1 error, but got %d errors", len(v.Errors))
	}
	if v.Errors["field2"] != "error2" {
		t.Errorf("Expected error message for field2 to be 'error2', but got '%s'", v.Errors["field2"])
	}
}

func TestIn(t *testing.T) {
	tests := []struct {
		value string
		list  []string
		want  bool
	}{
		{"apple", []string{"apple", "banana", "orange"}, true},
		{"banana", []string{"apple", "banana", "orange"}, true},
		{"orange", []string{"apple", "banana", "orange"}, true},
		{"grape", []string{"apple", "banana", "orange"}, false},
		{"", []string{"apple", "banana", "orange"}, false},
		{"apple", []string{}, false},
	}

	for _, tt := range tests {
		got := In(tt.value, tt.list...)
		if got != tt.want {
			t.Errorf("In(%q, %q...) = %v, want %v", tt.value, tt.list, got, tt.want)
		}
	}
}

func TestMatches(t *testing.T) {
	tests := []struct {
		value string
		rx    *regexp.Regexp
		want  bool
	}{
		{"abc", regexp.MustCompile(`[a-z]+`), true},
		{"123", regexp.MustCompile(`[a-z]+`), false},
		{"abc123", regexp.MustCompile(`[a-z]+`), true},
		{"ABC", regexp.MustCompile(`[a-z]+`), false},
	}

	for _, tt := range tests {
		got := Matches(tt.value, tt.rx)
		if got != tt.want {
			t.Errorf("Matches(%q, %v) = %v, want %v", tt.value, tt.rx, got, tt.want)
		}
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		values []string
		want   bool
	}{
		{[]string{"apple", "banana", "orange"}, true},
		{[]string{"apple", "apple", "orange"}, false},
		{[]string{"apple", "banana", "banana"}, false},
		{[]string{"apple", "banana", "apple"}, false},
		{[]string{"apple", "banana", "apple", "banana"}, false},
		{[]string{"apple", "banana", "orange", "grape"}, true},
		{[]string{}, true},
	}

	for _, tt := range tests {
		got := Unique(tt.values)
		if got != tt.want {
			t.Errorf("Unique(%v) = %v, want %v", tt.values, got, tt.want)
		}
	}
}
