package gomorphy

import "testing"

// inversePatrCase — one expectation row: input (any case), expected Nom,
// expected gender.
type inversePatrCase struct {
	in    string
	wantS string
	wantG Gender
}

// TestInversePatronymicMasc — all 4 masculine patronymic types × 6 cases.
func TestInversePatronymicMasc(t *testing.T) {
	cases := []inversePatrCase{
		// -ович (Иванович).
		{"иванович", "иванович", Masculine},
		{"ивановича", "иванович", Masculine},
		{"ивановичу", "иванович", Masculine},
		{"ивановичем", "иванович", Masculine},
		{"ивановиче", "иванович", Masculine},
		// -евич (Сергеевич).
		{"сергеевич", "сергеевич", Masculine},
		{"сергеевича", "сергеевич", Masculine},
		{"сергеевичу", "сергеевич", Masculine},
		{"сергеевичем", "сергеевич", Masculine},
		{"сергеевиче", "сергеевич", Masculine},
		// -ич (Никитич).
		{"никитич", "никитич", Masculine},
		{"никитича", "никитич", Masculine},
		{"никитичу", "никитич", Masculine},
		{"никитичем", "никитич", Masculine},
		{"никитиче", "никитич", Masculine},
		// -ьич (Ильич).
		{"ильич", "ильич", Masculine},
		{"ильича", "ильич", Masculine},
		{"ильичу", "ильич", Masculine},
		{"ильичем", "ильич", Masculine},
		{"ильиче", "ильич", Masculine},
	}
	runInversePatrCases(t, cases)
}

// TestInversePatronymicFemn — all 4 feminine types × 6 cases.
func TestInversePatronymicFemn(t *testing.T) {
	cases := []inversePatrCase{
		// -овна (Ивановна).
		{"ивановна", "ивановна", Feminine},
		{"ивановны", "ивановна", Feminine},
		{"ивановне", "ивановна", Feminine},
		{"ивановну", "ивановна", Feminine},
		{"ивановной", "ивановна", Feminine},
		// -евна (Сергеевна).
		{"сергеевна", "сергеевна", Feminine},
		{"сергеевны", "сергеевна", Feminine},
		{"сергеевне", "сергеевна", Feminine},
		{"сергеевну", "сергеевна", Feminine},
		{"сергеевной", "сергеевна", Feminine},
		// -ична (Никитична).
		{"никитична", "никитична", Feminine},
		{"никитичны", "никитична", Feminine},
		{"никитичне", "никитична", Feminine},
		{"никитичну", "никитична", Feminine},
		{"никитичной", "никитична", Feminine},
		// -инична (Ильинична). The -иничн stem must be recognized before
		// -ичн, otherwise "Ильиничны" would be truncated to "Ильинича"
		// instead of "Ильинична".
		{"ильинична", "ильинична", Feminine},
		{"ильиничны", "ильинична", Feminine},
		{"ильиничне", "ильинична", Feminine},
		{"ильиничну", "ильинична", Feminine},
		{"ильиничной", "ильинична", Feminine},
	}
	runInversePatrCases(t, cases)
}

// TestInversePatronymicReject — tokens that do not look like a patronymic.
func TestInversePatronymicReject(t *testing.T) {
	for _, in := range []string{"иванов", "анна", "скрипка", ""} {
		t.Run(in, func(t *testing.T) {
			if _, _, ok := inversePatronymic(in); ok {
				t.Fatalf("inversePatronymic(%q) — must return ok=false", in)
			}
		})
	}
}

// inverseSurnCase — input, gHint, expected Nom, expected gender.
type inverseSurnCase struct {
	in    string
	hint  Gender
	wantS string
	wantG Gender
}

