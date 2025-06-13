package errors

import "testing"

func TestNew(t *testing.T) {
	e := New("c", "msg", 500)
	if e.Code != "c" || e.Message != "msg" || e.Status != 500 {
		t.Fatalf("unexpected AppError: %+v", e)
	}
	if e.Details != nil {
		t.Fatalf("expected nil details")
	}
	if e.Error() != "msg" {
		t.Fatalf("expected error message 'msg', got %q", e.Error())
	}
}

func TestNewWithDetails(t *testing.T) {
	d := map[string]any{"foo": "bar"}
	e := NewWithDetails("c", "msg", 400, d)
	if e.Code != "c" || e.Message != "msg" || e.Status != 400 {
		t.Fatalf("unexpected AppError: %+v", e)
	}
	if e.Details["foo"] != "bar" {
		t.Fatalf("expected details to contain foo=bar")
	}
}

func TestCommonErrors(t *testing.T) {
	if ErrUserNotFound.Status != 404 {
		t.Errorf("expected 404 status")
	}
	if ErrInvalidInput.Status != 400 {
		t.Errorf("expected 400 status")
	}
}
