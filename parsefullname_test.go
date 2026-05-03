package gomorphy

import (
	"strings"
	"testing"
)

// TestToNominativeRussianMascAll — a complete masculine full name in
// arbitrary cases is reduced to Nom.
func TestToNominativeRussianMascAll(t *testing.T) {
	wantNom := FullName{Last: "Иванов", First: "Иван", Patronymic: "Иванович"}
	cases := []FullName{
		// Nom (self-check).
		{Last: "Иванов", First: "Иван", Patronymic: "Иванович"},
		// Gen.
		{Last: "Иванова", First: "Ивана", Patronymic: "Ивановича"},
		// Dat.
		{Last: "Иванову", First: "Ивану", Patronymic: "Ивановичу"},
		// Acc.
		{Last: "Иванова", First: "Ивана", Patronymic: "Ивановича"},
		// Inst.
		{Last: "Ивановым", First: "Иваном", Patronymic: "Ивановичем"},
		// Prep.
		{Last: "Иванове", First: "Иване", Patronymic: "Ивановиче"},
	}
	for _, in := range cases {
		t.Run(joinFullName(in), func(t *testing.T) {
			got, err := ToNominative(in)
			if err != nil {
				t.Fatalf("ToNominative(%v) error: %v", in, err)
			}
			if got != wantNom {
				t.Fatalf("ToNominative(%v) = %v, expected %v", in, got, wantNom)
			}
		})
	}
}

// TestToNominativeRussianFemnAll — a feminine full name in arbitrary cases.
func TestToNominativeRussianFemnAll(t *testing.T) {
	wantNom := FullName{Last: "Иванова", First: "Анна", Patronymic: "Сергеевна"}
	cases := []FullName{
		{Last: "Иванова", First: "Анна", Patronymic: "Сергеевна"},
		{Last: "Ивановой", First: "Анны", Patronymic: "Сергеевны"},
		{Last: "Ивановой", First: "Анне", Patronymic: "Сергеевне"},
		{Last: "Иванову", First: "Анну", Patronymic: "Сергеевну"},
		{Last: "Ивановой", First: "Анной", Patronymic: "Сергеевной"},
		{Last: "Ивановой", First: "Анне", Patronymic: "Сергеевне"},
	}
	for _, in := range cases {
		t.Run(joinFullName(in), func(t *testing.T) {
			got, err := ToNominative(in)
			if err != nil {
				t.Fatalf("ToNominative(%v) error: %v", in, err)
			}
			if got != wantNom {
				t.Fatalf("ToNominative(%v) = %v, expected %v", in, got, wantNom)
			}
		})
	}
}

// TestToNominativeRoundTrip — symmetry check:
// ToNominative(DeclineFullName(nom, c)) returns the same nom for all cases.
func TestToNominativeRoundTrip(t *testing.T) {
	noms := []FullName{
		{Last: "Иванов", First: "Иван", Patronymic: "Иванович"},
		{Last: "Иванова", First: "Анна", Patronymic: "Сергеевна"},
		{Last: "Достоевский", First: "Иван"},
		{Last: "Достоевская", First: "Анна"},
		{Last: "Иванов", First: "Иван"},
		{Last: "Иванова", First: "Анна"},
	}
	for _, nom := range noms {
		for _, c := range allCases() {
			t.Run(joinFullName(nom)+"/"+caseTag(c), func(t *testing.T) {
				declined, err := DeclineFullName(nom, c)
				if err != nil {
					t.Fatalf("DeclineFullName(%v, %v) error: %v", nom, c, err)
				}
				got, err := ToNominative(declined)
				if err != nil {
					t.Fatalf("ToNominative(%v) error: %v", declined, err)
				}
				if got != nom {
					t.Fatalf("round-trip %v → %v → %v, expected %v",
						nom, declined, got, nom)
				}
			})
		}
	}
}

// TestToNominativeIndecl — indeclinable Дюма (Surn+Fixd) — Nom = Дюма.
func TestToNominativeIndecl(t *testing.T) {
	in := FullName{Last: "Дюма", First: "Александр"}
	got, err := ToNominative(in)
	if err != nil {
		t.Fatalf("ToNominative(%v) error: %v", in, err)
	}
	want := FullName{Last: "Дюма", First: "Александр"}
	if got != want {
		t.Fatalf("ToNominative(%v) = %v, expected %v", in, got, want)
	}
}

// TestToNominativeEmpty — an empty FullName → error.
func TestToNominativeEmpty(t *testing.T) {
	if _, err := ToNominative(FullName{}); err == nil {
		t.Fatal("ToNominative(FullName{}) must return an error")
	}
}

// TestParseFullNameRussianTriple — the Russian "Last First Patr" order.
func TestParseFullNameRussianTriple(t *testing.T) {
	want := FullName{Last: "Иванов", First: "Иван", Patronymic: "Иванович"}
	cases := []string{
		"Иванов Иван Иванович",
		"Иванова Ивана Ивановича",   // Gen
		"Иванову Ивану Ивановичу",   // Dat
		"Ивановым Иваном Ивановичем", // Inst
		"Иванове Иване Ивановиче",   // Prep
	}
	for _, s := range cases {
		t.Run(s, func(t *testing.T) {
			got, err := ParseFullName(s)
			if err != nil {
				t.Fatalf("ParseFullName(%q) error: %v", s, err)
			}
			if got != want {
				t.Fatalf("ParseFullName(%q) = %v, expected %v", s, got, want)
			}
		})
	}
}

