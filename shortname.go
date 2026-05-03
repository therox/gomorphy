package gomorphy

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// ShortName returns the short form of a full name: surname plus initials
// for the given name and patronymic. Empty fields are omitted; surrounding
// whitespace inside fields is trimmed.
//
// Examples:
//
//	{Last: "Иванов",  First: "Иван", Patronymic: "Иванович"}  → "Иванов И. И."
//	{Last: "Иванова", First: "Анна", Patronymic: "Сергеевна"} → "Иванова А. С."
//	{Last: "Иванова", First: "Анна"}                          → "Иванова А."
//	{Last: "Иванов"}                                          → "Иванов"
//	{First: "Иван", Patronymic: "Иванович"}                   → "И. И."
//	{}                                                        → ""
//
// The initial is the first letter of the component, uppercased, followed
// by a period. A component whose first rune is not a letter contributes
// nothing.
//
// ShortName does not normalize case — surname and patronymic are copied as
// given. To get a short form in a non-nominative case, pass already
// inflected components (e.g. from DeclineFullName):
//
//	gen, _ := DeclineFullName(nom, Genitive)
//	ShortName(gen) // "Ивановой А. С."
func ShortName(name FullName) string {
	parts := make([]string, 0, 3)
	if last := strings.TrimSpace(name.Last); last != "" {
		parts = append(parts, last)
	}
	if init := initial(name.First); init != "" {
		parts = append(parts, init)
	}
	if init := initial(name.Patronymic); init != "" {
		parts = append(parts, init)
	}
	return strings.Join(parts, " ")
}

// initial returns the first letter of s, uppercased and followed by a
// period. Returns "" for an empty string or a leading non-letter rune.
func initial(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	r, _ := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError || !unicode.IsLetter(r) {
		return ""
	}
	return string(unicode.ToUpper(r)) + "."
}
