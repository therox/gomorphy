package gomorphy

// File: heuristic for inflecting given names that are not in the
// dictionary. Mirrors the surname/patronymic heuristics in surname.go and
// patronymic.go but tailored to first-name declension classes.
//
// Russian first names follow standard noun declension classes:
//   - consonant or -й / -ь: 2nd masculine declension (Иван, Игорь, Андрей)
//   - -а: 1st declension (Анна, Никита) — both genders use it
//   - -я / -ия: 1st declension (Илья, Мария)
//   - other final vowels: indeclinable foreign (Жозе, Дзё, Бруно)

import "strings"

// declineFirstHeuristic builds a first-name form by suffix/declension
// rules. The argument first must be lowercase. Returns a lowercase form
// and true when a pattern matches; ("", false) means "leave as is".
//
// Used as a fallback in declineFirst when the dictionary has no entry or
// the lemma's paradigm is incomplete.
func declineFirstHeuristic(first string, c Case, g Gender) (string, bool) {
	rs := []rune(first)
	n := len(rs)
	if n == 0 {
		return first, true
	}
	last := rs[n-1]

	// Indeclinable foreign: vowels other than -а/-я (Жозе, Дзё, Бруно, Жюли).
	if isVowel(last) && last != 'а' && last != 'я' {
		return first, true
	}

	if last == 'а' {
		// Both Анна (femn) and Никита (masc) use 1st declension.
		return decline1stA(rs, c), true
	}
	if last == 'я' {
		// Both Мария / Юлия (femn -ия) and Илья (masc) use 1st declension.
		return decline1stYa(rs, c), true
	}

	// Consonant or -й / -ь — masculine 2nd declension. Feminine names with
	// such endings in Russian are rare (foreign) and indeclinable.
	if g != Masculine {
		return first, true
	}
	return decline2ndMasc(rs, c), true
}

// inverseFirstNameHeuristic reduces a first-name form in an arbitrary case
// to the nominative. Returns (Nom_lowercase, gender, ok).
//
// Coverage is intentionally narrower than the forward direction — without
// a dictionary the inverse is far more ambiguous (Игоря/Александра,
// Кати/Юлии). Patterns handled:
//
//   - Masculine on a consonant: case endings -а / -у / -ом / -е strip to
//     the consonant base (Иванa → Иван, Иваном → Иван, Иване → Иван).
//   - Feminine -а stem: -ы / -е / -у / -ой strip to the -а stem
//     (Анны → Анна).
//   - Feminine -я / -ия stem: doubled -ии → -ия (Юлии → Юлия); -ю → -я
//     (Юлию → Юлия); -ей / -ёй → -я (Юлией → Юлия).
//
// Forms outside this set are returned as is — the caller must accept that
// the round-trip cannot be completed without a dictionary.
func inverseFirstNameHeuristic(form string, gHint Gender) (string, Gender, bool) {
	norm := strings.ToLower(form)
	rs := []rune(norm)
	n := len(rs)
	if n == 0 {
		return norm, GenderUnknown, false
	}
	last := rs[n-1]

	// Indeclinable foreign: vowels other than -а / -я.
	if isVowel(last) && last != 'а' && last != 'я' && last != 'е' &&
		last != 'у' && last != 'ы' && last != 'и' && last != 'ю' {
		return norm, gHint, true
	}

	if gHint == Masculine {
		// -ом Inst → consonant base (Иваном → Иван).
		if n >= 2 && rs[n-2] == 'о' && last == 'м' {
			return string(rs[:n-2]), Masculine, true
		}
		// -а Gen/Acc, -у Dat, -е Prep → strip → consonant base.
		// The -ь / -й classes (Игоря, Андрея) collide with bare-consonant
		// ones (Александра) and are NOT recognized here on purpose.
		switch last {
		case 'а', 'у', 'е':
			return string(rs[:n-1]), Masculine, true
		}
		// Already in Nom (consonant / -ь / -й), or unknown form.
		return norm, Masculine, true
	}

	if gHint == Feminine {
		// Already in -а / -я Nom (Анна, Юлия).
		if last == 'а' || last == 'я' {
			return norm, Feminine, true
		}
		// -ой Inst → -а (Анной → Анна).
		if n >= 2 && rs[n-2] == 'о' && last == 'й' {
			return string(rs[:n-2]) + "а", Feminine, true
		}
		// -ей / -ёй Inst → -я (Юлией → Юлия, Зоей → Зоя).
		if n >= 2 && (rs[n-2] == 'е' || rs[n-2] == 'ё') && last == 'й' {
			return string(rs[:n-2]) + "я", Feminine, true
		}
		switch last {
		case 'ы':
			// Анны → Анна.
			return string(rs[:n-1]) + "а", Feminine, true
		case 'е':
			// -е Dat/Prep — by default to -а (Анне → Анна). Кате → Катя is
			// missed; common Russian -я names are in the dictionary.
			return string(rs[:n-1]) + "а", Feminine, true
		case 'у':
			// Анну → Анна.
			return string(rs[:n-1]) + "а", Feminine, true
		case 'ю':
			// Юлию → Юлия, Катю → Катя.
			return string(rs[:n-1]) + "я", Feminine, true
		case 'и':
			// -ии Gen/Dat/Prep → -ия (Юлии → Юлия).
			if n >= 2 && rs[n-2] == 'и' {
				return string(rs[:n-1]) + "я", Feminine, true
			}
			// -ши / -жи / -чи / -щи / -ки / -ги / -хи Gen → -а (Маши → Маша).
			if n >= 2 && isHushOrKGH(rs[n-2]) {
				return string(rs[:n-1]) + "а", Feminine, true
			}
			// Other -и → -я (Кати → Катя).
			return string(rs[:n-1]) + "я", Feminine, true
		}
	}

	// Without a usable hint, no safe inverse.
	return norm, GenderUnknown, false
}
