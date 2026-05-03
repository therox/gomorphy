package gomorphy

// File: implementation of DeclineFullName.
// Algorithm:
//   1. detectGender by priority: patronymic → first name → surname.
//   2. Each component is inflected separately: first / patronymic / surname.
//      Inside: dictionary → heuristic (for patronymics and surnames).
//   3. The case of the first letter of the source is carried over to the form.

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/therox/gomorphy/internal/dict"
)

// declineFullName is the internal implementation of the public DeclineFullName.
// Empty FullName fields stay empty and do not cause errors.
func declineFullName(name FullName, c Case) (FullName, error) {
	if name.Last == "" && name.First == "" && name.Patronymic == "" {
		return FullName{}, errors.New("DeclineFullName: empty name")
	}
	g, err := detectGender(name)
	if err != nil {
		return FullName{}, fmt.Errorf("declension of full name %q: %w", joinFullName(name), err)
	}

	out := FullName{}
	if name.Last != "" {
		s, err := declineLast(name.Last, c, g)
		if err != nil {
			return FullName{}, fmt.Errorf("declension of full name %q: %w", joinFullName(name), err)
		}
		out.Last = s
	}
	if name.First != "" {
		s, err := declineFirst(name.First, c, g)
		if err != nil {
			return FullName{}, fmt.Errorf("declension of full name %q: %w", joinFullName(name), err)
		}
		out.First = s
	}
	if name.Patronymic != "" {
		s, err := declinePatronymic(name.Patronymic, c, g)
		if err != nil {
			return FullName{}, fmt.Errorf("declension of full name %q: %w", joinFullName(name), err)
		}
		out.Patronymic = s
	}
	return out, nil
}

// detectGender determines gender using a priority of sources: patronymic →
// first name → surname. Each source is skipped if the corresponding field
// is empty. If no source yields a definite gender — an error is returned.
//
// Common (ms-f) from a dictionary surname is treated as a valid result —
// it is used as "any" in genderCompatible: for indeclinable foreign
// surnames such as "Дюма" that is enough.
func detectGender(name FullName) (Gender, error) {
	if name.Patronymic != "" {
		if g := genderFromPatronymic(name.Patronymic); g != GenderUnknown {
			return g, nil
		}
	}
	if name.First != "" {
		if g := genderFromFirstName(name.First); g != GenderUnknown && g != Common {
			return g, nil
		}
	}
	if name.Last != "" {
		if g := genderFromSurname(name.Last); g != GenderUnknown {
			return g, nil
		}
	}
	return GenderUnknown, errors.New("could not determine gender")
}

// genderFromFirstName returns the gender from the dictionary Name lemma.
// If the given name is not in the dictionary — GenderUnknown.
func genderFromFirstName(first string) Gender {
	d, err := getDict()
	if err != nil {
		return GenderUnknown
	}
	for _, e := range d.Lookup(first) {
		l := &d.Lemmas[e.LemmaID]
		if !hasTagStr(d, l.LemmaTags, "Name") {
			continue
		}
		for _, t := range l.LemmaTags {
			if g := parseGender(d.TagString(t)); g != GenderUnknown {
				return g
			}
		}
	}
	return GenderUnknown
}

// genderFromSurname returns the gender from the dictionary Surn lemma,
// or the result of the suffix heuristic.
//
// The heuristic is checked first: otherwise a dictionary lookup for
// "Иванова" would first find the genitive form of the masculine lemma
// "Иванов" and would return Masculine — which is wrong for the input
// nominative form "Иванова".
func genderFromSurname(last string) Gender {
	if g := genderFromSurnameHeuristic(last); g != GenderUnknown {
		return g
	}
	d, err := getDict()
	if err != nil {
		return GenderUnknown
	}
	for _, e := range d.Lookup(last) {
		l := &d.Lemmas[e.LemmaID]
		if !hasTagStr(d, l.LemmaTags, "Surn") {
			continue
		}
		// Consider only lemmas whose base form equals the input.
		// Otherwise for "Иванова" we would pick the masculine lemma "иванов".
		if dict.NormalizeForm(l.Lemma) != dict.NormalizeForm(last) {
			continue
		}
		for _, t := range l.LemmaTags {
			if g := parseGender(d.TagString(t)); g != GenderUnknown {
				return g
			}
		}
	}
	return GenderUnknown
}

