package gomorphy

import (
	"fmt"

	"github.com/therox/gomorphy/internal/dict"
)

// Decline inflects the noun word into case c and number n.
// The word is looked up in the dictionary's reverse index; the first noun
// lemma is used. For indeclinable nouns (Fixd) the lemma form is returned
// for any case/number. For pluralia tantum (Pltm) a singular request → error.
func Decline(word string, c Case, n Number) (string, error) {
	d, err := getDict()
	if err != nil {
		return "", err
	}
	entries := d.Lookup(word)
	if len(entries) == 0 {
		return "", fmt.Errorf("declension of %q: word not found in dictionary", word)
	}
	for _, e := range entries {
		lemma := &d.Lemmas[e.LemmaID]
		if lemmaPOS(d, lemma.LemmaTags) != POSNoun {
			continue
		}
		s, err := declineNoun(d, lemma, c, n)
		if err != nil {
			return "", fmt.Errorf("declension of %q: %w", word, err)
		}
		return s, nil
	}
	return "", fmt.Errorf("declension of %q: no noun among interpretations", word)
}

// declineNoun looks up the form in the lemma's paradigm by the given case and number.
// Shared helper used by both Decline and Agree.
func declineNoun(d *dict.Dict, lemma *dict.Lemma, c Case, n Number) (string, error) {
	// Indeclinable noun: "кофе", "метро" — the single form is returned
	// for any case/number.
	if hasTagStr(d, lemma.LemmaTags, "Fixd") {
		return lemma.Lemma, nil
	}
	// Pluralia tantum: "ножницы", "сани" — no singular forms.
	if n == Singular && hasTagStr(d, lemma.LemmaTags, "Pltm") {
		return "", fmt.Errorf("lemma %q is pluralia tantum, no singular forms", lemma.Lemma)
	}

	cTag := caseTag(c)
	nTag := numberTag(n)
	if cTag == "" || nTag == "" {
		return "", fmt.Errorf("unsupported case/number combination for lemma %q", lemma.Lemma)
	}
	for _, f := range lemma.Paradigm.Forms {
		if hasTagStr(d, f.Tags, cTag) && hasTagStr(d, f.Tags, nTag) {
			return f.Text, nil
		}
	}
	return "", fmt.Errorf("paradigm of %q has no form %s/%s", lemma.Lemma, cTag, nTag)
}
