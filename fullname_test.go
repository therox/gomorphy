package gomorphy

import (
	"strings"
	"testing"
)

// fullCase — one expectation row: input full name, case, expected full name.
type fullCase struct {
	in   FullName
	c    Case
	want FullName
}

// TestDeclineFullNameMascAll — a complete masculine full name in six cases.
// "Иванов Иван Иванович". Иванов is in the dictionary up to Inst, Иванович —
// up to Dat; for the remaining cases the heuristic kicks in and matches
// the expected form.
func TestDeclineFullNameMascAll(t *testing.T) {
	in := FullName{Last: "Иванов", First: "Иван", Patronymic: "Иванович"}
	cases := []fullCase{
		{in, Nominative, FullName{Last: "Иванов", First: "Иван", Patronymic: "Иванович"}},
		{in, Genitive, FullName{Last: "Иванова", First: "Ивана", Patronymic: "Ивановича"}},
		{in, Dative, FullName{Last: "Иванову", First: "Ивану", Patronymic: "Ивановичу"}},
		{in, Accusative, FullName{Last: "Иванова", First: "Ивана", Patronymic: "Ивановича"}},
		{in, Instrumental, FullName{Last: "Ивановым", First: "Иваном", Patronymic: "Ивановичем"}},
		{in, Prepositional, FullName{Last: "Иванове", First: "Иване", Patronymic: "Ивановиче"}},
	}
	runFullnameCases(t, cases)
}

// TestDeclineFullNameFemnAll — a complete feminine full name in six cases.
// "Иванова Анна Сергеевна". Сергеевна is in the dictionary only in Nom —
// for the remaining cases the heuristic kicks in.
func TestDeclineFullNameFemnAll(t *testing.T) {
	in := FullName{Last: "Иванова", First: "Анна", Patronymic: "Сергеевна"}
	cases := []fullCase{
		{in, Nominative, FullName{Last: "Иванова", First: "Анна", Patronymic: "Сергеевна"}},
		{in, Genitive, FullName{Last: "Ивановой", First: "Анны", Patronymic: "Сергеевны"}},
		{in, Dative, FullName{Last: "Ивановой", First: "Анне", Patronymic: "Сергеевне"}},
		{in, Accusative, FullName{Last: "Иванову", First: "Анну", Patronymic: "Сергеевну"}},
		{in, Instrumental, FullName{Last: "Ивановой", First: "Анной", Patronymic: "Сергеевной"}},
		{in, Prepositional, FullName{Last: "Ивановой", First: "Анне", Patronymic: "Сергеевне"}},
	}
	runFullnameCases(t, cases)
}

// TestDeclineFullNameLastFirstNoPatr — the most common real-world database
// case: surname + given name, no patronymic.
func TestDeclineFullNameLastFirstNoPatr(t *testing.T) {
	cases := []fullCase{
		// Masculine: "Иванов Иван".
		{
			FullName{Last: "Иванов", First: "Иван"}, Genitive,
			FullName{Last: "Иванова", First: "Ивана"},
		},
		{
			FullName{Last: "Иванов", First: "Иван"}, Dative,
			FullName{Last: "Иванову", First: "Ивану"},
		},
		{
			FullName{Last: "Иванов", First: "Иван"}, Instrumental,
			FullName{Last: "Ивановым", First: "Иваном"},
		},
		// Feminine: "Иванова Анна".
		{
			FullName{Last: "Иванова", First: "Анна"}, Genitive,
			FullName{Last: "Ивановой", First: "Анны"},
		},
		{
			FullName{Last: "Иванова", First: "Анна"}, Accusative,
			FullName{Last: "Иванову", First: "Анну"},
		},
		{
			FullName{Last: "Иванова", First: "Анна"}, Instrumental,
			FullName{Last: "Ивановой", First: "Анной"},
		},
	}
	runFullnameCases(t, cases)
}

