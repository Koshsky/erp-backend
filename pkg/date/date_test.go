package date_test

import (
	"encoding/json"
	"testing"

	"github.com/Koshsky/erp-backend/pkg/date"
)

func TestDateRoundTrip(t *testing.T) {
	t.Parallel()

	d, err := date.Parse("2026-07-15")
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `"2026-07-15"` {
		t.Fatalf("marshal = %s", b)
	}
	var got date.Date
	err = json.Unmarshal(b, &got)
	if err != nil {
		t.Fatal(err)
	}
	if got != d {
		t.Fatalf("round-trip mismatch: %q", got)
	}
}

func TestDateUnmarshalRFC3339Compat(t *testing.T) {
	t.Parallel()

	var d date.Date
	if err := json.Unmarshal([]byte(`"2026-07-15T00:00:00Z"`), &d); err != nil {
		t.Fatal(err)
	}
	if d.String() != "2026-07-15" {
		t.Fatalf("got %q", d)
	}
}

func TestDateRejectsBadFormat(t *testing.T) {
	t.Parallel()

	var d date.Date
	if err := json.Unmarshal([]byte(`"2026/07/15"`), &d); err == nil {
		t.Fatal("expected error")
	}
}

func TestDateTimeIsMidnightUTC(t *testing.T) {
	t.Parallel()

	d, err := date.Parse("2026-07-15")
	if err != nil {
		t.Fatal(err)
	}
	if got := d.Time(); got.UTC().Format(date.Layout) != "2026-07-15" {
		t.Fatalf("got %v", got)
	}
}
