package jsonlog

import (
	"bytes"
	"errors"
	"runtime/debug"
	"strings"
	"testing"
)

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{LevelInfo, "INFO"},
		{LevelError, "ERROR"},
		{LevelFatal, "FATAL"},
		{Level(100), ""},
	}

	for _, test := range tests {
		result := test.level.String()
		if result != test.expected {
			t.Errorf("Level.String() returned %q, expected %q", result, test.expected)
		}
	}
}

func TestNew(t *testing.T) {
	out := &bytes.Buffer{}
	minLevel := LevelInfo

	logger := New(out, minLevel)

	if logger.out != out {
		t.Errorf("New() did not set the correct 'out' field")
	}

	if logger.minLevel != minLevel {
		t.Errorf("New() did not set the correct 'minLevel' field")
	}
}

func TestPrintInfo(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf, LevelInfo)

	message := "This is an info message"
	properties := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	logger.PrintInfo(message, properties)

	expectedPrefix := `{"level":"INFO",`
	expectedSuffix := `,"properties":{"key1":"value1","key2":"value2"}}` + "\n"
	actualOutput := buf.String()

	if !strings.HasPrefix(actualOutput, expectedPrefix) {
		t.Errorf("actualOutput does not have the expected prefix. pref: %q, suf: %q, got: %q", expectedPrefix, expectedSuffix, actualOutput)
	}

	if !strings.HasSuffix(actualOutput, expectedSuffix) {
		t.Errorf("actualOutput does not have the expected suffix. pref: %q, suf: %q, got: %q", expectedPrefix, expectedSuffix, actualOutput)
	}
}

func TestPrintError(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf, LevelInfo)

	err := errors.New("This is an error")
	properties := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	logger.PrintError(err, properties)

	expectedPrefix := `{"level":"ERROR",`
	actualOutput := buf.String()

	if !strings.HasPrefix(actualOutput, expectedPrefix) {
		t.Errorf("actualOutput does not have the expected prefix. pref: %q, got: %q", expectedPrefix, actualOutput)
	}
}
func TestLogger_Write(t *testing.T) {
	tests := []struct {
		message  []byte
		expected string
	}{
		{[]byte("This is an error message"), `{"level":"ERROR","message":"This is an error message","properties":{}}` + "\n"},
		{[]byte("Another error message"), `{"level":"ERROR","message":"Another error message","properties":{}}` + "\n"},
	}

	for _, test := range tests {
		var buf bytes.Buffer
		logger := New(&buf, LevelInfo)

		_, err := logger.Write(test.message)
		if err != nil {
			t.Errorf("Logger.Write() returned an unexpected error: %v", err)
		}

		expectedPrefix := `{"level":"ERROR",`
		actualOutput := buf.String()
		if !strings.HasPrefix(actualOutput, expectedPrefix) {
			t.Errorf("actualOutput does not have the expected prefix. pref: %q, got: %q", expectedPrefix, actualOutput)
		}
	}
}

func TestLogger_print(t *testing.T) {
	tests := []struct {
		level      Level
		message    string
		properties map[string]string
		expected   string
	}{
		{
			level:      LevelInfo,
			message:    "This is an info message",
			properties: map[string]string{"key1": "value1", "key2": "value2"},
			expected:   `{"level":"INFO","time":"2022-01-01T00:00:00Z","message":"This is an info message","properties":{"key1":"value1","key2":"value2"}}` + "\n",
		},
		{
			level:      LevelError,
			message:    "This is an error",
			properties: map[string]string{"key1": "value1", "key2": "value2"},
			expected:   `{"level":"ERROR","time":"2022-01-01T00:00:00Z","message":"This is an error","properties":{"key1":"value1","key2":"value2"},"trace":"` + string(debug.Stack()) + `"` + "\n",
		},
	}

	for _, test := range tests {
		var buf bytes.Buffer
		logger := &Logger{
			out:      &buf,
			minLevel: LevelInfo,
		}

		_, err := logger.print(test.level, test.message, test.properties)
		if err != nil {
			t.Errorf("Logger.print() returned an unexpected error: %v", err)
		}

		actualOutput := buf.String()

		if test.level == LevelInfo {
			expectedPrefix := `{"level":"INFO",`
			if !strings.HasPrefix(actualOutput, expectedPrefix) {
				t.Errorf("actualOutput does not have the expected prefix. pref: %q, got: %q", expectedPrefix, actualOutput)
			}
		} else if test.level == LevelError {
			expectedPrefix := `{"level":"ERROR",`
			if !strings.HasPrefix(actualOutput, expectedPrefix) {
				t.Errorf("actualOutput does not have the expected prefix. pref: %q, got: %q", expectedPrefix, actualOutput)
			}
		} else {
			t.Errorf("unexpected level: %v", test.level)
		}
	}
}