// TestDeclineFullNameFirstPatrNoLast — first name and patronymic only.
func TestDeclineFullNameFirstPatrNoLast(t *testing.T) {
	cases := []fullCase{
		{
			FullName{First: "Иван", Patronymic: "Иванович"}, Genitive,
			FullName{First: "Ивана", Patronymic: "Ивановича"},
		},
		{
			FullName{First: "Иван", Patronymic: "Иванович"}, Instrumental,
			FullName{First: "Иваном", Patronymic: "Ивановичем"},
		},
		{
			FullName{First: "Анна", Patronymic: "Сергеевна"}, Genitive,
			FullName{First: "Анны", Patronymic: "Сергеевны"},
		},
		{
			FullName{First: "Анна", Patronymic: "Сергеевна"}, Dative,
			FullName{First: "Анне", Patronymic: "Сергеевне"},
		},
	}
	runFullnameCases(t, cases)
}

// TestDeclineFullNameFirstOnly — first name only.
func TestDeclineFullNameFirstOnly(t *testing.T) {
	cases := []fullCase{
		{FullName{First: "Иван"}, Genitive, FullName{First: "Ивана"}},
		{FullName{First: "Иван"}, Instrumental, FullName{First: "Иваном"}},
		{FullName{First: "Анна"}, Genitive, FullName{First: "Анны"}},
		{FullName{First: "Анна"}, Accusative, FullName{First: "Анну"}},
	}
	runFullnameCases(t, cases)
}

// TestDeclineFullNameLastOnly — surname only.
// Gender is determined from an informative ending (-ов / -ова).
func TestDeclineFullNameLastOnly(t *testing.T) {
	cases := []fullCase{
		{FullName{Last: "Иванов"}, Genitive, FullName{Last: "Иванова"}},
		{FullName{Last: "Иванов"}, Instrumental, FullName{Last: "Ивановым"}},
		{FullName{Last: "Иванова"}, Genitive, FullName{Last: "Ивановой"}},
		{FullName{Last: "Иванова"}, Accusative, FullName{Last: "Иванову"}},
	}
	runFullnameCases(t, cases)
}

// TestDeclineFullNameIndeclSurnameDict — dictionary path for the
// indeclinable foreign surname Дюма (Surn+Fixd, common gender ms-f). It
// returns as is for any case in any combination with a given name or
// without one.
func TestDeclineFullNameIndeclSurnameDict(t *testing.T) {
	cases := []fullCase{
		// With the given name Александр (masc).
		{
			FullName{Last: "Дюма", First: "Александр"}, Genitive,
			FullName{Last: "Дюма", First: "Александра"},
		},
		{
			FullName{Last: "Дюма", First: "Александр"}, Instrumental,
			FullName{Last: "Дюма", First: "Александром"},
		},
		// Without the given name — gender comes from the dictionary
		// (ms-f, treated as compatible).
		{FullName{Last: "Дюма"}, Genitive, FullName{Last: "Дюма"}},
		{FullName{Last: "Дюма"}, Accusative, FullName{Last: "Дюма"}},
	}
	runFullnameCases(t, cases)
}

// TestDeclineFullNameAlphaSurnameHeuristic — a -а-ending surname outside the
// dictionary. Demonstrates the fallback to the 1st-declension heuristic.
// Окуджава is not in sample.xml, so declineLast falls into
// declineSurnameHeuristic.
//
// This covers the "-а/-я surname outside the dictionary" class: without
// stress detection the heuristic inflects such a surname by the 1st
// declension (Окуджавы/Окуджаве/...). See docs/DESIGN.md §9.3 — the
// "stressed vowels" simplification.
func TestDeclineFullNameAlphaSurnameHeuristic(t *testing.T) {
	in := FullName{Last: "Окуджава", First: "Александр"}
	cases := []fullCase{
		{in, Nominative, FullName{Last: "Окуджава", First: "Александр"}},
		{in, Genitive, FullName{Last: "Окуджавы", First: "Александра"}},
		{in, Dative, FullName{Last: "Окуджаве", First: "Александру"}},
		{in, Accusative, FullName{Last: "Окуджаву", First: "Александра"}},
		{in, Instrumental, FullName{Last: "Окуджавой", First: "Александром"}},
		{in, Prepositional, FullName{Last: "Окуджаве", First: "Александре"}},
	}
	runFullnameCases(t, cases)
}

// TestDeclineFullNameEmpty — an empty FullName → error.
func TestDeclineFullNameEmpty(t *testing.T) {
	_, err := DeclineFullName(FullName{}, Genitive)
	if err == nil {
		t.Fatal("DeclineFullName(FullName{}, ...) must return an error")
	}
	if !strings.Contains(err.Error(), "empty name") {
		t.Fatalf("expected mention of \"empty name\", got: %v", err)
	}
}