// declineFirst inflects a given name.
// The name is searched among lemmas tagged Name with the desired (or
// compatible) gender. If the dictionary has nothing or the required form
// is missing in the paradigm — the name is returned as is (indeclinable
// foreign).
func declineFirst(first string, c Case, g Gender) (string, error) {
	d, err := getDict()
	if err != nil {
		return "", err
	}
	for _, e := range d.Lookup(first) {
		l := &d.Lemmas[e.LemmaID]
		if !hasTagStr(d, l.LemmaTags, "Name") {
			continue
		}
		if !genderCompatible(d, l, g) {
			continue
		}
		s, err := declineNoun(d, l, c, Singular)
		if err == nil {
			return applyCase(first, s), nil
		}
		// Incomplete paradigm — return the original form (without an error).
		return first, nil
	}
	return first, nil
}

// declinePatronymic inflects a patronymic: dictionary first, otherwise heuristic.
func declinePatronymic(patr string, c Case, g Gender) (string, error) {
	d, err := getDict()
	if err != nil {
		return "", err
	}
	for _, e := range d.Lookup(patr) {
		l := &d.Lemmas[e.LemmaID]
		if !hasTagStr(d, l.LemmaTags, "Patr") {
			continue
		}
		if !genderCompatible(d, l, g) {
			continue
		}
		if s, err := declineNoun(d, l, c, Singular); err == nil {
			return applyCase(patr, s), nil
		}
		break
	}
	s, ok := declinePatronymicHeuristic(strings.ToLower(patr), c, g)
	if !ok {
		return patr, nil
	}
	return applyCase(patr, s), nil
}

// declineLast inflects a surname: dictionary first, otherwise heuristic.
func declineLast(last string, c Case, g Gender) (string, error) {
	d, err := getDict()
	if err != nil {
		return "", err
	}
	for _, e := range d.Lookup(last) {
		l := &d.Lemmas[e.LemmaID]
		if !hasTagStr(d, l.LemmaTags, "Surn") {
			continue
		}
		if !genderCompatible(d, l, g) {
			continue
		}
		if s, err := declineNoun(d, l, c, Singular); err == nil {
			return applyCase(last, s), nil
		}
		break
	}
	s, ok := declineSurnameHeuristic(strings.ToLower(last), c, g)
	if !ok {
		return last, nil
	}
	return applyCase(last, s), nil
}

// genderCompatible — the lemma fits the requested gender.
// Common (ms-f) is compatible with any want; GenderUnknown is also accepted.
func genderCompatible(d *dict.Dict, l *dict.Lemma, want Gender) bool {
	if want == GenderUnknown || want == Common {
		return true
	}
	var lemmaG Gender
	for _, t := range l.LemmaTags {
		if g := parseGender(d.TagString(t)); g != GenderUnknown {
			lemmaG = g
			break
		}
	}
	if lemmaG == GenderUnknown || lemmaG == Common {
		return true
	}
	return lemmaG == want
}

// applyCase copies the case of the first letter of orig onto formed.
// If the first letter of orig is uppercase, it uppercases the first letter
// of formed. Used to preserve "Иванов" → "Иванова", "иванов" → "иванова".
func applyCase(orig, formed string) string {
	if orig == "" || formed == "" {
		return formed
	}
	rOrig, _ := utf8.DecodeRuneInString(orig)
	if !unicode.IsUpper(rOrig) {
		return formed
	}
	rFormed, sz := utf8.DecodeRuneInString(formed)
	return string(unicode.ToUpper(rFormed)) + formed[sz:]
}

// joinFullName joins the components of a full name for error messages.
// Empty fields are skipped.
func joinFullName(n FullName) string {
	parts := make([]string, 0, 3)
	if n.Last != "" {
		parts = append(parts, n.Last)
	}
	if n.First != "" {
		parts = append(parts, n.First)
	}
	if n.Patronymic != "" {
		parts = append(parts, n.Patronymic)
	}
	return strings.Join(parts, " ")
}

