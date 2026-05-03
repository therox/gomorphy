package gomorphy

// File: heuristic for inflecting out-of-dictionary surnames by suffix.
// Used as a fallback when a surname is not found in the dictionary or the
// lemma has no required form. Coverage — see the table in
// docs/DESIGN.md (Phase 3).

import "strings"

// genderFromSurnameHeuristic determines gender by an informative surname
// ending. For uninformative endings (Гайдай, Скрипка, Долгих, Шевченко)
// it returns GenderUnknown — the gender must be determined from another
// source (patronymic or given name).
func genderFromSurnameHeuristic(last string) Gender {
	norm := strings.ToLower(last)
	switch {
	case strings.HasSuffix(norm, "ова"),
		strings.HasSuffix(norm, "ева"),
		strings.HasSuffix(norm, "ёва"),
		strings.HasSuffix(norm, "ина"),
		strings.HasSuffix(norm, "ына"),
		strings.HasSuffix(norm, "ская"),
		strings.HasSuffix(norm, "цкая"):
		return Feminine
	case strings.HasSuffix(norm, "ов"),
		strings.HasSuffix(norm, "ев"),
		strings.HasSuffix(norm, "ёв"),
		strings.HasSuffix(norm, "ин"),
		strings.HasSuffix(norm, "ын"),
		strings.HasSuffix(norm, "ский"),
		strings.HasSuffix(norm, "цкий"):
		return Masculine
	}
	return GenderUnknown
}

// declineSurnameHeuristic builds a surname form from the heuristic table.
// The argument last must be lowercase, c is the desired case, g is the gender.
// Returns a lowercase form and true; ("", false) if the surname does not
// match any pattern or the gender did not match the pattern.
//
// Indeclinable patterns (-ых/-их, vowels other than -а/-я, etc.) return
// the original string with true — that is a successful "do not inflect"
// outcome.
func declineSurnameHeuristic(last string, c Case, g Gender) (string, bool) {
	rs := []rune(last)
	n := len(rs)
	if n == 0 {
		return last, true
	}

	// Indeclinable: -ых, -их (Долгих, Седых) — common gender.
	if hasRuneSuffix(rs, "ых") || hasRuneSuffix(rs, "их") {
		return last, true
	}

	// Adjectival -ский / -цкий — masculine.
	if hasRuneSuffix(rs, "ский") || hasRuneSuffix(rs, "цкий") {
		if g != Masculine {
			return "", false
		}
		return adjMascDecline(rs, c), true
	}
	// Adjectival -ская / -цкая — feminine.
	if hasRuneSuffix(rs, "ская") || hasRuneSuffix(rs, "цкая") {
		if g != Feminine {
			return "", false
		}
		return adjFemnDecline(rs, c), true
	}

	// Possessive feminine: -ова/-ева/-ёва/-ина/-ына.
	if hasRuneSuffix(rs, "ова") || hasRuneSuffix(rs, "ева") ||
		hasRuneSuffix(rs, "ёва") || hasRuneSuffix(rs, "ина") ||
		hasRuneSuffix(rs, "ына") {
		if g != Feminine {
			return "", false
		}
		base := string(rs[:n-1])
		return base + possFemnSuffix(c), true
	}
	// Possessive masculine: -ов/-ев/-ёв/-ин/-ын.
	if hasRuneSuffix(rs, "ов") || hasRuneSuffix(rs, "ев") ||
		hasRuneSuffix(rs, "ёв") || hasRuneSuffix(rs, "ин") ||
		hasRuneSuffix(rs, "ын") {
		if g != Masculine {
			return "", false
		}
		return last + possMascSuffix(c), true
	}

	// General adjectival: -ой/-ый/-ий (M), -ая/-яя (F).
	if hasRuneSuffix(rs, "ой") || hasRuneSuffix(rs, "ый") || hasRuneSuffix(rs, "ий") {
		if g != Masculine {
			return "", false
		}
		return adjMascDecline(rs, c), true
	}
	if hasRuneSuffix(rs, "ая") || hasRuneSuffix(rs, "яя") {
		if g != Feminine {
			return "", false
		}
		return adjFemnDecline(rs, c), true
	}

	// -а / -я → 1st declension (Скрипка, Заря).
	// NOTE: this includes the Дюма class — surnames with a stressed -а
	// that are in fact indeclinable. The heuristic does not detect stress
	// and inflects them (Дюма → Дюмы). For correct behavior add such
	// surnames to the dictionary with a Fixd tag. See docs/DESIGN.md §9.3.
	if rs[n-1] == 'а' {
		return decline1stA(rs, c), true
	}
	if rs[n-1] == 'я' {
		return decline1stYa(rs, c), true
	}

	// Other vowels (-о/-е/-у/-ю/-и/-ы/-э/-ё) — indeclinable.
	if isVowel(rs[n-1]) {
		return last, true
	}

	// Consonant remains — inflected only for masculine
	// (Гайдай → Гайдая, Шевчук → Шевчука). For feminine it is indeclinable.
	if g != Masculine {
		return last, true
	}
	return decline2ndMasc(rs, c), true
}

