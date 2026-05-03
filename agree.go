package gomorphy

import "fmt"

// Agree agrees the noun word with the cardinal numeral count
// according to the Russian language rule:
//
//	count mod 100 ∈ [11..14]                   → genitive plural    (12 яблок)
//	otherwise count mod 10 == 1                → nominative singular (1 яблоко)
//	otherwise count mod 10 ∈ [2..4]            → genitive singular   (3 яблока)
//	otherwise                                  → genitive plural     (5 яблок)
//
// Negative numbers are taken by absolute value.
func Agree(word string, count int) (string, error) {
	d, err := getDict()
	if err != nil {
		return "", err
	}
	entries := d.Lookup(word)
	if len(entries) == 0 {
		return "", fmt.Errorf("agreement %q: word not found in dictionary", word)
	}

	c, n := agreeCaseNumber(count)

	for _, e := range entries {
		lemma := &d.Lemmas[e.LemmaID]
		if lemmaPOS(d, lemma.LemmaTags) != POSNoun {
			continue
		}
		s, err := declineNoun(d, lemma, c, n)
		if err != nil {
			return "", fmt.Errorf("agreement %q: %w", word, err)
		}
		return s, nil
	}
	return "", fmt.Errorf("agreement %q: no noun among interpretations", word)
}

// agreeCaseNumber computes the required case and number from the numeral.
// Extracted separately so it is testable and independent of the dictionary.
func agreeCaseNumber(count int) (Case, Number) {
	c := count
	if c < 0 {
		c = -c
	}
	if c%100 >= 11 && c%100 <= 14 {
		return Genitive, Plural
	}
	switch c % 10 {
	case 1:
		return Nominative, Singular
	case 2, 3, 4:
		return Genitive, Singular
	default:
		return Genitive, Plural
	}
}