// toNominativeImpl is the implementation of the public ToNominative.
// Each component is reduced to the nominative independently. Algorithm:
//  1. Patronymic: inversePatronymic — yields Nom and a gender hint.
//  2. First name: dictionary (Parse → Lemma); if absent, leave as is.
//  3. Surname: dictionary (Lookup + Surn) with gender hint, otherwise
//     inverseSurnameHeuristic with the gender hint.
//
// The case of the first letter of each component is carried over from the
// original form.
func toNominativeImpl(name FullName) (FullName, error) {
	if name.Last == "" && name.First == "" && name.Patronymic == "" {
		return FullName{}, errors.New("ToNominative: empty name")
	}

	out := FullName{}
	var gHint Gender

	if name.Patronymic != "" {
		if nom, g, ok := inversePatronymic(name.Patronymic); ok {
			out.Patronymic = applyCase(name.Patronymic, nom)
			gHint = g
		} else {
			out.Patronymic = name.Patronymic
		}
	}

	if name.First != "" {
		nom, g, ok := firstNameNom(name.First)
		if ok {
			out.First = applyCase(name.First, nom)
			if gHint == GenderUnknown && g != GenderUnknown && g != Common {
				gHint = g
			}
		} else {
			out.First = name.First
		}
	}

	if name.Last != "" {
		nom, g, ok := surnameNom(name.Last, gHint)
		if ok {
			out.Last = applyCase(name.Last, nom)
			if gHint == GenderUnknown && g != GenderUnknown && g != Common {
				gHint = g
			}
		} else {
			out.Last = name.Last
		}
	}

	return out, nil
}

// firstNameNom reduces a given name to the nominative via the dictionary (Parse).
// Out-of-dictionary names are returned as (s, Unknown, false), signalling
// "leave as is".
func firstNameNom(first string) (string, Gender, bool) {
	d, err := getDict()
	if err != nil {
		return first, GenderUnknown, false
	}
	for _, e := range d.Lookup(first) {
		l := &d.Lemmas[e.LemmaID]
		if !hasTagStr(d, l.LemmaTags, "Name") {
			continue
		}
		var g Gender
		for _, t := range l.LemmaTags {
			if gg := parseGender(d.TagString(t)); gg != GenderUnknown {
				g = gg
				break
			}
		}
		return l.Lemma, g, true
	}
	return first, GenderUnknown, false
}

// surnameNom reduces a surname to the nominative. Dictionary first,
// otherwise the heuristic. Honors gHint (a gender hint from the other
// components) — for the dictionary it filters lemmas, for the heuristic it
// disambiguates ambiguous forms.
//
// When gHint is absent — the lemma whose base form matches the input is
// preferred (so "Иванова" returns the feminine lemma "иванова" rather than
// the masculine "иванов" with the genitive form "иванова").
func surnameNom(last string, gHint Gender) (string, Gender, bool) {
	d, err := getDict()
	if err == nil {
		// Two-pass search: first look for a lemma whose base form equals
		// the input ("input is already in Nom" priority), then for any
		// compatible one.
		normIn := dict.NormalizeForm(last)
		for pass := 0; pass < 2; pass++ {
			for _, e := range d.Lookup(last) {
				l := &d.Lemmas[e.LemmaID]
				if !hasTagStr(d, l.LemmaTags, "Surn") {
					continue
				}
				if pass == 0 && dict.NormalizeForm(l.Lemma) != normIn {
					continue
				}
				var g Gender
				for _, t := range l.LemmaTags {
					if gg := parseGender(d.TagString(t)); gg != GenderUnknown {
						g = gg
						break
					}
				}
				// If gHint is set and the lemma has a definite gender — filter.
				// Common (ms-f) is accepted with any hint.
				if gHint != GenderUnknown && gHint != Common &&
					g != GenderUnknown && g != Common && g != gHint {
					continue
				}
				return l.Lemma, g, true
			}
		}
	}
	return inverseSurnameHeuristic(last, gHint)
}

