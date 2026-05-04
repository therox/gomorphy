package gomorphy

import "testing"

// TestDeclineFirstHeuristicMascConsonant — out-of-dictionary masculine
// name on a consonant. Used as a regression for the user-reported
// "Элендиль Арогорн Араторович" case.
func TestDeclineFirstHeuristicMascConsonant(t *testing.T) {
	in := FullName{Last: "Элендиль", First: "Арогорн", Patronymic: "Араторович"}
	cases := []fullCase{
		{in, Nominative, FullName{Last: "Элендиль", First: "Арогорн", Patronymic: "Араторович"}},
		{in, Genitive, FullName{Last: "Элендиля", First: "Арогорна", Patronymic: "Араторовича"}},
		{in, Dative, FullName{Last: "Элендилю", First: "Арогорну", Patronymic: "Араторовичу"}},
		{in, Accusative, FullName{Last: "Элендиля", First: "Арогорна", Patronymic: "Араторовича"}},
		{in, Instrumental, FullName{Last: "Элендилем", First: "Арогорном", Patronymic: "Араторовичем"}},
		{in, Prepositional, FullName{Last: "Элендиле", First: "Арогорне", Patronymic: "Араторовиче"}},
	}
	runFullnameCases(t, cases)
}

// TestDeclineFirstHeuristicMascSoft — masculine name on -ь / -й is treated
// by the 2nd-declension heuristic when out of dictionary.
func TestDeclineFirstHeuristicMascSoft(t *testing.T) {
	cases := []fullCase{
		// Беорь — fictional name on -ь. Patronymic supplies gender.
		{
			FullName{First: "Беорь", Patronymic: "Иванович"}, Genitive,
			FullName{First: "Беоря", Patronymic: "Ивановича"},
		},
		{
			FullName{First: "Беорь", Patronymic: "Иванович"}, Instrumental,
			FullName{First: "Беорем", Patronymic: "Ивановичем"},
		},
	}
	runFullnameCases(t, cases)
}

// TestDeclineFirstHeuristicFemnA — out-of-dictionary feminine name on -а
// uses the 1st declension; -ы/-и after hush is honoured by decline1stA.
func TestDeclineFirstHeuristicFemnA(t *testing.T) {
	cases := []fullCase{
		// Эйра — fictional feminine name; patronymic supplies gender.
		{
			FullName{First: "Эйра", Patronymic: "Сергеевна"}, Genitive,
			FullName{First: "Эйры", Patronymic: "Сергеевны"},
		},
		{
			FullName{First: "Эйра", Patronymic: "Сергеевна"}, Dative,
			FullName{First: "Эйре", Patronymic: "Сергеевне"},
		},
		{
			FullName{First: "Эйра", Patronymic: "Сергеевна"}, Instrumental,
			FullName{First: "Эйрой", Patronymic: "Сергеевной"},
		},
	}
	runFullnameCases(t, cases)
}

// TestDeclineFirstHeuristicFemnYa — out-of-dictionary feminine name on
// -ия uses the special -ия subtype: dat/prep → -ии.
func TestDeclineFirstHeuristicFemnYa(t *testing.T) {
	cases := []fullCase{
		{
			FullName{First: "Нэрвия", Patronymic: "Сергеевна"}, Genitive,
			FullName{First: "Нэрвии", Patronymic: "Сергеевны"},
		},
		{
			FullName{First: "Нэрвия", Patronymic: "Сергеевна"}, Dative,
			FullName{First: "Нэрвии", Patronymic: "Сергеевне"},
		},
		{
			FullName{First: "Нэрвия", Patronymic: "Сергеевна"}, Instrumental,
			FullName{First: "Нэрвией", Patronymic: "Сергеевной"},
		},
	}
	runFullnameCases(t, cases)
}

// TestDeclineFirstHeuristicIndecl — first names ending in a non-а/-я
// vowel are foreign and indeclinable.
func TestDeclineFirstHeuristicIndecl(t *testing.T) {
	cases := []fullCase{
		{
			FullName{First: "Бруно", Patronymic: "Иванович"}, Genitive,
			FullName{First: "Бруно", Patronymic: "Ивановича"},
		},
		{
			FullName{First: "Жозе", Patronymic: "Иванович"}, Instrumental,
			FullName{First: "Жозе", Patronymic: "Ивановичем"},
		},
	}
	runFullnameCases(t, cases)
}

// TestToNominativeFirstNameHeuristic — inverse direction: the name
// "Арогорн" round-trips through DeclineFullName + ToNominative for a
// complete masculine triple.
func TestToNominativeFirstNameHeuristic(t *testing.T) {
	nom := FullName{Last: "Элендиль", First: "Арогорн", Patronymic: "Араторович"}
	for _, c := range allCases() {
		t.Run(caseTag(c), func(t *testing.T) {
			declined, err := DeclineFullName(nom, c)
			if err != nil {
				t.Fatalf("DeclineFullName(%v, %v) error: %v", nom, c, err)
			}
			got, err := ToNominative(declined)
			if err != nil {
				t.Fatalf("ToNominative(%v) error: %v", declined, err)
			}
			if got.First != nom.First {
				t.Fatalf("First round-trip: %q → %q → %q (case %v)",
					nom.First, declined.First, got.First, c)
			}
		})
	}
}

// TestInverseFirstNameHeuristicFemn — inverse for a femn -а stem name.
func TestInverseFirstNameHeuristicFemn(t *testing.T) {
	// Round-trip via DeclineFullName + ToNominative.
	nom := FullName{First: "Эйра", Patronymic: "Сергеевна"}
	for _, c := range allCases() {
		t.Run(caseTag(c), func(t *testing.T) {
			declined, err := DeclineFullName(nom, c)
			if err != nil {
				t.Fatalf("DeclineFullName(%v, %v) error: %v", nom, c, err)
			}
			got, err := ToNominative(declined)
			if err != nil {
				t.Fatalf("ToNominative(%v) error: %v", declined, err)
			}
			if got.First != nom.First {
				t.Fatalf("First round-trip: %q → %q → %q (case %v)",
					nom.First, declined.First, got.First, c)
			}
		})
	}
}
