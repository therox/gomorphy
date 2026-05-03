package gomorphy

import "testing"

// surnCase — one expectation row: input surname (lowercase), case, gender,
// expected form.
type surnCase struct {
	last string
	c    Case
	g    Gender
	want string
}

// TestSurnameHeuristicPossMasc — possessive masculine surnames:
// -ов, -ев, -ёв, -ин, -ын. One example per suffix in six cases (per the
// plan: "one example per pattern from the table above, in 6 cases").
func TestSurnameHeuristicPossMasc(t *testing.T) {
	cases := []surnCase{
		// -ов (Иванов).
		{"иванов", Nominative, Masculine, "иванов"},
		{"иванов", Genitive, Masculine, "иванова"},
		{"иванов", Dative, Masculine, "иванову"},
		{"иванов", Accusative, Masculine, "иванова"},
		{"иванов", Instrumental, Masculine, "ивановым"},
		{"иванов", Prepositional, Masculine, "иванове"},

		// -ев (Сергеев).
		{"сергеев", Nominative, Masculine, "сергеев"},
		{"сергеев", Genitive, Masculine, "сергеева"},
		{"сергеев", Dative, Masculine, "сергееву"},
		{"сергеев", Accusative, Masculine, "сергеева"},
		{"сергеев", Instrumental, Masculine, "сергеевым"},
		{"сергеев", Prepositional, Masculine, "сергееве"},

		// -ёв (Соловьёв).
		{"соловьёв", Nominative, Masculine, "соловьёв"},
		{"соловьёв", Genitive, Masculine, "соловьёва"},
		{"соловьёв", Dative, Masculine, "соловьёву"},
		{"соловьёв", Accusative, Masculine, "соловьёва"},
		{"соловьёв", Instrumental, Masculine, "соловьёвым"},
		{"соловьёв", Prepositional, Masculine, "соловьёве"},

		// -ин (Пушкин).
		{"пушкин", Nominative, Masculine, "пушкин"},
		{"пушкин", Genitive, Masculine, "пушкина"},
		{"пушкин", Dative, Masculine, "пушкину"},
		{"пушкин", Accusative, Masculine, "пушкина"},
		{"пушкин", Instrumental, Masculine, "пушкиным"},
		{"пушкин", Prepositional, Masculine, "пушкине"},

		// -ын (Куницын).
		{"куницын", Nominative, Masculine, "куницын"},
		{"куницын", Genitive, Masculine, "куницына"},
		{"куницын", Dative, Masculine, "куницыну"},
		{"куницын", Accusative, Masculine, "куницына"},
		{"куницын", Instrumental, Masculine, "куницыным"},
		{"куницын", Prepositional, Masculine, "куницыне"},
	}
	runSurnameCases(t, cases)
}

// TestSurnameHeuristicPossFemn — possessive feminine surnames:
// -ова, -ева, -ёва, -ина, -ына. One example × 6 cases.
func TestSurnameHeuristicPossFemn(t *testing.T) {
	cases := []surnCase{
		// -ова (Иванова).
		{"иванова", Nominative, Feminine, "иванова"},
		{"иванова", Genitive, Feminine, "ивановой"},
		{"иванова", Dative, Feminine, "ивановой"},
		{"иванова", Accusative, Feminine, "иванову"},
		{"иванова", Instrumental, Feminine, "ивановой"},
		{"иванова", Prepositional, Feminine, "ивановой"},

		// -ева (Сергеева).
		{"сергеева", Nominative, Feminine, "сергеева"},
		{"сергеева", Genitive, Feminine, "сергеевой"},
		{"сергеева", Dative, Feminine, "сергеевой"},
		{"сергеева", Accusative, Feminine, "сергееву"},
		{"сергеева", Instrumental, Feminine, "сергеевой"},
		{"сергеева", Prepositional, Feminine, "сергеевой"},

		// -ёва (Соловьёва).
		{"соловьёва", Nominative, Feminine, "соловьёва"},
		{"соловьёва", Genitive, Feminine, "соловьёвой"},
		{"соловьёва", Dative, Feminine, "соловьёвой"},
		{"соловьёва", Accusative, Feminine, "соловьёву"},
		{"соловьёва", Instrumental, Feminine, "соловьёвой"},
		{"соловьёва", Prepositional, Feminine, "соловьёвой"},

		// -ина (Пушкина).
		{"пушкина", Nominative, Feminine, "пушкина"},
		{"пушкина", Genitive, Feminine, "пушкиной"},
		{"пушкина", Dative, Feminine, "пушкиной"},
		{"пушкина", Accusative, Feminine, "пушкину"},
		{"пушкина", Instrumental, Feminine, "пушкиной"},
		{"пушкина", Prepositional, Feminine, "пушкиной"},

		// -ына (Куницына).
		{"куницына", Nominative, Feminine, "куницына"},
		{"куницына", Genitive, Feminine, "куницыной"},
		{"куницына", Dative, Feminine, "куницыной"},
		{"куницына", Accusative, Feminine, "куницыну"},
		{"куницына", Instrumental, Feminine, "куницыной"},
		{"куницына", Prepositional, Feminine, "куницыной"},
	}
	runSurnameCases(t, cases)
}

