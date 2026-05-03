package gomorphy

import "testing"

// TestParseSingle — an unambiguous word form returns a single analysis.
func TestParseSingle(t *testing.T) {
	got, err := Parse("аппетита")
	if err != nil {
		t.Fatalf("Parse(\"аппетита\") error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected exactly one analysis, got %d: %+v", len(got), got)
	}
	a := got[0]
	if a.Lemma != "аппетит" || a.POS != POSNoun ||
		a.Case != Genitive || a.Number != Singular ||
		a.Gender != Masculine || a.Animate {
		t.Fatalf("unexpected analysis: %+v", a)
	}
}

// TestParseHomonymy — "стали" has several analyses inside the single
// lemma "сталь": gent.sing, datv.sing, loct.sing, nomn.plur, accs.plur.
func TestParseHomonymy(t *testing.T) {
	got, err := Parse("стали")
	if err != nil {
		t.Fatalf("Parse(\"стали\") error: %v", err)
	}
	if len(got) < 2 {
		t.Fatalf("expected several analyses, got %d: %+v", len(got), got)
	}

	// Map of observed "case/number" pairs for convenient assertion that
	// specific analyses are present.
	type key struct {
		c Case
		n Number
	}
	seen := map[key]bool{}
	for _, a := range got {
		if a.Lemma != "сталь" {
			t.Fatalf("expected lemma \"сталь\", got %q", a.Lemma)
		}
		if a.POS != POSNoun || a.Gender != Feminine {
			t.Fatalf("unexpected analysis: %+v", a)
		}
		seen[key{a.Case, a.Number}] = true
	}
	required := []key{
		{Genitive, Singular},
		{Dative, Singular},
		{Prepositional, Singular},
		{Nominative, Plural},
	}
	for _, r := range required {
		if !seen[r] {
			t.Fatalf("the analysis is missing the expected pair %v/%v: %+v", r.c, r.n, got)
		}
	}
}

// TestParseAdjective — for an adjective the gender is derived from the
// form tags (for singular).
func TestParseAdjective(t *testing.T) {
	got, err := Parse("красная")
	if err != nil {
		t.Fatalf("Parse(\"красная\") error: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("expected at least one analysis")
	}
	found := false
	for _, a := range got {
		if a.Lemma == "красный" && a.POS == POSAdjf &&
			a.Case == Nominative && a.Number == Singular && a.Gender == Feminine {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected analysis красная/femn/sing/nomn not found: %+v", got)
	}
}

// TestParseAnimate — animacy is carried over from the lemma tags.
func TestParseAnimate(t *testing.T) {
	got, err := Parse("ивана")
	if err != nil {
		t.Fatalf("Parse(\"ивана\") error: %v", err)
	}
	if len(got) == 0 {
		t.Fatal("expected at least one analysis")
	}
	hasAnim := false
	for _, a := range got {
		if a.Animate && a.Lemma == "иван" {
			hasAnim = true
			break
		}
	}
	if !hasAnim {
		t.Fatalf("expected Animate=true for the lemma \"иван\": %+v", got)
	}
}

// TestParseNotFound — a word outside the dictionary returns an empty
// slice and a nil error.
func TestParseNotFound(t *testing.T) {
	got, err := Parse("абракадабра")
	if err != nil {
		t.Fatalf("expected a nil error, got %v", err)
	}
	if got == nil {
		t.Fatal("expected a non-nil zero-length slice, got nil")
	}
	if len(got) != 0 {
		t.Fatalf("expected an empty slice, got %+v", got)
	}
}

// TestParseEYoNormalization — lookup works the same way for "ё" and "е".
func TestParseEYoNormalization(t *testing.T) {
	a, err := Parse("ёж")
	if err != nil {
		t.Fatalf("Parse(\"ёж\") error: %v", err)
	}
	b, err := Parse("еж")
	if err != nil {
		t.Fatalf("Parse(\"еж\") error: %v", err)
	}
	if len(a) == 0 || len(b) == 0 {
		t.Fatalf("expected analyses for both forms, got: a=%+v, b=%+v", a, b)
	}
	if a[0].Lemma != b[0].Lemma {
		t.Fatalf("ё/е are parsed differently: %+v vs %+v", a[0], b[0])
	}
}
