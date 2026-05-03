package gomorphy

import "fmt"

// GenderOf returns the gender of word from the dictionary.
// The first noun lemma among the interpretations is used.
// For pluralia tantum the gender is not specified in the dictionary →
// GenderUnknown is returned without an error.
func GenderOf(word string) (Gender, error) {
	d, err := getDict()
	if err != nil {
		return GenderUnknown, err
	}
	entries := d.Lookup(word)
	if len(entries) == 0 {
		return GenderUnknown, fmt.Errorf("gender lookup for %q: word not found in dictionary", word)
	}
	for _, e := range entries {
		lemma := &d.Lemmas[e.LemmaID]
		if lemmaPOS(d, lemma.LemmaTags) != POSNoun {
			continue
		}
		for _, t := range lemma.LemmaTags {
			if g := parseGender(d.TagString(t)); g != GenderUnknown {
				return g, nil
			}
		}
		// Gender is absent on the lemma — for example, on pluralia tantum.
		return GenderUnknown, nil
	}
	return GenderUnknown, fmt.Errorf("gender lookup for %q: no noun among interpretations", word)
}
