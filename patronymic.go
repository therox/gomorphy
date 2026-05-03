package gomorphy

// File: heuristic for inflecting patronymics.
// Patronymics are inflected strictly regularly by the suffix of the base
// form. Used as a fallback when a patronymic is not in the dictionary or
// the desired form is missing in its paradigm.

import "strings"

// genderFromPatronymic determines gender from the patronymic suffix.
// Returns GenderUnknown if the suffix is not recognized.
func genderFromPatronymic(patr string) Gender {
	norm := strings.ToLower(patr)
	switch {
	case strings.HasSuffix(norm, "овна"),
		strings.HasSuffix(norm, "евна"),
		strings.HasSuffix(norm, "инична"),
		strings.HasSuffix(norm, "ична"):
		return Feminine
	case strings.HasSuffix(norm, "ович"),
		strings.HasSuffix(norm, "евич"),
		strings.HasSuffix(norm, "ьич"),
		strings.HasSuffix(norm, "ич"):
		return Masculine
	}
	return GenderUnknown
}

// declinePatronymicHeuristic builds a patronymic form from the rules.
// The argument patr must be lowercase. Returns a lowercase form together
// with true, or ("", false) when the pattern did not match.
//
// Masculine type (-ович/-евич/-ич/-ьич): drop "ич", append the case
// suffix from patrMascSuffix. Feminine type (-овна/-евна/-ична/-инична):
// drop "а", append the suffix from patrFemnSuffix.
func declinePatronymicHeuristic(patr string, c Case, g Gender) (string, bool) {
	rs := []rune(patr)
	n := len(rs)
	if n < 2 {
		return "", false
	}
	switch g {
	case Masculine:
		if rs[n-2] != 'и' || rs[n-1] != 'ч' {
			return "", false
		}
		base := string(rs[:n-2])
		suf := patrMascSuffix(c)
		if suf == "" {
			return "", false
		}
		return base + suf, true
	case Feminine:
		if rs[n-1] != 'а' {
			return "", false
		}
		base := string(rs[:n-1])
		suf := patrFemnSuffix(c)
		if suf == "" {
			return "", false
		}
		return base + suf, true
	}
	return "", false
}

// patrMascSuffix — endings of masculine patronymics in six cases.
// The animate accusative coincides with the genitive.
func patrMascSuffix(c Case) string {
	switch c {
	case Nominative:
		return "ич"
	case Genitive:
		return "ича"
	case Dative:
		return "ичу"
	case Accusative:
		return "ича"
	case Instrumental:
		return "ичем"
	case Prepositional:
		return "иче"
	}
	return ""
}

// patrFemnSuffix — endings of feminine patronymics in six cases.
func patrFemnSuffix(c Case) string {
	switch c {
	case Nominative:
		return "а"
	case Genitive:
		return "ы"
	case Dative:
		return "е"
	case Accusative:
		return "у"
	case Instrumental:
		return "ой"
	case Prepositional:
		return "е"
	}
	return ""
}

// inversePatronymic reduces a patronymic in an arbitrary case to the
// nominative. Returns (form_in_Nom_lowercase, gender, ok). When ok=false
// the input could not be recognized as a patronymic — the caller must
// leave the component as is.
//
// Algorithm: recognize the pair (stem, case ending).
//   - Masculine: stem ends in "ич", case endings
//     {"", "а", "у", "е", "ем"}.
//   - Feminine: stem ∈ {"иничн", "овн", "евн", "ичн"}, case endings
//     {"а", "ы", "е", "у", "ой"}.
//
// The "иничн" stem is checked BEFORE "ичн" — otherwise "Ильиничны" would
// be wrongly truncated to "Ильинича" instead of "Ильинична".
func inversePatronymic(s string) (string, Gender, bool) {
	norm := strings.ToLower(s)

	// Masculine: check case endings from longest to shortest.
	// The empty ending (Nom) is checked last so that longer
	// "ем"/"а"/"у"/"е" take priority.
	for _, end := range []string{"ем", "а", "у", "е", ""} {
		if strings.HasSuffix(norm, "ич"+end) {
			return strings.TrimSuffix(norm, end), Masculine, true
		}
	}

	// Feminine: outer loop is over stems (more specific ones first),
	// the inner loop is over case endings.
	for _, stem := range []string{"иничн", "овн", "евн", "ичн"} {
		for _, end := range []string{"ой", "а", "ы", "е", "у"} {
			if strings.HasSuffix(norm, stem+end) {
				return strings.TrimSuffix(norm, end) + "а", Feminine, true
			}
		}
	}

	return "", GenderUnknown, false
}
