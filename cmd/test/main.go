package main

import (
	"fmt"
	"log"

	"github.com/therox/gomorphy"
)

func main() {
	// The local version of the library is used automatically — this file
	// lives inside the github.com/therox/gomorphy module, and Go takes the
	// package straight from the root without go get / replace.
	//
	// Before running, build the dictionary:
	//   go run ./cmd/builddict -in testdata/sample.xml -out dict-data/dict.bin
	if err := gomorphy.Init("dict-data/dict.bin"); err != nil {
		log.Fatal(err)
	}

	// === DeclineFullName: Nom → any case. ===
	out, err := gomorphy.DeclineFullName(
		gomorphy.FullName{Last: "Ваткова", First: "Элеанора", Patronymic: "Рузвельтовна"},
		gomorphy.Genitive,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out.Last, out.First, out.Patronymic)

	// === Parse: word form → all possible interpretations. ===
	// Decline reconstructs the word form back — round-trip.
	for _, form := range []string{"стали", "толике"} {
		analyses, err := gomorphy.Parse(form)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("analysis of %q:\n", form)
		for _, a := range analyses {
			back, _ := gomorphy.Decline(a.Lemma, a.Case, a.Number)
			fmt.Printf("  form=%s ← lemma=%s POS=%v case=%v number=%v gender=%v animate=%v\n",
				back, a.Lemma, a.POS, a.Case, a.Number, a.Gender, a.Animate)
		}
	}

	// === ParseFullName + DeclineFullName: full name from any case to any. ===
	// Step 1. Parse the string "Ивановой Анне Сергеевне" (dative case):
	// tokenization + role recognition + reduction of every component to Nom.
	nom, err := gomorphy.ParseFullName("Ивановой Анне Сергеевне")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ParseFullName(\"Ивановой Анне Сергеевне\") → nom %s %s %s\n",
		nom.Last, nom.First, nom.Patronymic)

	// Step 2. Inflect the reconstructed Nom into the desired case.
	for _, c := range []gomorphy.Case{
		gomorphy.Nominative, gomorphy.Genitive, gomorphy.Dative,
		gomorphy.Accusative, gomorphy.Instrumental, gomorphy.Prepositional,
	} {
		got, err := gomorphy.DeclineFullName(nom, c)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("  %s: %s %s %s\n", c, got.Last, got.First, got.Patronymic)
	}
}