// adjMascDecline inflects a masculine adjectival surname (Толстой,
// Достоевский, Тверской, Горький). Takes a slice of lowercase runes.
//
// Spelling rule: after к/г/х/ж/ш/щ/ч/ц the instrumental ending uses "и"
// instead of "ы" (Достоевским, not Достоевскым).
func adjMascDecline(rs []rune, c Case) string {
	n := len(rs)
	if n < 2 {
		return string(rs)
	}
	nomEnd := string([]rune{rs[n-2], rs[n-1]})
	base := rs[:n-2]
	var lastBase rune
	if len(base) > 0 {
		lastBase = base[len(base)-1]
	}
	ablt := "ым"
	if isHushOrKGH(lastBase) {
		ablt = "им"
	}
	var ending string
	switch c {
	case Nominative:
		ending = nomEnd
	case Genitive, Accusative:
		ending = "ого"
	case Dative:
		ending = "ому"
	case Instrumental:
		ending = ablt
	case Prepositional:
		ending = "ом"
	}
	return string(base) + ending
}

// adjFemnDecline inflects a feminine adjectival surname (Толстая,
// Достоевская). The soft type (-яя) uses the endings -ей/-юю; the hard
// type uses -ой/-ую.
func adjFemnDecline(rs []rune, c Case) string {
	n := len(rs)
	if n < 2 {
		return string(rs)
	}
	nomEnd := string([]rune{rs[n-2], rs[n-1]})
	base := rs[:n-2]
	soft := nomEnd == "яя"
	var ending string
	switch c {
	case Nominative:
		ending = nomEnd
	case Genitive, Dative, Instrumental, Prepositional:
		if soft {
			ending = "ей"
		} else {
			ending = "ой"
		}
	case Accusative:
		if soft {
			ending = "юю"
		} else {
			ending = "ую"
		}
	}
	return string(base) + ending
}

// possMascSuffix — endings for masculine possessive surnames
// (-ов/-ев/-ин/...). The animate accusative coincides with the genitive.
func possMascSuffix(c Case) string {
	switch c {
	case Nominative:
		return ""
	case Genitive, Accusative:
		return "а"
	case Dative:
		return "у"
	case Instrumental:
		return "ым"
	case Prepositional:
		return "е"
	}
	return ""
}

// possFemnSuffix — endings for feminine possessive surnames
// (-ова/-ева/-ина/...).
func possFemnSuffix(c Case) string {
	switch c {
	case Nominative:
		return "а"
	case Genitive, Dative, Instrumental, Prepositional:
		return "ой"
	case Accusative:
		return "у"
	}
	return ""
}

// decline1stA inflects a noun-style surname ending in -а (Скрипка, Гора).
// After hushing and back-lingual consonants the genitive ending is "и"
// instead of "ы".
func decline1stA(rs []rune, c Case) string {
	n := len(rs)
	base := rs[:n-1]
	var lastBase rune
	if len(base) > 0 {
		lastBase = base[len(base)-1]
	}
	gent := "ы"
	if isHushOrKGH(lastBase) {
		gent = "и"
	}
	var ending string
	switch c {
	case Nominative:
		ending = "а"
	case Genitive:
		ending = gent
	case Dative, Prepositional:
		ending = "е"
	case Accusative:
		ending = "у"
	case Instrumental:
		ending = "ой"
	}
	return string(base) + ending
}

// decline1stYa inflects a noun-style surname ending in -я (Заря) and the
// -ия subtype. The -ия subtype has a special paradigm: dat./prep. → -ии
// (like "Мария").
func decline1stYa(rs []rune, c Case) string {
	n := len(rs)
	base := rs[:n-1]
	iya := n >= 2 && rs[n-2] == 'и'
	var ending string
	if iya {
		switch c {
		case Nominative:
			ending = "я"
		case Genitive, Dative, Prepositional:
			ending = "и"
		case Accusative:
			ending = "ю"
		case Instrumental:
			ending = "ей"
		}
	} else {
		switch c {
		case Nominative:
			ending = "я"
		case Genitive:
			ending = "и"
		case Dative, Prepositional:
			ending = "е"
		case Accusative:
			ending = "ю"
		case Instrumental:
			ending = "ей"
		}
	}
	return string(base) + ending
}