// TestParseFullNameWesternTriple — the Western "First Patr Last" order.
func TestParseFullNameWesternTriple(t *testing.T) {
	want := FullName{Last: "Иванов", First: "Иван", Patronymic: "Иванович"}
	cases := []string{
		"Иван Иванович Иванов",
		"Ивана Ивановича Иванова",
		"Ивану Ивановичу Иванову",
	}
	for _, s := range cases {
		t.Run(s, func(t *testing.T) {
			got, err := ParseFullName(s)
			if err != nil {
				t.Fatalf("ParseFullName(%q) error: %v", s, err)
			}
			if got != want {
				t.Fatalf("ParseFullName(%q) = %v, expected %v", s, got, want)
			}
		})
	}
}

// TestParseFullNameDouble — two tokens: surname + given name in any order.
func TestParseFullNameDouble(t *testing.T) {
	want := FullName{Last: "Иванов", First: "Иван"}
	cases := []string{
		"Иванов Иван",
		"Иван Иванов",   // Иван → Name from the dictionary, Иванов → Last
		"Иванова Ивана", // both in Gen
	}
	for _, s := range cases {
		t.Run(s, func(t *testing.T) {
			got, err := ParseFullName(s)
			if err != nil {
				t.Fatalf("ParseFullName(%q) error: %v", s, err)
			}
			if got != want {
				t.Fatalf("ParseFullName(%q) = %v, expected %v", s, got, want)
			}
		})
	}
}

// TestParseFullNameNameAndPatr — given name + patronymic, no surname.
func TestParseFullNameNameAndPatr(t *testing.T) {
	cases := []struct {
		in   string
		want FullName
	}{
		{"Иван Иванович", FullName{First: "Иван", Patronymic: "Иванович"}},
		{"Анна Сергеевна", FullName{First: "Анна", Patronymic: "Сергеевна"}},
		{"Ивана Ивановича", FullName{First: "Иван", Patronymic: "Иванович"}},
		{"Анной Сергеевной", FullName{First: "Анна", Patronymic: "Сергеевна"}},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got, err := ParseFullName(tc.in)
			if err != nil {
				t.Fatalf("ParseFullName(%q) error: %v", tc.in, err)
			}
			if got != tc.want {
				t.Fatalf("ParseFullName(%q) = %v, expected %v", tc.in, got, tc.want)
			}
		})
	}
}

// TestParseFullNameSingle — a single token is recognized by the heuristic.
func TestParseFullNameSingle(t *testing.T) {
	cases := []struct {
		in   string
		want FullName
	}{
		{"Иванович", FullName{Patronymic: "Иванович"}},
		{"Ивановича", FullName{Patronymic: "Иванович"}},
		{"Сергеевне", FullName{Patronymic: "Сергеевна"}},
		{"Иванов", FullName{Last: "Иванов"}},
		{"Иванова", FullName{Last: "Иванова"}},
		{"Достоевский", FullName{Last: "Достоевский"}},
		{"Иван", FullName{First: "Иван"}},
		{"Анна", FullName{First: "Анна"}},
	}
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got, err := ParseFullName(tc.in)
			if err != nil {
				t.Fatalf("ParseFullName(%q) error: %v", tc.in, err)
			}
			if got != tc.want {
				t.Fatalf("ParseFullName(%q) = %v, expected %v", tc.in, got, tc.want)
			}
		})
	}
}

// TestParseFullNameRoundTrip — the ParseFullName + DeclineFullName chain
// gives a consistent result for all cases.
func TestParseFullNameRoundTrip(t *testing.T) {
	// Feed it a dative form, demand the nominative, then inflect into the genitive.
	got, err := ParseFullName("Ивановой Анне Сергеевне")
	if err != nil {
		t.Fatalf("ParseFullName error: %v", err)
	}
	wantNom := FullName{Last: "Иванова", First: "Анна", Patronymic: "Сергеевна"}
	if got != wantNom {
		t.Fatalf("ParseFullName(dat) = %v, expected %v", got, wantNom)
	}

	gen, err := DeclineFullName(got, Genitive)
	if err != nil {
		t.Fatalf("DeclineFullName error: %v", err)
	}
	wantGen := FullName{Last: "Ивановой", First: "Анны", Patronymic: "Сергеевны"}
	if gen != wantGen {
		t.Fatalf("DeclineFullName(%v, Gen) = %v, expected %v", got, gen, wantGen)
	}
}

// TestParseFullNameErrors — an empty string and ≥4 tokens are rejected.
func TestParseFullNameErrors(t *testing.T) {
	cases := []string{
		"",
		"   ",
		"А Б В Г",                        // 4 tokens
		"Иванов Петров Сидоров Кузнецов", // 4 tokens
	}
	for _, s := range cases {
		t.Run(s, func(t *testing.T) {
			_, err := ParseFullName(s)
			if err == nil {
				t.Fatalf("ParseFullName(%q) — expected an error", s)
			}
		})
	}
}

// TestParseFullNamePatronymicFirstError — a patronymic at position 0 is
// an error (it is neither the Russian nor the Western order).
func TestParseFullNamePatronymicFirstError(t *testing.T) {
	_, err := ParseFullName("Иванович Иван Иванов")
	if err == nil {
		t.Fatal("ParseFullName with a patronymic in the first position must return an error")
	}
	if !strings.Contains(err.Error(), "patronymic cannot come first") {
		t.Fatalf("unexpected error message: %v", err)
	}
}
