package board

import (
	"errors"
	"strings"
	"testing"
)

func TestParseEPDLine(t *testing.T) {
	line := `4k3/8/8/8/8/8/8/4K2R w K - 0 1 ;D1 15 ; D2 66 ;D3 1197; id "BT2630-14"; note some value with spaces `

	got, err := ParseEPDLine(line)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	if got.FEN != "4k3/8/8/8/8/8/8/4K2R w K - 0 1" {
		t.Errorf("FEN mismatch: %v", got.FEN)
	}

	want := []Operation{
		{Key: "D1", Value: "15"},
		{Key: "D2", Value: "66"},
		{Key: "D3", Value: "1197"},
		{Key: "id", Value: "BT2630-14"},
		{Key: "note", Value: "some value with spaces"},
	}

	if len(got.Operations) != len(want) {
		t.Errorf("length mismatch: got %d, want %d", len(got.Operations), len(want))
	}

	for i := range want {
		if got.Operations[i] != want[i] {
			t.Errorf("Operations[%d]: got %v, want %v", i, got.Operations[i], want[i])
		}
	}
}

func TestParseEPDLine_NoOperations(t *testing.T) {
	got, err := ParseEPDLine("  " + StartingFEN + "  ")
	if err != nil {
		t.Fatalf("unexpected error = %v", err)
	}
	if got.FEN != StartingFEN {
		t.Errorf("FEN mismatch: got %v want %v", got.FEN, StartingFEN)
	}
	if len(got.Operations) != 0 {
		t.Errorf("length mismatch: got %d, want 0", len(got.Operations))
	}
}

func TestParseEPDLine_OperationWithoutValue(t *testing.T) {
	got, err := ParseEPDLine(StartingFEN + "; noop;")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	want := Operation{Key: "noop", Value: ""}
	if len(got.Operations) != 1 {
		t.Fatalf("length mismatch: got %d, want 1", len(got.Operations))
	}
	if got.Operations[0] != want {
		t.Errorf("Operations[0]: got %v, want %v", got.Operations[0], want)
	}
}

func TestParseEPD(t *testing.T) {
	input := "\n" + StartingFEN + "; D1 20;\n" + StartingFEN + "; D2 400\n"

	var got []EPD

	err := ParseEPD(strings.NewReader(input), func(epd EPD) error {
		got = append(got, epd)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("length mismatch: got %d, want 2", len(got))
	}
	if got[0].Operations[0] != (Operation{Key: "D1", Value: "20"}) {
		t.Errorf("Operations[0]: got %v, want %v", got[0].Operations[0], Operation{Key: "D1", Value: "20"})
	}
	if got[1].Operations[0] != (Operation{Key: "D2", Value: "400"}) {
		t.Errorf("Operations[1]: got %v, want %v", got[1].Operations[0], Operation{Key: "D2", Value: "400"})
	}
}

func TestParseEPD_PropagatesCallbackError(t *testing.T) {
	wantErr := errors.New("stop")
	err := ParseEPD(strings.NewReader(StartingFEN), func(EPD) error {
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Errorf("unexpected error: got %v, want %v", err, wantErr)
	}
}
