// Package gomorphy provides a public API for inflecting Russian words
// and morphological analysis of word forms. The implementation is dictionary
// based (OpenCorpora).
package gomorphy

// POS is a part of speech.
type POS uint8

const (
	// POSUnknown — unknown or undefined part of speech.
	POSUnknown POS = iota
	// POSNoun — noun (including given names, patronymics, surnames).
	POSNoun
	// POSAdjf — adjective (full form).
	POSAdjf
)

// String returns the short string name of the part of speech (as in OpenCorpora).
func (p POS) String() string {
	switch p {
	case POSNoun:
		return "NOUN"
	case POSAdjf:
		return "ADJF"
	default:
		return "UNKN"
	}
}

// Case is a Russian grammatical case (the six basic ones).
type Case uint8

const (
	// CaseUnknown — undefined case.
	CaseUnknown Case = iota
	// Nominative — nominative (кто? что?).
	Nominative
	// Genitive — genitive (кого? чего?).
	Genitive
	// Dative — dative (кому? чему?).
	Dative
	// Accusative — accusative (кого? что?).
	Accusative
	// Instrumental — instrumental (кем? чем?).
	Instrumental
	// Prepositional — prepositional (о ком? о чём?).
	Prepositional
)

// String returns the short Russian grade-school abbreviation of the case:
// "им.п.", "род.п.", and so on.
func (c Case) String() string {
	switch c {
	case Nominative:
		return "им.п."
	case Genitive:
		return "род.п."
	case Dative:
		return "дат.п."
	case Accusative:
		return "вин.п."
	case Instrumental:
		return "тв.п."
	case Prepositional:
		return "пр.п."
	}
	return "?"
}

// Number is a grammatical number.
type Number uint8

const (
	// NumberUnknown — undefined number.
	NumberUnknown Number = iota
	// Singular — singular number.
	Singular
	// Plural — plural number.
	Plural
)

// String returns the short abbreviation of the number: "ед.ч." / "мн.ч.".
func (n Number) String() string {
	switch n {
	case Singular:
		return "ед.ч."
	case Plural:
		return "мн.ч."
	}
	return "?"
}

// Gender is a grammatical gender.
type Gender uint8

const (
	// GenderUnknown — gender is undefined (e.g. for pluralia tantum).
	GenderUnknown Gender = iota
	// Masculine — masculine gender.
	Masculine
	// Feminine — feminine gender.
	Feminine
	// Neuter — neuter gender.
	Neuter
	// Common — common gender (сирота, плакса).
	Common
)

// String returns the short abbreviation of the gender: "м.р.", "ж.р.", "ср.р.", "общ.р.".
func (g Gender) String() string {
	switch g {
	case Masculine:
		return "м.р."
	case Feminine:
		return "ж.р."
	case Neuter:
		return "ср.р."
	case Common:
		return "общ.р."
	}
	return "?"
}

// FullName is a person's full name (last / first / patronymic).
// Empty fields stay empty during inflection.
type FullName struct {
	// Last is the family name (surname).
	Last string
	// First is the given name.
	First string
	// Patronymic is the patronymic.
	Patronymic string
}

// Analysis is the result of morphological analysis of a single word form.
// For an ambiguous form Parse returns several Analysis entries.
type Analysis struct {
	// Lemma is the base form (nominative singular).
	Lemma string
	// POS is the part of speech.
	POS POS
	// Case is the case of the analyzed form.
	Case Case
	// Number is the number of the analyzed form.
	Number Number
	// Gender is the gender of the lemma (for NOUN) or the agreed gender of the
	// form (for ADJF).
	Gender Gender
	// Animate is animacy (relevant for nouns and adjectives).
	Animate bool
}
