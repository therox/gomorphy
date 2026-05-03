package gomorphy

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/therox/gomorphy/internal/dict"
)

// Global dictionary state.
// Init populates the pointer once; a second call returns an error.
// If Init has not been called, the first API access checks the
// GOMORPHY_DICT environment variable and loads the dictionary from it.
//
// The complexity is intentionally minimal: a single pointer plus a mutex.
// No LRU, no asynchronous reload.
var (
	dictMu      sync.Mutex
	dictPtr     *dict.Dict
	initInvoked bool
)

// Init loads a .bin dictionary into the global singleton.
// A second call returns an error — reinitialization is not supported.
func Init(path string) error {
	dictMu.Lock()
	defer dictMu.Unlock()

	if initInvoked {
		return errors.New("gomorphy: dictionary is already initialized")
	}
	d, err := dict.LoadFile(path)
	if err != nil {
		return fmt.Errorf("gomorphy: loading dictionary %q: %w", path, err)
	}
	dictPtr = d
	initInvoked = true
	return nil
}

// getDict returns a pointer to the dictionary, loading it automatically
// from the file pointed to by GOMORPHY_DICT when needed.
// Every public API function uses it.
func getDict() (*dict.Dict, error) {
	dictMu.Lock()
	defer dictMu.Unlock()

	if dictPtr != nil {
		return dictPtr, nil
	}
	path := os.Getenv("GOMORPHY_DICT")
	if path == "" {
		return nil, errors.New("gomorphy: dictionary is not initialized, call Init or set GOMORPHY_DICT")
	}
	d, err := dict.LoadFile(path)
	if err != nil {
		return nil, fmt.Errorf("gomorphy: autoload of dictionary %q: %w", path, err)
	}
	dictPtr = d
	initInvoked = true
	return d, nil
}

// useDict is an internal helper for tests: it swaps the dictionary pointer,
// bypassing LoadFile. Not exported.
func useDict(d *dict.Dict) {
	dictMu.Lock()
	defer dictMu.Unlock()
	dictPtr = d
	initInvoked = true
}

// ===== Helpers for working with OpenCorpora tags =====
//
// Tag names follow the OpenCorpora convention: nomn/gent/datv/accs/ablt/loct,
// sing/plur, masc/femn/neut/ms-f, NOUN/ADJF, anim/inan, Pltm, Fixd, Name/Patr/Surn.

// caseTag maps a Case to the corresponding OpenCorpora string tag.
func caseTag(c Case) string {
	switch c {
	case Nominative:
		return "nomn"
	case Genitive:
		return "gent"
	case Dative:
		return "datv"
	case Accusative:
		return "accs"
	case Instrumental:
		return "ablt"
	case Prepositional:
		return "loct"
	}
	return ""
}

// numberTag maps a Number to the corresponding OpenCorpora string tag.
func numberTag(n Number) string {
	switch n {
	case Singular:
		return "sing"
	case Plural:
		return "plur"
	}
	return ""
}

// genderTag maps a Gender to the corresponding OpenCorpora string tag.
func genderTag(g Gender) string {
	switch g {
	case Masculine:
		return "masc"
	case Feminine:
		return "femn"
	case Neuter:
		return "neut"
	case Common:
		return "ms-f"
	}
	return ""
}

// parseCase extracts a Case from an OpenCorpora tag; returns CaseUnknown if it is not a case.
func parseCase(s string) Case {
	switch s {
	case "nomn":
		return Nominative
	case "gent":
		return Genitive
	case "datv":
		return Dative
	case "accs":
		return Accusative
	case "ablt":
		return Instrumental
	case "loct":
		return Prepositional
	}
	return CaseUnknown
}

// parseNumber extracts a Number from an OpenCorpora tag.
func parseNumber(s string) Number {
	switch s {
	case "sing":
		return Singular
	case "plur":
		return Plural
	}
	return NumberUnknown
}

// parseGender extracts a Gender from an OpenCorpora tag.
func parseGender(s string) Gender {
	switch s {
	case "masc":
		return Masculine
	case "femn":
		return Feminine
	case "neut":
		return Neuter
	case "ms-f":
		return Common
	}
	return GenderUnknown
}

// hasTagStr reports whether the tags slice contains the string tag s.
func hasTagStr(d *dict.Dict, tags []dict.Tag, s string) bool {
	for _, t := range tags {
		if d.TagString(t) == s {
			return true
		}
	}
	return false
}

// lemmaPOS determines the part of speech of a lemma from its tags.
func lemmaPOS(d *dict.Dict, tags []dict.Tag) POS {
	for _, t := range tags {
		switch d.TagString(t) {
		case "NOUN":
			return POSNoun
		case "ADJF":
			return POSAdjf
		}
	}
	return POSUnknown
}