// TestInverseSurnamePossMasc — possessive masculine (-ов/-ев/-ёв/-ин/-ын).
func TestInverseSurnamePossMasc(t *testing.T) {
	cases := []inverseSurnCase{
		// -ов: Nom + 5 case forms. With gHint=Masculine the ambiguous
		// -а/-у are resolved.
		{"иванов", Masculine, "иванов", Masculine},
		{"иванова", Masculine, "иванов", Masculine},
		{"иванову", Masculine, "иванов", Masculine},
		{"ивановым", Masculine, "иванов", Masculine},
		{"иванове", Masculine, "иванов", Masculine},
		// -ев (Сергеев).
		{"сергеев", Masculine, "сергеев", Masculine},
		{"сергеева", Masculine, "сергеев", Masculine},
		{"сергееву", Masculine, "сергеев", Masculine},
		{"сергеевым", Masculine, "сергеев", Masculine},
		{"сергееве", Masculine, "сергеев", Masculine},
		// -ёв (Соловьёв).
		{"соловьёв", Masculine, "соловьёв", Masculine},
		{"соловьёва", Masculine, "соловьёв", Masculine},
		{"соловьёву", Masculine, "соловьёв", Masculine},
		{"соловьёвым", Masculine, "соловьёв", Masculine},
		{"соловьёве", Masculine, "соловьёв", Masculine},
		// -ин (Пушкин).
		{"пушкин", Masculine, "пушкин", Masculine},
		{"пушкина", Masculine, "пушкин", Masculine},
		{"пушкину", Masculine, "пушкин", Masculine},
		{"пушкиным", Masculine, "пушкин", Masculine},
		{"пушкине", Masculine, "пушкин", Masculine},
		// -ын (Куницын).
		{"куницын", Masculine, "куницын", Masculine},
		{"куницына", Masculine, "куницын", Masculine},
		{"куницыну", Masculine, "куницын", Masculine},
		{"куницыным", Masculine, "куницын", Masculine},
		{"куницыне", Masculine, "куницын", Masculine},
	}
	runInverseSurnCases(t, cases)
}

// TestInverseSurnamePossFemn — possessive feminine (-ова/-ева/-ёва/-ина/-ына).
// gHint=Feminine resolves the ambiguous -а/-у in favor of F.
func TestInverseSurnamePossFemn(t *testing.T) {
	cases := []inverseSurnCase{
		// -ова: Nom + 5 case forms.
		{"иванова", Feminine, "иванова", Feminine},
		{"ивановой", Feminine, "иванова", Feminine},
		{"иванову", Feminine, "иванова", Feminine},
		// -ева (Сергеева).
		{"сергеева", Feminine, "сергеева", Feminine},
		{"сергеевой", Feminine, "сергеева", Feminine},
		{"сергееву", Feminine, "сергеева", Feminine},
		// -ёва (Соловьёва).
		{"соловьёва", Feminine, "соловьёва", Feminine},
		{"соловьёвой", Feminine, "соловьёва", Feminine},
		{"соловьёву", Feminine, "соловьёва", Feminine},
		// -ина (Пушкина).
		{"пушкина", Feminine, "пушкина", Feminine},
		{"пушкиной", Feminine, "пушкина", Feminine},
		{"пушкину", Feminine, "пушкина", Feminine},
		// -ына (Куницына).
		{"куницына", Feminine, "куницына", Feminine},
		{"куницыной", Feminine, "куницына", Feminine},
		{"куницыну", Feminine, "куницына", Feminine},
	}
	runInverseSurnCases(t, cases)
}

// TestInverseSurnameAdjMasc — adjectival masculine (-ский/-цкий/-ой/-ий).
func TestInverseSurnameAdjMasc(t *testing.T) {
	cases := []inverseSurnCase{
		// -ский (Достоевский) — five case forms each.
		{"достоевский", Masculine, "достоевский", Masculine},
		{"достоевского", Masculine, "достоевский", Masculine},
		{"достоевскому", Masculine, "достоевский", Masculine},
		{"достоевским", Masculine, "достоевский", Masculine},
		{"достоевском", Masculine, "достоевский", Masculine},
		// -цкий (Высоцкий).
		{"высоцкий", Masculine, "высоцкий", Masculine},
		{"высоцкого", Masculine, "высоцкий", Masculine},
		{"высоцкому", Masculine, "высоцкий", Masculine},
		{"высоцким", Masculine, "высоцкий", Masculine},
		{"высоцком", Masculine, "высоцкий", Masculine},
		// -ой (Толстой) — stem "толст", after "т" Nom = -ой.
		{"толстой", Masculine, "толстой", Masculine},
		{"толстого", Masculine, "толстой", Masculine},
		{"толстому", Masculine, "толстой", Masculine},
		{"толстым", Masculine, "толстой", Masculine},
		{"толстом", Masculine, "толстой", Masculine},
		// -ий after "к" (Горький).
		{"горький", Masculine, "горький", Masculine},
		{"горького", Masculine, "горький", Masculine},
		{"горьким", Masculine, "горький", Masculine},
	}
	runInverseSurnCases(t, cases)
}