// TestSurnameHeuristicAdjMasc — adjectival masculine: -ский, -цкий, -ой,
// -ый, -ий. One example × 6 cases.
func TestSurnameHeuristicAdjMasc(t *testing.T) {
	cases := []surnCase{
		// -ский (Достоевский) — after "к" the instrumental ending is "им".
		{"достоевский", Nominative, Masculine, "достоевский"},
		{"достоевский", Genitive, Masculine, "достоевского"},
		{"достоевский", Dative, Masculine, "достоевскому"},
		{"достоевский", Accusative, Masculine, "достоевского"},
		{"достоевский", Instrumental, Masculine, "достоевским"},
		{"достоевский", Prepositional, Masculine, "достоевском"},

		// -цкий (Высоцкий).
		{"высоцкий", Nominative, Masculine, "высоцкий"},
		{"высоцкий", Genitive, Masculine, "высоцкого"},
		{"высоцкий", Dative, Masculine, "высоцкому"},
		{"высоцкий", Accusative, Masculine, "высоцкого"},
		{"высоцкий", Instrumental, Masculine, "высоцким"},
		{"высоцкий", Prepositional, Masculine, "высоцком"},

		// -ой (Толстой) — after "т" the instrumental ending is "ым".
		{"толстой", Nominative, Masculine, "толстой"},
		{"толстой", Genitive, Masculine, "толстого"},
		{"толстой", Dative, Masculine, "толстому"},
		{"толстой", Accusative, Masculine, "толстого"},
		{"толстой", Instrumental, Masculine, "толстым"},
		{"толстой", Prepositional, Masculine, "толстом"},

		// -ый (Белый) — after "л" the instrumental ending is "ым".
		{"белый", Nominative, Masculine, "белый"},
		{"белый", Genitive, Masculine, "белого"},
		{"белый", Dative, Masculine, "белому"},
		{"белый", Accusative, Masculine, "белого"},
		{"белый", Instrumental, Masculine, "белым"},
		{"белый", Prepositional, Masculine, "белом"},

		// -ий (Горький) — after "к" the instrumental ending is "им".
		{"горький", Nominative, Masculine, "горький"},
		{"горький", Genitive, Masculine, "горького"},
		{"горький", Dative, Masculine, "горькому"},
		{"горький", Accusative, Masculine, "горького"},
		{"горький", Instrumental, Masculine, "горьким"},
		{"горький", Prepositional, Masculine, "горьком"},
	}
	runSurnameCases(t, cases)
}

// TestSurnameHeuristicAdjFemn — adjectival feminine: -ская, -цкая, -ая, -яя.
// One example × 6 cases.
func TestSurnameHeuristicAdjFemn(t *testing.T) {
	cases := []surnCase{
		// -ская (Достоевская).
		{"достоевская", Nominative, Feminine, "достоевская"},
		{"достоевская", Genitive, Feminine, "достоевской"},
		{"достоевская", Dative, Feminine, "достоевской"},
		{"достоевская", Accusative, Feminine, "достоевскую"},
		{"достоевская", Instrumental, Feminine, "достоевской"},
		{"достоевская", Prepositional, Feminine, "достоевской"},

		// -цкая (Высоцкая).
		{"высоцкая", Nominative, Feminine, "высоцкая"},
		{"высоцкая", Genitive, Feminine, "высоцкой"},
		{"высоцкая", Dative, Feminine, "высоцкой"},
		{"высоцкая", Accusative, Feminine, "высоцкую"},
		{"высоцкая", Instrumental, Feminine, "высоцкой"},
		{"высоцкая", Prepositional, Feminine, "высоцкой"},

		// -ая (Толстая).
		{"толстая", Nominative, Feminine, "толстая"},
		{"толстая", Genitive, Feminine, "толстой"},
		{"толстая", Dative, Feminine, "толстой"},
		{"толстая", Accusative, Feminine, "толстую"},
		{"толстая", Instrumental, Feminine, "толстой"},
		{"толстая", Prepositional, Feminine, "толстой"},

		// -яя (Зимняя — soft type, endings -ей/-юю).
		{"зимняя", Nominative, Feminine, "зимняя"},
		{"зимняя", Genitive, Feminine, "зимней"},
		{"зимняя", Dative, Feminine, "зимней"},
		{"зимняя", Accusative, Feminine, "зимнюю"},
		{"зимняя", Instrumental, Feminine, "зимней"},
		{"зимняя", Prepositional, Feminine, "зимней"},
	}
	runSurnameCases(t, cases)
}

