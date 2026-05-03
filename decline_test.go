package gomorphy

import (
	"strings"
	"testing"
)

// declineCase — one expectation row: input word (lemma), case, number, expected form.
type declineCase struct {
	word string
	c    Case
	n    Number
	want string
}

// TestDecline covers 6 cases × 2 numbers for seven nouns
// of different declension types plus pluralia tantum.
func TestDecline(t *testing.T) {
	cases := []declineCase{
		// аппетит — masculine, ends in a consonant (1st declension).
		{"аппетит", Nominative, Singular, "аппетит"},
		{"аппетит", Genitive, Singular, "аппетита"},
		{"аппетит", Dative, Singular, "аппетиту"},
		{"аппетит", Accusative, Singular, "аппетит"},
		{"аппетит", Instrumental, Singular, "аппетитом"},
		{"аппетит", Prepositional, Singular, "аппетите"},
		{"аппетит", Nominative, Plural, "аппетиты"},
		{"аппетит", Genitive, Plural, "аппетитов"},
		{"аппетит", Dative, Plural, "аппетитам"},
		{"аппетит", Accusative, Plural, "аппетиты"},
		{"аппетит", Instrumental, Plural, "аппетитами"},
		{"аппетит", Prepositional, Plural, "аппетитах"},

		// книга — feminine, ends in -а.
		{"книга", Nominative, Singular, "книга"},
		{"книга", Genitive, Singular, "книги"},
		{"книга", Dative, Singular, "книге"},
		{"книга", Accusative, Singular, "книгу"},
		{"книга", Instrumental, Singular, "книгой"},
		{"книга", Prepositional, Singular, "книге"},
		{"книга", Nominative, Plural, "книги"},
		{"книга", Genitive, Plural, "книг"},
		{"книга", Dative, Plural, "книгам"},
		{"книга", Accusative, Plural, "книги"},
		{"книга", Instrumental, Plural, "книгами"},
		{"книга", Prepositional, Plural, "книгах"},

		// окно — neuter, ends in -о.
		{"окно", Nominative, Singular, "окно"},
		{"окно", Genitive, Singular, "окна"},
		{"окно", Dative, Singular, "окну"},
		{"окно", Accusative, Singular, "окно"},
		{"окно", Instrumental, Singular, "окном"},
		{"окно", Prepositional, Singular, "окне"},
		{"окно", Nominative, Plural, "окна"},
		{"окно", Genitive, Plural, "окон"},
		{"окно", Dative, Plural, "окнам"},
		{"окно", Accusative, Plural, "окна"},
		{"окно", Instrumental, Plural, "окнами"},
		{"окно", Prepositional, Plural, "окнах"},

		// путь — heteroclitic.
		{"путь", Nominative, Singular, "путь"},
		{"путь", Genitive, Singular, "пути"},
		{"путь", Dative, Singular, "пути"},
		{"путь", Accusative, Singular, "путь"},
		{"путь", Instrumental, Singular, "путём"},
		{"путь", Prepositional, Singular, "пути"},
		{"путь", Nominative, Plural, "пути"},
		{"путь", Genitive, Plural, "путей"},
		{"путь", Dative, Plural, "путям"},
		{"путь", Accusative, Plural, "пути"},
		{"путь", Instrumental, Plural, "путями"},
		{"путь", Prepositional, Plural, "путях"},

		// время — heteroclitic, ends in -мя.
		{"время", Nominative, Singular, "время"},
		{"время", Genitive, Singular, "времени"},
		{"время", Dative, Singular, "времени"},
		{"время", Accusative, Singular, "время"},
		{"время", Instrumental, Singular, "временем"},
		{"время", Prepositional, Singular, "времени"},
		{"время", Nominative, Plural, "времена"},
		{"время", Genitive, Plural, "времён"},
		{"время", Dative, Plural, "временам"},
		{"время", Accusative, Plural, "времена"},
		{"время", Instrumental, Plural, "временами"},
		{"время", Prepositional, Plural, "временах"},

		// дитя — heteroclitic with a suppletive plural.
		{"дитя", Nominative, Singular, "дитя"},
		{"дитя", Genitive, Singular, "дитяти"},
		{"дитя", Dative, Singular, "дитяти"},
		{"дитя", Accusative, Singular, "дитя"},
		{"дитя", Instrumental, Singular, "дитятей"},
		{"дитя", Prepositional, Singular, "дитяти"},
		{"дитя", Nominative, Plural, "дети"},
		{"дитя", Genitive, Plural, "детей"},
		{"дитя", Dative, Plural, "детям"},
		{"дитя", Accusative, Plural, "детей"},
		{"дитя", Instrumental, Plural, "детьми"},
		{"дитя", Prepositional, Plural, "детях"},

		// мать — 3rd declension ending in -ь, animate.
		{"мать", Nominative, Singular, "мать"},
		{"мать", Genitive, Singular, "матери"},
		{"мать", Dative, Singular, "матери"},
		{"мать", Accusative, Singular, "мать"},
		{"мать", Instrumental, Singular, "матерью"},
		{"мать", Prepositional, Singular, "матери"},
		{"мать", Nominative, Plural, "матери"},
		{"мать", Genitive, Plural, "матерей"},
		{"мать", Dative, Plural, "матерям"},
		{"мать", Accusative, Plural, "матерей"},
		{"мать", Instrumental, Plural, "матерями"},
		{"мать", Prepositional, Plural, "матерях"},

		// ножницы — pluralia tantum, plural fully covered.
		{"ножницы", Nominative, Plural, "ножницы"},
		{"ножницы", Genitive, Plural, "ножниц"},
		{"ножницы", Dative, Plural, "ножницам"},
		{"ножницы", Accusative, Plural, "ножницы"},
		{"ножницы", Instrumental, Plural, "ножницами"},
		{"ножницы", Prepositional, Plural, "ножницах"},
	}

	for _, tc := range cases {
		name := tc.word + "/" + caseTag(tc.c) + "/" + numberTag(tc.n)
		t.Run(name, func(t *testing.T) {
			got, err := Decline(tc.word, tc.c, tc.n)
			if err != nil {
				t.Fatalf("Decline(%q, %v, %v) error: %v", tc.word, tc.c, tc.n, err)
			}
			if got != tc.want {
				t.Fatalf("Decline(%q, %v, %v) = %q, expected %q",
					tc.word, tc.c, tc.n, got, tc.want)
			}
		})
	}
}

// TestDeclinePluraliaTantumSingularError — for pluralia tantum every singular
// form must return an error with a clear message.
func TestDeclinePluraliaTantumSingularError(t *testing.T) {
	cases := []Case{Nominative, Genitive, Dative, Accusative, Instrumental, Prepositional}
	for _, c := range cases {
		t.Run("ножницы/sing/"+caseTag(c), func(t *testing.T) {
			_, err := Decline("ножницы", c, Singular)
			if err == nil {
				t.Fatalf("Decline(\"ножницы\", %v, Singular) should have returned an error", c)
			}
			if !strings.Contains(err.Error(), "pluralia tantum") {
				t.Fatalf("Decline(\"ножницы\", %v, Singular): expected mention of pluralia tantum, got %v", c, err)
			}
		})
	}
}

// TestDeclineWordNotFound — a word outside the dictionary returns an error.
func TestDeclineWordNotFound(t *testing.T) {
	_, err := Decline("абракадабра", Genitive, Singular)
	if err == nil {
		t.Fatal("expected an error for an unknown word")
	}
}