// decline2ndMasc inflects a masculine 2nd-declension noun-style surname:
// ending in a consonant (Шевчук), in -й (Гайдай), or in -ь (Соболь).
// After hushing consonants and ц the instrumental ending is "ем" rather
// than "ом".
func decline2ndMasc(rs []rune, c Case) string {
	n := len(rs)
	last := rs[n-1]
	if last == 'й' || last == 'ь' {
		base := rs[:n-1]
		var ending string
		switch c {
		case Nominative:
			ending = string(last)
		case Genitive, Accusative:
			ending = "я"
		case Dative:
			ending = "ю"
		case Instrumental:
			ending = "ем"
		case Prepositional:
			ending = "е"
		}
		return string(base) + ending
	}
	hush := last == 'ж' || last == 'ч' || last == 'ш' || last == 'щ' || last == 'ц'
	var ending string
	switch c {
	case Nominative:
		ending = ""
	case Genitive, Accusative:
		ending = "а"
	case Dative:
		ending = "у"
	case Instrumental:
		if hush {
			ending = "ем"
		} else {
			ending = "ом"
		}
	case Prepositional:
		ending = "е"
	}
	return string(rs) + ending
}

// hasRuneSuffix checks whether the trailing runes of rs match suffix.
// Works on []rune to avoid relying on byte boundaries of Cyrillic.
func hasRuneSuffix(rs []rune, suffix string) bool {
	sr := []rune(suffix)
	if len(rs) < len(sr) {
		return false
	}
	off := len(rs) - len(sr)
	for i, r := range sr {
		if rs[off+i] != r {
			return false
		}
	}
	return true
}

// isHushOrKGH — true for hushing and back-lingual consonants
// (к/г/х/ж/ш/щ/ч/ц). Used by the rules "after hushing/к-г-х use 'и'
// instead of 'ы'".
func isHushOrKGH(r rune) bool {
	switch r {
	case 'к', 'г', 'х', 'ж', 'ш', 'щ', 'ч', 'ц':
		return true
	}
	return false
}

// isVowel — true for Russian vowels.
func isVowel(r rune) bool {
	switch r {
	case 'а', 'я', 'о', 'ё', 'е', 'и', 'ы', 'у', 'ю', 'э':
		return true
	}
	return false
}