// TestSurnameHeuristicIndeclYkh — indeclinable -ых/-их (Долгих, Седых).
// Common gender: both masculine and feminine inputs return the original
// form in every case.
func TestSurnameHeuristicIndeclYkh(t *testing.T) {
	for _, last := range []string{"долгих", "седых"} {
		for _, c := range allCases() {
			t.Run(last+"/"+caseTag(c), func(t *testing.T) {
				got, ok := declineSurnameHeuristic(last, c, Masculine)
				if !ok || got != last {
					t.Fatalf("declineSurnameHeuristic(%q, %v, M) = (%q, %v), expected (%q, true)",
						last, c, got, ok, last)
				}
				got, ok = declineSurnameHeuristic(last, c, Feminine)
				if !ok || got != last {
					t.Fatalf("declineSurnameHeuristic(%q, %v, F) = (%q, %v), expected (%q, true)",
						last, c, got, ok, last)
				}
			})
		}
	}
}

// TestSurnameHeuristicConsMasc — masculine surname ending in a consonant,
// -й or -ь. One example × 6 cases.
func TestSurnameHeuristicConsMasc(t *testing.T) {
	cases := []surnCase{
		// -й (Гайдай) — soft type.
		{"гайдай", Nominative, Masculine, "гайдай"},
		{"гайдай", Genitive, Masculine, "гайдая"},
		{"гайдай", Dative, Masculine, "гайдаю"},
		{"гайдай", Accusative, Masculine, "гайдая"},
		{"гайдай", Instrumental, Masculine, "гайдаем"},
		{"гайдай", Prepositional, Masculine, "гайдае"},

		// -ь (Соболь).
		{"соболь", Nominative, Masculine, "соболь"},
		{"соболь", Genitive, Masculine, "соболя"},
		{"соболь", Dative, Masculine, "соболю"},
		{"соболь", Accusative, Masculine, "соболя"},
		{"соболь", Instrumental, Masculine, "соболем"},
		{"соболь", Prepositional, Masculine, "соболе"},

		// Consonant (Шевчук).
		{"шевчук", Nominative, Masculine, "шевчук"},
		{"шевчук", Genitive, Masculine, "шевчука"},
		{"шевчук", Dative, Masculine, "шевчуку"},
		{"шевчук", Accusative, Masculine, "шевчука"},
		{"шевчук", Instrumental, Masculine, "шевчуком"},
		{"шевчук", Prepositional, Masculine, "шевчуке"},
	}
	runSurnameCases(t, cases)
}

