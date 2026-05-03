package gomorphy

import "testing"

// TestShortName — table tests for ShortName covering complete and partial
// full names plus a few edge cases (lowercase input, leading whitespace,
// empty struct).
func TestShortName(t *testing.T) {
	cases := []struct {
		name string
		in   FullName
		want string
	}{
		{
			"masc full",
			FullName{Last: "Иванов", First: "Иван", Patronymic: "Иванович"},
			"Иванов И. И.",
		},
		{
			"femn full",
			FullName{Last: "Иванова", First: "Анна", Patronymic: "Сергеевна"},
			"Иванова А. С.",
		},
		{
			"last + first",
			FullName{Last: "Иванова", First: "Анна"},
			"Иванова А.",
		},
		{
			"last + patronymic",
			FullName{Last: "Иванов", Patronymic: "Иванович"},
			"Иванов И.",
		},
		{
			"last only",
			FullName{Last: "Иванов"},
			"Иванов",
		},
		{
			"first + patronymic, no last",
			FullName{First: "Иван", Patronymic: "Иванович"},
			"И. И.",
		},
		{
			"first only",
			FullName{First: "Иван"},
			"И.",
		},
		{
			"patronymic only",
			FullName{Patronymic: "Иванович"},
			"И.",
		},
		{
			"lowercase first → uppercased initial, surname kept as is",
			FullName{Last: "иванов", First: "иван"},
			"иванов И.",
		},
		{
			"surrounding whitespace is trimmed",
			FullName{Last: "  Иванов ", First: " Иван ", Patronymic: " Иванович "},
			"Иванов И. И.",
		},
		{
			"empty struct",
			FullName{},
			"",
		},
		{
			"latin input",
			FullName{Last: "Smith", First: "John", Patronymic: "Edward"},
			"Smith J. E.",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ShortName(tc.in)
			if got != tc.want {
				t.Fatalf("ShortName(%+v) = %q, expected %q", tc.in, got, tc.want)
			}
		})
	}
}

// TestShortNamePreservesCase — ShortName itself does not normalize case;
// it copies surname and patronymic as given and only uppercases the
// initials. Pass already-inflected components (e.g. straight from
// DeclineFullName) to get the short form in that case.
func TestShortNamePreservesCase(t *testing.T) {
	cases := []struct {
		name string
		in   FullName
		want string
	}{
		{
			"genitive",
			FullName{Last: "Ивановой", First: "Анны", Patronymic: "Сергеевны"},
			"Ивановой А. С.",
		},
		{
			"dative",
			FullName{Last: "Ивановой", First: "Анне", Patronymic: "Сергеевне"},
			"Ивановой А. С.",
		},
		{
			"instrumental",
			FullName{Last: "Ивановой", First: "Анной", Patronymic: "Сергеевной"},
			"Ивановой А. С.",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ShortName(tc.in)
			if got != tc.want {
				t.Fatalf("ShortName(%+v) = %q, expected %q", tc.in, got, tc.want)
			}
		})
	}
}

// TestShortNameAfterParseFullName — common usage chain: a free-form string
// is parsed and then reduced to the short form.
func TestShortNameAfterParseFullName(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Иванов Иван Иванович", "Иванов И. И."},
		{"Иван Иванович Иванов", "Иванов И. И."},
		{"Ивановой Анне Сергеевне", "Иванова А. С."},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			fn, err := ParseFullName(tc.in)
			if err != nil {
				t.Fatalf("ParseFullName(%q) error: %v", tc.in, err)
			}
			got := ShortName(fn)
			if got != tc.want {
				t.Fatalf("ShortName(ParseFullName(%q)) = %q, expected %q",
					tc.in, got, tc.want)
			}
		})
	}
}
