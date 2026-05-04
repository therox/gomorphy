package gomorphy_test

import (
	"fmt"
	"log"

	"github.com/therox/gomorphy"
)

// The dictionary is loaded once for the whole test binary by TestMain in
// init_test.go via the unexported useDict helper, so the examples below
// can call the public API without an explicit Init.
//
// In production code the caller must initialise the dictionary first
// (gomorphy.Init or the GOMORPHY_DICT environment variable) — see
// ExampleInit.

// ExampleInit shows the typical bootstrap: load a prebuilt .bin file and
// call any API afterwards. The example is compiled but not executed,
// because the global singleton is already populated by the test runner.
func ExampleInit() {
	if err := gomorphy.Init("dict.bin"); err != nil {
		log.Fatal(err)
	}
	word, _ := gomorphy.Decline("аппетит", gomorphy.Genitive, gomorphy.Singular)
	fmt.Println(word)
}

// ExampleDecline shows declension of a noun in singular and plural.
func ExampleDecline() {
	gen, _ := gomorphy.Decline("аппетит", gomorphy.Genitive, gomorphy.Singular)
	fmt.Println(gen)

	genPlur, _ := gomorphy.Decline("аппетит", gomorphy.Genitive, gomorphy.Plural)
	fmt.Println(genPlur)
	// Output:
	// аппетита
	// аппетитов
}

// ExampleDeclineAdj shows an adjective in different gender / case / animacy
// combinations. Animacy only matters for accusative masculine singular and
// accusative plural — it is ignored otherwise.
func ExampleDeclineAdj() {
	// Feminine nominative singular.
	femnNom, _ := gomorphy.DeclineAdj("красный",
		gomorphy.Nominative, gomorphy.Singular, gomorphy.Feminine, false)
	fmt.Println(femnNom)

	// Masculine accusative singular: animate vs inanimate.
	mascAccAnim, _ := gomorphy.DeclineAdj("красный",
		gomorphy.Accusative, gomorphy.Singular, gomorphy.Masculine, true)
	mascAccInan, _ := gomorphy.DeclineAdj("красный",
		gomorphy.Accusative, gomorphy.Singular, gomorphy.Masculine, false)
	fmt.Println(mascAccAnim, mascAccInan)
	// Output:
	// красная
	// красного красный
}

// ExampleAgree shows numeral agreement: a noun is put into the form
// required after a numeral (1 / 2–4 / 5+ / 11–14 special case).
func ExampleAgree() {
	for _, n := range []int{1, 2, 5, 12, 21} {
		form, _ := gomorphy.Agree("яблоко", n)
		fmt.Printf("%d %s\n", n, form)
	}
	// Output:
	// 1 яблоко
	// 2 яблока
	// 5 яблок
	// 12 яблок
	// 21 яблоко
}

// ExampleGenderOf shows lemma-level gender lookup for nouns of every
// gender, including common (сирота).
func ExampleGenderOf() {
	for _, w := range []string{"стол", "книга", "окно", "сирота"} {
		g, _ := gomorphy.GenderOf(w)
		fmt.Printf("%s — %s\n", w, g)
	}
	// Output:
	// стол — м.р.
	// книга — ж.р.
	// окно — ср.р.
	// сирота — общ.р.
}

// ExamplePluralOf shows nominative plural lookup, including suppletive
// pairs (человек → люди, ребёнок → дети) and pluralia tantum (the lemma
// itself is already plural).
func ExamplePluralOf() {
	for _, w := range []string{"стол", "ребёнок", "ножницы"} {
		p, _ := gomorphy.PluralOf(w)
		fmt.Printf("%s → %s\n", w, p)
	}
	// Output:
	// стол → столы
	// ребёнок → дети
	// ножницы → ножницы
}

// ExampleParse shows morphological analysis of an unambiguous word form:
// the lemma plus the grammatical features of the analyzed form.
func ExampleParse() {
	got, _ := gomorphy.Parse("аппетита")
	for _, a := range got {
		fmt.Printf("lemma=%s pos=%s case=%s number=%s gender=%s animate=%t\n",
			a.Lemma, a.POS, a.Case, a.Number, a.Gender, a.Animate)
	}
	// Output:
	// lemma=аппетит pos=NOUN case=род.п. number=ед.ч. gender=м.р. animate=false
}

