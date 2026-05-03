package gomorphy

// Parse performs morphological analysis of the word form word.
// Returns all possible interpretations (homonyms); if the form is not
// found, an empty non-nil slice and a nil error are returned (see
// DESIGN.md, section 7).
//
// The Animate field is derived from the anim/inan tag: either on the
// form or on the lemma. The Gender field for NOUN is taken from the
// lemma; for ADJF it is taken from the form (for singular) or remains
// GenderUnknown (for plural, where gender is not distinguished in the
// paradigm).
func Parse(word string) ([]Analysis, error) {
	d, err := getDict()
	if err != nil {
		return nil, err
	}
	entries := d.Lookup(word)
	out := make([]Analysis, 0, len(entries))
	for _, e := range entries {
		lemma := &d.Lemmas[e.LemmaID]
		pos := lemmaPOS(d, lemma.LemmaTags)

		a := Analysis{
			Lemma: lemma.Lemma,
			POS:   pos,
		}

		// Case and number come from the form tags; for indeclinables the
		// form tags only cover one position, so Case/Number may stay
		// Unknown, which correctly conveys "no grammatical information".
		for _, t := range e.FormTags {
			s := d.TagString(t)
			if c := parseCase(s); c != CaseUnknown {
				a.Case = c
			}
			if n := parseNumber(s); n != NumberUnknown {
				a.Number = n
			}
			if pos == POSAdjf {
				if g := parseGender(s); g != GenderUnknown {
					a.Gender = g
				}
			}
			if s == "anim" {
				a.Animate = true
			}
		}

		// For nouns, gender is a property of the lemma.
		if pos == POSNoun {
			for _, t := range lemma.LemmaTags {
				if g := parseGender(d.TagString(t)); g != GenderUnknown {
					a.Gender = g
					break
				}
			}
		}
		// Animacy may be marked only on the lemma (the typical case).
		if !a.Animate && hasTagStr(d, lemma.LemmaTags, "anim") {
			a.Animate = true
		}

		out = append(out, a)
	}
	return out, nil
}
