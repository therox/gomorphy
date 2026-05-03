package dict

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

// xmlGrammeme is a grammeme (tag) inside <l> or <f>.
type xmlGrammeme struct {
	V string `xml:"v,attr"`
}

// xmlLemmaHead is the contents of <l>: the base form + the lemma's grammemes.
type xmlLemmaHead struct {
	T string        `xml:"t,attr"`
	G []xmlGrammeme `xml:"g"`
}

// xmlForm is the contents of <f>: the word form + the form's grammemes.
type xmlForm struct {
	T string        `xml:"t,attr"`
	G []xmlGrammeme `xml:"g"`
}

// xmlLemma is the entire <lemma> element.
type xmlLemma struct {
	ID    string       `xml:"id,attr"`
	Head  xmlLemmaHead `xml:"l"`
	Forms []xmlForm    `xml:"f"`
}

// BuildStats are counters describing the result of parsing the XML dictionary.
type BuildStats struct {
	// Seen is how many lemmas were encountered overall.
	Seen int
	// Kept is how many lemmas were kept after the POS filter.
	Kept int
	// Forms is how many word forms were added to the index.
	Forms int
}

// BuildFromXML parses the OpenCorpora XML dictionary from r and assembles a Dict.
// Filter: lemmas tagged NOUN (including the subtypes Name/Surn/Patr) and
// ADJF are kept.
func BuildFromXML(r io.Reader) (*Dict, BuildStats, error) {
	b := NewBuilder()
	var stats BuildStats

	dec := xml.NewDecoder(r)
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, stats, fmt.Errorf("parsing XML: %w", err)
		}
		se, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		if se.Name.Local != "lemma" {
			continue
		}

		var lm xmlLemma
		if err := dec.DecodeElement(&lm, &se); err != nil {
			return nil, stats, fmt.Errorf("parsing <lemma>: %w", err)
		}
		stats.Seen++

		if !shouldKeep(lm.Head.G) {
			continue
		}

		lemmaTagStrs := grammemeStrings(lm.Head.G)
		lemma := Lemma{
			Lemma:     strings.ToLower(lm.Head.T),
			LemmaTags: b.InternTags(lemmaTagStrs),
		}

		paradigm := Paradigm{
			Forms: make([]Form, 0, len(lm.Forms)),
		}
		for _, f := range lm.Forms {
			formTagStrs := normalizeCaseAliases(grammemeStrings(f.G))
			paradigm.Forms = append(paradigm.Forms, Form{
				Text: strings.ToLower(f.T),
				Tags: b.InternTags(formTagStrs),
			})
			stats.Forms++
		}
		lemma.Paradigm = paradigm

		b.AddLemma(lemma)
		stats.Kept++
	}

	return b.Dict(), stats, nil
}

// BuildSmall parses the XML at xmlPath and returns a Dict.
// A convenience wrapper for tests: lets you build a mini dictionary from
// testdata without a full OpenCorpora dump.
func BuildSmall(xmlPath string) (*Dict, error) {
	f, err := os.Open(xmlPath)
	if err != nil {
		return nil, fmt.Errorf("opening XML %q: %w", xmlPath, err)
	}
	defer f.Close()

	d, _, err := BuildFromXML(f)
	if err != nil {
		return nil, fmt.Errorf("parsing XML %q: %w", xmlPath, err)
	}
	return d, nil
}

// shouldKeep decides whether to keep the lemma in the dictionary.
// We keep NOUN (including Name/Surn/Patr through subtags) and ADJF.
func shouldKeep(gs []xmlGrammeme) bool {
	for _, g := range gs {
		switch g.V {
		case "NOUN", "ADJF":
			return true
		}
	}
	return false
}

// grammemeStrings extracts the v= attributes as a slice of strings.
func grammemeStrings(gs []xmlGrammeme) []string {
	if len(gs) == 0 {
		return nil
	}
	out := make([]string, len(gs))
	for i, g := range gs {
		out[i] = g.V
	}
	return out
}

// normalizeCaseAliases maps the OpenCorpora "numbered" case aliases
// (gen1/loc1/acc1) onto the canonical tags (gent/loct/accs).
//
// In OpenCorpora the gen1 tag is set on the canonical genitive form for
// words that have a separate partitive gen2 form ("чашка чаю").
// The same goes for loc1/loc2 (prepositional/locative) and acc1/acc2
// (accusative). The partitive/locative (gen2/loc2/acc2) are left as is —
// they carry different morphological meanings and Decline must not
// return them.
func normalizeCaseAliases(tags []string) []string {
	for i, s := range tags {
		switch s {
		case "gen1":
			tags[i] = "gent"
		case "loc1":
			tags[i] = "loct"
		case "acc1":
			tags[i] = "accs"
		}
	}
	return tags
}
