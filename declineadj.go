package gomorphy

import (
	"fmt"

	"github.com/therox/gomorphy/internal/dict"
)

// DeclineAdj inflects the adjective word into case c, number n, gender g.
// For masculine singular accusative and for plural accusative the form is
// chosen by the animacy rule: animate=true → genitive form, otherwise →
// nominative form. For feminine and neuter singular accusative the form is
// taken straight from the paradigm.
func DeclineAdj(word string, c Case, n Number, g Gender, animate bool) (string, error) {
	d, err := getDict()
	if err != nil {
		return "", err
	}
	entries := d.Lookup(word)
	if len(entries) == 0 {
		return "", fmt.Errorf("declension of adjective %q: word not found in dictionary", word)
	}
	var lemma *dict.Lemma
	for _, e := range entries {
		l := &d.Lemmas[e.LemmaID]
		if lemmaPOS(d, l.LemmaTags) == POSAdjf {
			lemma = l
			break
		}
	}
	if lemma == nil {
		return "", fmt.Errorf("declension of adjective %q: no ADJF among interpretations", word)
	}

	// Accusative case with substitution: masculine singular and any plural.
	// When animate take genitive, otherwise nominative.
	effectiveCase := c
	if c == Accusative {
		switch {
		case n == Singular && g == Masculine:
			if animate {
				effectiveCase = Genitive
			} else {
				effectiveCase = Nominative
			}
		case n == Plural:
			if animate {
				effectiveCase = Genitive
			} else {
				effectiveCase = Nominative
			}
		}
	}

	cTag := caseTag(effectiveCase)
	nTag := numberTag(n)
	gTag := genderTag(g)
	if cTag == "" || nTag == "" {
		return "", fmt.Errorf("declension of adjective %q: unsupported case/number", word)
	}

	// For plural the gender is not specified in the paradigm — forms are common.
	for _, f := range lemma.Paradigm.Forms {
		if !hasTagStr(d, f.Tags, cTag) || !hasTagStr(d, f.Tags, nTag) {
			continue
		}
		if n == Singular && !hasTagStr(d, f.Tags, gTag) {
			continue
		}
		return f.Text, nil
	}
	return "", fmt.Errorf("paradigm of %q has no adjective form %s/%s/%s",
		lemma.Lemma, cTag, nTag, gTag)
}
