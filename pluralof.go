package gomorphy

import "fmt"

// PluralOf returns the nominative plural form of word.
// The lemma is determined from the input word form.
// For pluralia tantum the lemma itself is returned (it is already plural).
func PluralOf(word string) (string, error) {
	d, err := getDict()
	if err != nil {
		return "", err
	}
	entries := d.Lookup(word)
	if len(entries) == 0 {
		return "", fmt.Errorf("plural of %q: word not found in dictionary", word)
	}
	for _, e := range entries {
		lemma := &d.Lemmas[e.LemmaID]
		if lemmaPOS(d, lemma.LemmaTags) != POSNoun {
			continue
		}
		// Pluralia tantum: the lemma is already plural, there is no separate form.
		if hasTagStr(d, lemma.LemmaTags, "Pltm") {
			return lemma.Lemma, nil
		}
		s, err := declineNoun(d, lemma, Nominative, Plural)
		if err != nil {
			return "", fmt.Errorf("plural of %q: %w", word, err)
		}
		return s, nil
	}
	return "", fmt.Errorf("plural of %q: no noun among interpretations", word)
}