// TestSurnameHeuristic1stDecl — surnames in -а/-я using the 1st declension.
// Includes "Дюма" — a known simplification of the heuristic: stress is not
// detected, so Дюма is inflected like Скрипка (Дюмы/Дюме/Дюму...). For
// correct behavior "Дюма" should be in the dictionary as Fixd (see
// docs/DESIGN.md §9.3).
func TestSurnameHeuristic1stDecl(t *testing.T) {
	cases := []surnCase{
		// Скрипка — after "к" the genitive ending is "и".
		{"скрипка", Nominative, Masculine, "скрипка"},
		{"скрипка", Genitive, Masculine, "скрипки"},
		{"скрипка", Dative, Masculine, "скрипке"},
		{"скрипка", Accusative, Masculine, "скрипку"},
		{"скрипка", Instrumental, Masculine, "скрипкой"},
		{"скрипка", Prepositional, Masculine, "скрипке"},

		// Заря (-я).
		{"заря", Nominative, Feminine, "заря"},
		{"заря", Genitive, Feminine, "зари"},
		{"заря", Dative, Feminine, "заре"},
		{"заря", Accusative, Feminine, "зарю"},
		{"заря", Instrumental, Feminine, "зарей"},
		{"заря", Prepositional, Feminine, "заре"},

		// Дюма — heuristic simplification: after "м" the genitive is "ы".
		// In reality, a French surname with stressed -а is indeclinable;
		// without stress detection the heuristic still inflects it.
		{"дюма", Nominative, Masculine, "дюма"},
		{"дюма", Genitive, Masculine, "дюмы"},
		{"дюма", Dative, Masculine, "дюме"},
		{"дюма", Accusative, Masculine, "дюму"},
		{"дюма", Instrumental, Masculine, "дюмой"},
		{"дюма", Prepositional, Masculine, "дюме"},

		// -ия subtype (Бария — a hypothetical surname).
		{"бария", Nominative, Feminine, "бария"},
		{"бария", Genitive, Feminine, "барии"},
		{"бария", Dative, Feminine, "барии"},
		{"бария", Accusative, Feminine, "барию"},
		{"бария", Instrumental, Feminine, "барией"},
		{"бария", Prepositional, Feminine, "барии"},
	}
	runSurnameCases(t, cases)
}

// TestSurnameHeuristicIndeclVowel — indeclinable surnames ending in a
// final vowel other than -а/-я (which fall under the 1st declension).
// Covers -о, -е, -и. The -у, -ю, -ы, -э endings are checked through the
// same path — the switch over the last rune in isVowel is uniform.
func TestSurnameHeuristicIndeclVowel(t *testing.T) {
	for _, last := range []string{"шевченко", "кобо", "дали", "хосе"} {
		for _, c := range allCases() {
			t.Run(last+"/"+caseTag(c), func(t *testing.T) {
				got, ok := declineSurnameHeuristic(last, c, Masculine)
				if !ok || got != last {
					t.Fatalf("declineSurnameHeuristic(%q, %v, M) = (%q, %v), expected (%q, true)",
						last, c, got, ok, last)
				}
			})
		}
	}
}

// TestSurnameGenderDetection — gender from an informative ending.
func TestSurnameGenderDetection(t *testing.T) {
	cases := []struct {
		last string
		want Gender
	}{
		{"Иванов", Masculine},
		{"Петров", Masculine},
		{"Соловьёв", Masculine},
		{"Пушкин", Masculine},
		{"Куницын", Masculine},
		{"Достоевский", Masculine},
		{"Высоцкий", Masculine},
		{"Иванова", Feminine},
		{"Соловьёва", Feminine},
		{"Пушкина", Feminine},
		{"Куницына", Feminine},
		{"Достоевская", Feminine},
		{"Высоцкая", Feminine},
		// Uninformative endings.
		{"Гайдай", GenderUnknown},
		{"Скрипка", GenderUnknown},
		{"Долгих", GenderUnknown},
		{"Шевченко", GenderUnknown},
		{"Дюма", GenderUnknown},
	}
	for _, tc := range cases {
		t.Run(tc.last, func(t *testing.T) {
			got := genderFromSurnameHeuristic(tc.last)
			if got != tc.want {
				t.Fatalf("genderFromSurnameHeuristic(%q) = %v, expected %v",
					tc.last, got, tc.want)
			}
		})
	}
}

// runSurnameCases — shared runner for surnCase tables.
func runSurnameCases(t *testing.T, cases []surnCase) {
	t.Helper()
	for _, tc := range cases {
		name := tc.last + "/" + caseTag(tc.c) + "/" + genderTag(tc.g)
		t.Run(name, func(t *testing.T) {
			got, ok := declineSurnameHeuristic(tc.last, tc.c, tc.g)
			if !ok {
				t.Fatalf("declineSurnameHeuristic(%q, %v, %v) — pattern did not match",
					tc.last, tc.c, tc.g)
			}
			if got != tc.want {
				t.Fatalf("declineSurnameHeuristic(%q, %v, %v) = %q, expected %q",
					tc.last, tc.c, tc.g, got, tc.want)
			}
		})
	}
}

// allCases — the six basic cases for iteration in tests.
func allCases() []Case {
	return []Case{Nominative, Genitive, Dative, Accusative, Instrumental, Prepositional}
}
