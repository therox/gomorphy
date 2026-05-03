package gomorphy

// DeclineFullName inflects a full name as a whole into case c.
// Each component is inflected by its own rules (Name/Patr/Surn).
// Empty fields of name stay empty. Indeclinable names are returned as is.
//
// Algorithm and heuristics: see docs/DESIGN.md (Phase 3).
func DeclineFullName(name FullName, c Case) (FullName, error) {
	return declineFullName(name, c)
}

// ToNominative reduces a full name in an arbitrary case to the nominative.
// Each component is processed independently: first the dictionary is tried
// (Lookup → lemma), then a reverse heuristic by suffix. Empty fields stay
// empty. Out-of-dictionary components that do not fit the heuristic
// (non-dictionary given names, rare surname patterns) are returned as is.
//
// Used as a counterpart to DeclineFullName: round-trip
//
//	in   := FullName{Last: "Ивановой", First: "Анне", Patronymic: "Сергеевне"}
//	nom, _ := ToNominative(in)                 // → Иванова Анна Сергеевна
//	abl, _ := DeclineFullName(nom, Instrumental) // → Ивановой Анной Сергеевной
//
// Algorithm and tables: see docs/DESIGN.md (Phase 4).
func ToNominative(name FullName) (FullName, error) {
	return toNominativeImpl(name)
}

// ParseFullName splits the string s into components (Last/First/Patronymic)
// and reduces each one to the nominative case.
//
// Supported formats:
//   - "Иванов Иван Иванович" — Russian order (Last First Patr)
//   - "Иван Иванович Иванов" — Western order (First Patr Last)
//   - "Иванов Иван" / "Иван Иванов" — two tokens
//   - "Иван Ивановна" / "Иван Иванович" — given name + patronymic
//   - "Иванов" / "Иванович" / "Иван" — single token
//
// The order is inferred from the patronymic position; for two tokens
// without a patronymic the dictionary (Name tag) is additionally used to
// determine which one is the given name.
//
// Returns an error for an empty string or ≥4 tokens.
func ParseFullName(s string) (FullName, error) {
	return parseFullNameImpl(s)
}
