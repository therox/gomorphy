package gomorphy

import (
	"strconv"
	"testing"
)

// positiveValues is the full list of boundary values from the Phase 2 plan.
// The mirror negative set is built as `-v` for each `v != 0`
// (a separate 0 for negatives makes no sense — it coincides with 0).
var positiveValues = []int{
	0, 1, 2, 4, 5, 11, 12, 14, 15, 21, 22, 25, 101, 111, 121, 1000,
}

// expectedAgreeKey returns the expected case/number for |count| by the rule
// 11–14 / 1 / 2–4 / 5+. Intentionally duplicates the agreeCaseNumber logic
// independently so the test is self-checking (if the production rule is
// wrong, the table-driven check will catch the divergence).
func expectedAgreeKey(count int) (Case, Number) {
	v := count
	if v < 0 {
		v = -v
	}
	if v%100 >= 11 && v%100 <= 14 {
		return Genitive, Plural
	}
	switch v % 10 {
	case 1:
		return Nominative, Singular
	case 2, 3, 4:
		return Genitive, Singular
	default:
		return Genitive, Plural
	}
}

// allAgreeValues returns the positive boundary values together with their
// mirror negatives. Used by both table-driven tests.
func allAgreeValues() []int {
	out := make([]int, 0, len(positiveValues)*2)
	out = append(out, positiveValues...)
	for _, v := range positiveValues {
		if v == 0 {
			continue
		}
		out = append(out, -v)
	}
	return out
}

// TestAgreeCaseNumber is a table-driven check of the case/number selection
// rule without consulting the dictionary. Full positive set plus mirror
// negative set (per the plan: "for negatives — the same by absolute value").
func TestAgreeCaseNumber(t *testing.T) {
	for _, v := range allAgreeValues() {
		v := v
		t.Run(strconv.Itoa(v), func(t *testing.T) {
			gotC, gotN := agreeCaseNumber(v)
			wantC, wantN := expectedAgreeKey(v)
			if gotC != wantC || gotN != wantN {
				t.Fatalf("agreeCaseNumber(%d) = (%v,%v), expected (%v,%v)",
					v, gotC, gotN, wantC, wantN)
			}
		})
	}
}

// TestAgreeYabloko is an integration run against the dictionary for "яблоко"
// over the same mirror set of positive and negative values.
func TestAgreeYabloko(t *testing.T) {
	// Mapping case/number → expected form of "яблоко".
	wantText := func(c Case, n Number) string {
		switch {
		case c == Nominative && n == Singular:
			return "яблоко"
		case c == Genitive && n == Singular:
			return "яблока"
		case c == Genitive && n == Plural:
			return "яблок"
		default:
			return ""
		}
	}

	for _, v := range allAgreeValues() {
		v := v
		t.Run(strconv.Itoa(v), func(t *testing.T) {
			c, n := expectedAgreeKey(v)
			want := wantText(c, n)
			if want == "" {
				t.Fatalf("test has no expected form for (%v,%v)", c, n)
			}
			got, err := Agree("яблоко", v)
			if err != nil {
				t.Fatalf("Agree(\"яблоко\", %d) error: %v", v, err)
			}
			if got != want {
				t.Fatalf("Agree(\"яблоко\", %d) = %q, expected %q", v, got, want)
			}
		})
	}
}

// TestAgreeWordNotFound — an unknown word returns an error.
func TestAgreeWordNotFound(t *testing.T) {
	_, err := Agree("абракадабра", 5)
	if err == nil {
		t.Fatal("expected an error for an unknown word")
	}
}
