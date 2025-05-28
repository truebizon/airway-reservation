package util

import "testing"

func TestEncodeDecode(t *testing.T) {
	original := "a73b0f60-574c-409c-a9ab-3bebeb60dcfa"
	encoded := Encode(original)
	if encoded == "" {
		t.Fatalf("encode returned empty string")
	}
	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if decoded.String() != original {
		t.Errorf("expected %s got %s", original, decoded.String())
	}
}
