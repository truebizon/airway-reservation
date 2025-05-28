package myerror

import "testing"

func TestConvertToHTTPError(t *testing.T) {
	err := Wrap(BadRequest, Errorf("", "bad"), "bad request")
	status, msg := ConvertToHTTPError(err)
	if status != 400 {
		t.Errorf("expected 400 got %d", status)
	}
	if msg == "" {
		t.Errorf("expected message")
	}
}
