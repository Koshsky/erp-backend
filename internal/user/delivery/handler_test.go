package delivery

import (
	"testing"
)

// TestParseIDList covers the ids query param of the batch worker-days endpoint:
// empty entries are skipped, anything non-numeric or non-positive is rejected.
func TestParseIDList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     string
		want    []int64
		wantErr bool
	}{
		{name: "plain list", raw: "1,2,3", want: []int64{1, 2, 3}},
		{name: "spaces around values", raw: " 1 , 2 ", want: []int64{1, 2}},
		{name: "trailing comma", raw: "5,6,", want: []int64{5, 6}},
		{name: "single", raw: "42", want: []int64{42}},
		{name: "empty is an error", raw: "", wantErr: true},
		{name: "only separators is an error", raw: ",,", wantErr: true},
		{name: "zero is an error", raw: "1,0", wantErr: true},
		{name: "negative is an error", raw: "-3", wantErr: true},
		{name: "garbage is an error", raw: "1,abc", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseIDList(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseIDList(%q) = %v, want error", tt.raw, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseIDList(%q) unexpected error: %v", tt.raw, err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("parseIDList(%q) = %v, want %v", tt.raw, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("parseIDList(%q) = %v, want %v", tt.raw, got, tt.want)
				}
			}
		})
	}
}