// TestDeclineFullNameUndecidableGender — no patronymic, the given name is
// not in the dictionary, and the surname is uninformative. A
// gender-determination error must be returned.
func TestDeclineFullNameUndecidableGender(t *testing.T) {
	in := FullName{First: "Дзё", Last: "Кобо"}
	_, err := DeclineFullName(in, Genitive)
	if err == nil {
		t.Fatalf("expected gender-determination error for %v", in)
	}
	if !strings.Contains(err.Error(), "could not determine gender") {
		t.Fatalf("expected mention of \"could not determine gender\", got: %v", err)
	}
}

// TestDeclineFullNameEmptyComponentsStayEmpty — empty fields in the input
// stay empty in the result (rather than turning into an error or artifact).
func TestDeclineFullNameEmptyComponentsStayEmpty(t *testing.T) {
	in := FullName{Last: "Иванов", First: "Иван"}
	out, err := DeclineFullName(in, Dative)
	if err != nil {
		t.Fatalf("DeclineFullName(%v, Dative) error: %v", in, err)
	}
	if out.Patronymic != "" {
		t.Fatalf("Patronymic must stay empty, got %q", out.Patronymic)
	}
	if out.Last != "Иванову" || out.First != "Ивану" {
		t.Fatalf("unexpected result: %+v", out)
	}
}

// TestDeclineFullNameAdjectivalSurnames — a complete full name with
// adjectival surnames (Достоевский/Достоевская, Толстой/Толстая) for
// several cases.
func TestDeclineFullNameAdjectivalSurnames(t *testing.T) {
	cases := []fullCase{
		// Достоевский — the dictionary has Nom/Gen/Dat, the rest comes from the heuristic.
		{
			FullName{Last: "Достоевский", First: "Иван"}, Genitive,
			FullName{Last: "Достоевского", First: "Ивана"},
		},
		{
			FullName{Last: "Достоевский", First: "Иван"}, Instrumental,
			FullName{Last: "Достоевским", First: "Иваном"},
		},
		// Толстой.
		{
			FullName{Last: "Толстой", First: "Иван"}, Instrumental,
			FullName{Last: "Толстым", First: "Иваном"},
		},
		// Достоевская.
		{
			FullName{Last: "Достоевская", First: "Анна"}, Genitive,
			FullName{Last: "Достоевской", First: "Анны"},
		},
		{
			FullName{Last: "Достоевская", First: "Анна"}, Accusative,
			FullName{Last: "Достоевскую", First: "Анну"},
		},
	}
	runFullnameCases(t, cases)
}

// TestDeclineFullNameDifferentNames — other given names (Дмитрий, Ольга,
// Александр) with full paradigms in the dictionary.
func TestDeclineFullNameDifferentNames(t *testing.T) {
	cases := []fullCase{
		{
			FullName{First: "Дмитрий", Patronymic: "Иванович"}, Genitive,
			FullName{First: "Дмитрия", Patronymic: "Ивановича"},
		},
		{
			FullName{First: "Дмитрий", Patronymic: "Иванович"}, Instrumental,
			FullName{First: "Дмитрием", Patronymic: "Ивановичем"},
		},
		{
			FullName{First: "Ольга", Patronymic: "Сергеевна"}, Dative,
			FullName{First: "Ольге", Patronymic: "Сергеевне"},
		},
		{
			FullName{First: "Александр", Patronymic: "Никитич"}, Prepositional,
			FullName{First: "Александре", Patronymic: "Никитиче"},
		},
	}
	runFullnameCases(t, cases)
}

// runFullnameCases — table-driven runner for FullName tests.
func runFullnameCases(t *testing.T, cases []fullCase) {
	t.Helper()
	for _, tc := range cases {
		name := joinFullName(tc.in) + "/" + caseTag(tc.c)
		t.Run(name, func(t *testing.T) {
			got, err := DeclineFullName(tc.in, tc.c)
			if err != nil {
				t.Fatalf("DeclineFullName(%+v, %v) error: %v", tc.in, tc.c, err)
			}
			if got != tc.want {
				t.Fatalf("DeclineFullName(%+v, %v) = %+v, expected %+v",
					tc.in, tc.c, got, tc.want)
			}
		})
	}
}