// TestInverseSurnameAdjFemn — adjectival feminine (-ская/-цкая/-ая/-яя).
func TestInverseSurnameAdjFemn(t *testing.T) {
	cases := []inverseSurnCase{
		// -ская (Достоевская).
		{"достоевская", Feminine, "достоевская", Feminine},
		{"достоевской", Feminine, "достоевская", Feminine},
		{"достоевскую", Feminine, "достоевская", Feminine},
		// -цкая (Высоцкая).
		{"высоцкая", Feminine, "высоцкая", Feminine},
		{"высоцкой", Feminine, "высоцкая", Feminine},
		{"высоцкую", Feminine, "высоцкая", Feminine},
		// -ая (Толстая).
		{"толстая", Feminine, "толстая", Feminine},
		{"толстой", Feminine, "толстая", Feminine}, // gHint=F resolves Толстой → Толстая
		{"толстую", Feminine, "толстая", Feminine},
		// -яя (Зимняя).
		{"зимняя", Feminine, "зимняя", Feminine},
		{"зимнюю", Feminine, "зимняя", Feminine},
	}
	runInverseSurnCases(t, cases)
}

// TestInverseSurnameAmbiguousNoHint — ambiguous forms without a hint.
// By default -ова/-ову is interpreted as F Nom.
func TestInverseSurnameAmbiguousNoHint(t *testing.T) {
	cases := []inverseSurnCase{
		{"иванова", GenderUnknown, "иванова", Feminine},
		{"иванову", GenderUnknown, "иванова", Feminine},
		// "Толстой" without a hint: M Nom (as is).
		{"толстой", GenderUnknown, "толстой", Masculine},
	}
	runInverseSurnCases(t, cases)
}

// TestInverseSurnameIndeclYkh — indeclinable -ых/-их stay as is.
func TestInverseSurnameIndeclYkh(t *testing.T) {
	for _, in := range []string{"долгих", "седых"} {
		t.Run(in, func(t *testing.T) {
			got, _, ok := inverseSurnameHeuristic(in, GenderUnknown)
			if !ok || got != in {
				t.Fatalf("inverseSurnameHeuristic(%q) = (%q, _, %v), expected (%q, _, true)",
					in, got, ok, in)
			}
		})
	}
}

func runInversePatrCases(t *testing.T, cases []inversePatrCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.in, func(t *testing.T) {
			got, g, ok := inversePatronymic(tc.in)
			if !ok {
				t.Fatalf("inversePatronymic(%q) — not recognized", tc.in)
			}
			if got != tc.wantS {
				t.Fatalf("inversePatronymic(%q) = %q, expected %q", tc.in, got, tc.wantS)
			}
			if g != tc.wantG {
				t.Fatalf("inversePatronymic(%q) gender = %v, expected %v", tc.in, g, tc.wantG)
			}
		})
	}
}

func runInverseSurnCases(t *testing.T, cases []inverseSurnCase) {
	t.Helper()
	for _, tc := range cases {
		name := tc.in + "/" + genderTag(tc.hint)
		t.Run(name, func(t *testing.T) {
			got, g, ok := inverseSurnameHeuristic(tc.in, tc.hint)
			if !ok {
				t.Fatalf("inverseSurnameHeuristic(%q, %v) — not recognized", tc.in, tc.hint)
			}
			if got != tc.wantS {
				t.Fatalf("inverseSurnameHeuristic(%q, %v) = %q, expected %q",
					tc.in, tc.hint, got, tc.wantS)
			}
			if g != tc.wantG {
				t.Fatalf("inverseSurnameHeuristic(%q, %v) gender = %v, expected %v",
					tc.in, tc.hint, g, tc.wantG)
			}
		})
	}
}