// parseFullNameImpl is the implementation of the public ParseFullName.
// Algorithm:
//  1. Tokenization on spaces (1, 2, or 3 tokens).
//  2. Detect each token's role by suffix:
//     - a token with a patronymic suffix → Patronymic;
//     - the remaining two — Last/First; the order is inferred from the
//       patronymic position:
//       Patr at position [1] → Russian "Last First Patr";
//       Patr at position [2] → Western "First Patr Last";
//     - for two tokens without a patronymic: the one recognized as a Name
//       in the dictionary is First; the other is Last;
//     - for a single token: first try as Patronymic, then as Surn (via
//       dictionary or ending), otherwise as First.
//  3. ToNominative on the resulting FullName.
func parseFullNameImpl(s string) (FullName, error) {
	tokens := strings.Fields(strings.TrimSpace(s))
	if len(tokens) == 0 {
		return FullName{}, errors.New("ParseFullName: empty string")
	}
	if len(tokens) > 3 {
		return FullName{}, fmt.Errorf("ParseFullName: expected 1–3 tokens, got %d", len(tokens))
	}

	in := FullName{}
	switch len(tokens) {
	case 1:
		in = classifySingle(tokens[0])
	case 2:
		in = classifyDouble(tokens[0], tokens[1])
	case 3:
		var err error
		in, err = classifyTriple(tokens[0], tokens[1], tokens[2])
		if err != nil {
			return FullName{}, err
		}
	}

	return toNominativeImpl(in)
}

// classifySingle determines the role of a single token.
// Heuristic: patronymic → surname (if recognizable) → first name.
func classifySingle(tok string) FullName {
	if looksLikePatronymic(tok) {
		return FullName{Patronymic: tok}
	}
	if looksLikeSurname(tok) {
		return FullName{Last: tok}
	}
	return FullName{First: tok}
}

// classifyDouble — two tokens. If one of them is a patronymic, the other
// is classified as a given name (Russian full name without surname).
// Otherwise — "Last First" or "First Last".
func classifyDouble(a, b string) FullName {
	switch {
	case looksLikePatronymic(a):
		return FullName{First: b, Patronymic: a}
	case looksLikePatronymic(b):
		return FullName{First: a, Patronymic: b}
	}
	// Both are either names or surnames. Use the dictionary.
	aIsName := isDictName(a)
	bIsName := isDictName(b)
	switch {
	case aIsName && !bIsName:
		return FullName{First: a, Last: b}
	case bIsName && !aIsName:
		return FullName{Last: a, First: b}
	}
	// Default: "Last First" — the most common order in the database.
	return FullName{Last: a, First: b}
}

// classifyTriple — three tokens. The patronymic position determines the order:
// "Last First Patr" (Russian) or "First Patr Last" (Western).
func classifyTriple(a, b, c string) (FullName, error) {
	switch {
	case looksLikePatronymic(c):
		return FullName{Last: a, First: b, Patronymic: c}, nil
	case looksLikePatronymic(b):
		return FullName{First: a, Patronymic: b, Last: c}, nil
	case looksLikePatronymic(a):
		return FullName{}, fmt.Errorf("ParseFullName: patronymic cannot come first in %q", a+" "+b+" "+c)
	}
	// No patronymic among the three — there is no ordering signal. We
	// treat "Last First Patronymic-look-alike" as unworkable and return as
	// "Last First X" — the user will see the parse is ambiguous.
	return FullName{Last: a, First: b, Patronymic: c}, nil
}

// looksLikePatronymic — true for a token ending in a recognizable
// patronymic suffix or its case form.
func looksLikePatronymic(s string) bool {
	_, _, ok := inversePatronymic(s)
	return ok
}

// looksLikeSurname — true for a token with an informative surname ending
// (including dictionary lookup). It does not cover rare
// consonants/vowels for which one cannot tell "surname or given name"
// without context.
func looksLikeSurname(s string) bool {
	if isDictSurname(s) {
		return true
	}
	return genderFromSurnameHeuristic(s) != GenderUnknown
}

// isDictName — true if the dictionary contains a lemma tagged Name whose
// surface form matches s.
func isDictName(s string) bool {
	d, err := getDict()
	if err != nil {
		return false
	}
	for _, e := range d.Lookup(s) {
		l := &d.Lemmas[e.LemmaID]
		if hasTagStr(d, l.LemmaTags, "Name") {
			return true
		}
	}
	return false
}

// isDictSurname — true if the dictionary contains a lemma tagged Surn
// whose surface form matches s.
func isDictSurname(s string) bool {
	d, err := getDict()
	if err != nil {
		return false
	}
	for _, e := range d.Lookup(s) {
		l := &d.Lemmas[e.LemmaID]
		if hasTagStr(d, l.LemmaTags, "Surn") {
			return true
		}
	}
	return false
}
