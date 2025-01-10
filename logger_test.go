package logger

import (
	"bytes"
	"errors"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(WithOutput(buf))
	if l == nil {
		t.Fatal("Expected logger to be created")
	}

	err := l.Log("Hello, World!")
	if err != nil {
		t.Error("Expected no error")
	}
	if buf.String() != "[INFO]: Hello, World!\n" {
		t.Errorf("Expected %q, got %q", "[INFO]: Hello, World!\n", buf.String())
	}
}

func TestNewWithNilOutput(t *testing.T) {
	l := New(WithOutput(nil))
	if l == nil {
		t.Fatal("Expected logger to be created")
	}

	err := l.Log("Hello, World!")
	if err == nil {
		t.Errorf("Expected error %q, got nil", CannotWriteError)
	} else if err.Error() != InvalidOutputError.Error() {
		t.Errorf("Expected error %q, got %q", CannotWriteError, err)
	}
}

// CappedWriter is a mock io.Writer that only accepts writes up to a certain length.
type CappedWriter struct {
	cap int
}

func (c *CappedWriter) Write(p []byte) (n int, err error) {
	if len(p) > c.cap {
		return 0, errors.New("input too long")
	}
	return len(p), nil
}

func TestNewWithInvalidOutput(t *testing.T) {
	w := &CappedWriter{cap: 1}
	l := New(WithOutput(w))
	if l == nil {
		t.Fatal("Expected logger to be created")
	}

	err := l.Log("Hello, World!")
	if err == nil {
		t.Errorf("Expected error %q, got nil", CannotWriteError)
	} else if err.Error() != CannotWriteError.Error() {
		t.Errorf("Expected error %q, got %q", CannotWriteError, err)
	}
}

func TestNewWithPrefix(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(WithOutput(buf), WithPrefix("TEST"))
	if l == nil {
		t.Fatal("Expected logger to be created")
	}

	err := l.Log("Hello, World!")
	if err != nil {
		t.Error("Expected no error")
	}
	if buf.String() != "TEST [INFO]: Hello, World!\n" {
		t.Errorf("Expected %q, got %q", "TEST [INFO]: Hello, World!\n", buf.String())
	}
}

func TestNewWithTime(t *testing.T) {
	buf := new(bytes.Buffer)
	staticTime := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)
	mockTimeNow := func() time.Time { return staticTime }
	expected := "2023-01-01 12:00:00 [INFO]: Hello, World!\n"

	l := New(WithOutput(buf), WithTimeIncluded(), WithClock(mockTimeNow))
	if l == nil {
		t.Fatal("Expected logger to be created")
	}

	err := l.Log("Hello, World!")
	if err != nil {
		t.Error("Expected no error")
	}
	if buf.String() != expected {
		t.Errorf("Expected %s, got %s", expected, buf.String())
	}
}

func TestWithTimeAndPrefix(t *testing.T) {
	buf := new(bytes.Buffer)
	staticTime := time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC)
	mockTimeNow := func() time.Time { return staticTime }
	expected := "TEST/2023-01-01 12:00:00 [INFO]: Hello, World!\n"

	l := New(
		WithOutput(buf),
		WithTimeIncluded(),
		WithPrefix("TEST"),
		WithClock(mockTimeNow),
	)
	if l == nil {
		t.Fatal("Expected logger to be created")
	}

	err := l.Log("Hello, World!")
	if err != nil {
		t.Error("Expected no error")
	}
	if buf.String() != expected {
		t.Errorf("Expected 34 characters, got %d", len(buf.String()))
	}

}

func TestWarnLogWithSeverity(t *testing.T) {
	buf := new(bytes.Buffer)
	expected := "[WARN (High)]: This is a warning\n"

	l := New(WithOutput(buf))
	if l == nil {
		t.Fatal("Expected logger to be created")
	}

	err := l.Warn("This is a warning", WarnSevere)
	if err != nil {
		t.Error("Expected no error")
	}
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}

func TestErrorLog(t *testing.T) {
	buf := new(bytes.Buffer)
	l := New(WithOutput(buf), WithPrefix("TEST"))
	expected := "TEST [ERROR (500)]: this is an error\n"
	if l == nil {
		t.Fatal("Expected logger to be created")
	}

	err := l.Error(errors.New("this is an error"), 500)
	if err != nil {
		t.Error("Expected no error")
	}
	if buf.String() != expected {
		t.Errorf("Expected %q, got %q", expected, buf.String())
	}
}