// inverseSurnameHeuristic reduces a surname in an arbitrary case to the
// nominative. Returns (Nom_lowercase, detected_gender, ok). When ok=false
// the input could not be recognized as a surname — the caller must leave
// the component as is.
//
// gHint — a gender hint from other components (patronymic, given name).
// Needed to resolve ambiguous forms:
//   - "Иванова": F Nom ("Иванова") or M Gen ("Иванов").
//   - "Иванову": F Acc or M Dat.
//   - "Толстой": M Nom or F Gen/Dat/Inst/Prep (from "Толстая").
//
// If gHint=GenderUnknown, F Nom is chosen for ambiguous cases (the more
// common interpretation for an unfamiliar surname).
//
// The returned gender is the one consistent with the recognized pattern
// and the hint.
func inverseSurnameHeuristic(last string, gHint Gender) (string, Gender, bool) {
	norm := strings.ToLower(last)

	// Indeclinable: -ых/-их.
	if strings.HasSuffix(norm, "ых") || strings.HasSuffix(norm, "их") {
		return norm, GenderUnknown, true
	}

	// Adjectival -ская/-цкая (F): endings -ая/-ой/-ую.
	for _, stem := range []string{"ск", "цк"} {
		for _, end := range []string{"ой", "ую", "ая"} {
			if strings.HasSuffix(norm, stem+end) {
				return strings.TrimSuffix(norm, end) + "ая", Feminine, true
			}
		}
	}

	// Adjectival -ский/-цкий (M): endings -ий/-ого/-ому/-ом/-им.
	for _, stem := range []string{"ск", "цк"} {
		for _, end := range []string{"ого", "ому", "им", "ом", "ий"} {
			if strings.HasSuffix(norm, stem+end) {
				return strings.TrimSuffix(norm, end) + "ий", Masculine, true
			}
		}
	}

	// Possessive F with unambiguous endings:
	// -овой/-евой/-ёвой/-иной/-ыной. (-ой here is unambiguously F, unlike
	// -ова / -ову which are ambiguous.)
	for _, stem := range []string{"ов", "ев", "ёв", "ин", "ын"} {
		if strings.HasSuffix(norm, stem+"ой") {
			return strings.TrimSuffix(norm, "ой") + "а", Feminine, true
		}
	}

	// Possessive F with ambiguous -а/-у: F Nom (-ова) or M Gen/Acc (-ова),
	// F Acc (-ову) or M Dat (-ову). Decide using gHint.
	for _, stem := range []string{"ов", "ев", "ёв", "ин", "ын"} {
		for _, end := range []string{"а", "у"} {
			if strings.HasSuffix(norm, stem+end) {
				base := strings.TrimSuffix(norm, end)
				if gHint == Masculine {
					return base, Masculine, true
				}
				return base + "а", Feminine, true
			}
		}
	}

	// Possessive M with unambiguous endings:
	// -овым/-евым/-ёвым/-иным/-ыным (Inst).
	for _, stem := range []string{"ов", "ев", "ёв", "ин", "ын"} {
		if strings.HasSuffix(norm, stem+"ым") {
			return strings.TrimSuffix(norm, "ым"), Masculine, true
		}
	}

	// Possessive M Prep: -ове/-еве/-ёве/-ине/-ыне.
	for _, stem := range []string{"ов", "ев", "ёв", "ин", "ын"} {
		if strings.HasSuffix(norm, stem+"е") {
			return strings.TrimSuffix(norm, "е"), Masculine, true
		}
	}

	// Possessive M Nom: -ов/-ев/-ёв/-ин/-ын.
	for _, stem := range []string{"ов", "ев", "ёв", "ин", "ын"} {
		if strings.HasSuffix(norm, stem) {
			return norm, Masculine, true
		}
	}

	// General adjectival: -ая/-яя F Acc → -ую/-юю.
	if strings.HasSuffix(norm, "ую") {
		return strings.TrimSuffix(norm, "ую") + "ая", Feminine, true
	}
	if strings.HasSuffix(norm, "юю") {
		return strings.TrimSuffix(norm, "юю") + "яя", Feminine, true
	}
	// -ая/-яя F Nom (as is).
	if strings.HasSuffix(norm, "ая") || strings.HasSuffix(norm, "яя") {
		return norm, Feminine, true
	}

	// -ой M (Толстой) or F gen/dat/inst/prep from -ая (Толстой = "к Толстой").
	// Decided by gHint.
	if strings.HasSuffix(norm, "ой") {
		if gHint == Feminine {
			return strings.TrimSuffix(norm, "ой") + "ая", Feminine, true
		}
		return norm, Masculine, true
	}

	// Adjectival M -ого/-ому/-ым/-им/-ом — Nom -ой/-ый/-ий depending on the
	// stem. After к/г/х/ц → -ий (Достоевский, Горький); otherwise → -ой
	// (Толстой).
	for _, end := range []string{"ого", "ому", "ым", "им", "ом"} {
		if strings.HasSuffix(norm, end) {
			base := strings.TrimSuffix(norm, end)
			return base + adjMascNomEnding(base), Masculine, true
		}
	}

	// -ый/-ий M Nom (as is).
	if strings.HasSuffix(norm, "ый") || strings.HasSuffix(norm, "ий") {
		return norm, Masculine, true
	}

	// 1st declension -а/-я: different case endings → restore -а/-я.
	rs := []rune(norm)
	n := len(rs)
	if n > 0 {
		last := rs[n-1]
		// -ы/-и (Gen) → -а; -е (Dat/Prep) → -а; -у (Acc) → -а.
		if last == 'ы' || last == 'и' || last == 'е' || last == 'у' {
			return string(rs[:n-1]) + "а", GenderUnknown, true
		}
		// -ю (Acc from -я) → -я.
		if last == 'ю' {
			return string(rs[:n-1]) + "я", GenderUnknown, true
		}
		// -а/-я (Nom) — as is.
		if last == 'а' || last == 'я' {
			return norm, GenderUnknown, true
		}
		// Other vowels — indeclinable.
		if isVowel(last) {
			return norm, GenderUnknown, true
		}
	}

	// Consonant (including -й/-ь) — 2nd declension M, Nom = as is.
	// (Case forms -я/-ю/-ем/-е cannot be reversed unambiguously: for two
	// distinct lemmas like "Гайдай" and "Гайдая" you cannot restore
	// correctly without a dictionary, but the input token in Nom is
	// already correct.)
	return norm, Masculine, true
}

// adjMascNomEnding chooses the Nom ending for a masculine adjectival
// surname based on the last letter of the stem. After к/г/х/ц — "ий"
// (Достоевский, Горький), otherwise — "ой" (Толстой). The "-ый" branch
// (Белый) is intentionally omitted: -ый is rare in surnames, and -ой
// gives the right result for the overwhelming majority of real
// variants.
func adjMascNomEnding(stem string) string {
	rs := []rune(stem)
	if len(rs) == 0 {
		return "ий"
	}
	last := rs[len(rs)-1]
	switch last {
	case 'к', 'г', 'х', 'ц':
		return "ий"
	}
	return "ой"
}