// ExamplePluralOf_suppletive isolates the suppletive pair человек → люди.
func ExamplePluralOf_suppletive() {
	p, _ := gomorphy.PluralOf("человек")
	fmt.Println(p)
	// Output: люди
}

// ExampleDeclineFullName shows declension of a complete masculine full
// name. Each component is inflected by its own rules; the case of the
// first letter is preserved.
func ExampleDeclineFullName() {
	nom := gomorphy.FullName{
		Last:       "Иванов",
		First:      "Иван",
		Patronymic: "Иванович",
	}
	gen, err := gomorphy.DeclineFullName(nom, gomorphy.Genitive)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gen.Last, gen.First, gen.Patronymic)
	// Output: Иванова Ивана Ивановича
}

// ExampleToNominative reduces an inflected feminine full name back to the
// nominative case — useful as the inverse of DeclineFullName.
func ExampleToNominative() {
	in := gomorphy.FullName{
		Last:       "Ивановой",
		First:      "Анне",
		Patronymic: "Сергеевне",
	}
	nom, err := gomorphy.ToNominative(in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(nom.Last, nom.First, nom.Patronymic)
	// Output: Иванова Анна Сергеевна
}

// ExampleParseFullName parses a free-form string in any case into a
// FullName whose components are reduced to the nominative.
func ExampleParseFullName() {
	for _, s := range []string{
		"Иванов Иван Иванович",    // Russian order, Nom
		"Иван Иванович Иванов",    // Western order
		"Ивановой Анне Сергеевне", // Russian order, Dat
		"Иван Иванов",             // two tokens
	} {
		fn, err := gomorphy.ParseFullName(s)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%q → %s | %s | %s\n", s, fn.Last, fn.First, fn.Patronymic)
	}
	// Output:
	// "Иванов Иван Иванович" → Иванов | Иван | Иванович
	// "Иван Иванович Иванов" → Иванов | Иван | Иванович
	// "Ивановой Анне Сергеевне" → Иванова | Анна | Сергеевна
	// "Иван Иванов" → Иванов | Иван |
}

// ExampleShortName demonstrates the basic case: a full FullName struct is
// reduced to "Surname I. P." form.
func ExampleShortName() {
	fn := gomorphy.FullName{
		Last:       "Иванов",
		First:      "Иван",
		Patronymic: "Иванович",
	}
	fmt.Println(gomorphy.ShortName(fn))
	// Output: Иванов И. И.
}

// ExampleShortName_partial shows that empty fields are simply omitted
// from the result.
func ExampleShortName_partial() {
	fmt.Println(gomorphy.ShortName(gomorphy.FullName{
		Last:  "Иванова",
		First: "Анна",
	}))
	fmt.Println(gomorphy.ShortName(gomorphy.FullName{
		First:      "Иван",
		Patronymic: "Иванович",
	}))
	fmt.Println(gomorphy.ShortName(gomorphy.FullName{
		Last: "Иванов",
	}))
	// Output:
	// Иванова А.
	// И. И.
	// Иванов
}

// ExampleShortName_fromParsed shows the typical pipeline: a free-form
// string is parsed (and reduced to the nominative case) and then turned
// into the short form.
func ExampleShortName_fromParsed() {
	fn, err := gomorphy.ParseFullName("Ивановой Анне Сергеевне")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gomorphy.ShortName(fn))
	// Output: Иванова А. С.
}

// ExampleShortName_inCase shows that ShortName preserves whatever case is
// in its input. Combine it with DeclineFullName to get the short form in
// any case.
func ExampleShortName_inCase() {
	nom := gomorphy.FullName{
		Last:       "Иванова",
		First:      "Анна",
		Patronymic: "Сергеевна",
	}
	gen, err := gomorphy.DeclineFullName(nom, gomorphy.Genitive)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(gomorphy.ShortName(gen))
	// Output: Ивановой А. С.
}
