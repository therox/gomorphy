// Package dict contains the internal dictionary types and serialization.
// The binary dictionary format is encoding/gob.
package dict

import (
	"encoding/gob"
	"fmt"
	"io"
)

// magic is the .bin format signature, used to distinguish our file from a
// foreign gob.
const magic = "MORPHGO\x01"

// Tag is a compact identifier of a grammatical feature
// (part of speech, case, number, gender, animacy, the Name/Surn/Patr
// subtypes, etc.). Full OpenCorpora tag strings live in the dictionary's
// shared Tags table and are addressed by index to save memory.
type Tag uint16

// Form is a single word form of a lemma (e.g. "аппетита" for the lemma
// "аппетит").
type Form struct {
	// Text is the actual word form string in lower case; "ё" is preserved as is.
	Text string
	// Tags is the set of grammatical tag indexes into Dict.Tags.
	Tags []Tag
}

// Paradigm is the complete set of word forms of a single lemma.
// For indeclinable words it contains a single form equal to the lemma.
type Paradigm struct {
	// Forms is all forms (including the base form).
	Forms []Form
}

// Lemma is a lexeme: base form + grammatical features of the lemma + paradigm.
type Lemma struct {
	// ID is the lemma's sequence number in the dictionary (starting at 0).
	// Used in the reverse index.
	ID uint32
	// Lemma is the base form (nominative singular) in lower case.
	Lemma string
	// LemmaTags are tags of the lemma (part of speech, gender, animacy,
	// Name/Surn/Patr, etc.).
	LemmaTags []Tag
	// Paradigm is the inflectional paradigm. In Phase 1 it is filled by
	// the parser; in Phase 2 it will be used by the inflection functions.
	Paradigm Paradigm
}

// IndexEntry is a record in the "form → lemma" reverse index.
// FormTags holds features of the specific word form (case, number, etc.),
// while LemmaTags is reachable via Dict.Lemmas[LemmaID].
type IndexEntry struct {
	// LemmaID is a reference into Dict.Lemmas.
	LemmaID uint32
	// FormTags are tags of the specific word form.
	FormTags []Tag
}

// Dict is the dictionary, loaded into memory in its entirety.
type Dict struct {
	// Tags is the shared table of OpenCorpora tag strings.
	// The index in this slice is used as a Tag.
	Tags []string
	// Lemmas is every lemma in the dictionary.
	Lemmas []Lemma
	// FormIndex is the reverse index "normalized form → list of interpretations".
	// The key is normalized: lower-case + replacement of "ё" → "е" (see
	// NormalizeForm).
	FormIndex map[string][]IndexEntry
}

// New creates an empty dictionary with initialized maps.
func New() *Dict {
	return &Dict{
		Tags:      nil,
		Lemmas:    nil,
		FormIndex: make(map[string][]IndexEntry),
	}
}

// Encode serializes the dictionary into w.
// Format: magic signature + gob stream.
func (d *Dict) Encode(w io.Writer) error {
	if _, err := io.WriteString(w, magic); err != nil {
		return fmt.Errorf("writing signature: %w", err)
	}
	enc := gob.NewEncoder(w)
	if err := enc.Encode(d); err != nil {
		return fmt.Errorf("gob encoding: %w", err)
	}
	return nil
}

// Decode reads a dictionary from r, verifying the signature.
func Decode(r io.Reader) (*Dict, error) {
	buf := make([]byte, len(magic))
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, fmt.Errorf("reading signature: %w", err)
	}
	if string(buf) != magic {
		return nil, fmt.Errorf("invalid dictionary file signature")
	}
	d := &Dict{}
	dec := gob.NewDecoder(r)
	if err := dec.Decode(d); err != nil {
		return nil, fmt.Errorf("gob decoding: %w", err)
	}
	if d.FormIndex == nil {
		d.FormIndex = make(map[string][]IndexEntry)
	}
	return d, nil
}

// Builder helps accumulate lemmas and tags while deduplicating tag strings.
type Builder struct {
	dict   *Dict
	tagIdx map[string]Tag
}

// NewBuilder creates a fresh builder.
func NewBuilder() *Builder {
	return &Builder{
		dict:   New(),
		tagIdx: make(map[string]Tag),
	}
}

// InternTag returns the Tag for a tag string, adding it to the table if needed.
func (b *Builder) InternTag(s string) Tag {
	if t, ok := b.tagIdx[s]; ok {
		return t
	}
	t := Tag(len(b.dict.Tags))
	b.dict.Tags = append(b.dict.Tags, s)
	b.tagIdx[s] = t
	return t
}

// InternTags converts a slice of tag strings into a slice of Tag.
func (b *Builder) InternTags(ss []string) []Tag {
	out := make([]Tag, len(ss))
	for i, s := range ss {
		out[i] = b.InternTag(s)
	}
	return out
}

// AddLemma adds a lemma to the dictionary, indexing its forms in
// FormIndex along the way. Returns the ID assigned to the added lemma.
func (b *Builder) AddLemma(l Lemma) uint32 {
	l.ID = uint32(len(b.dict.Lemmas))
	b.dict.Lemmas = append(b.dict.Lemmas, l)
	for _, f := range l.Paradigm.Forms {
		key := NormalizeForm(f.Text)
		b.dict.FormIndex[key] = append(b.dict.FormIndex[key], IndexEntry{
			LemmaID:  l.ID,
			FormTags: f.Tags,
		})
	}
	return l.ID
}

// Dict returns the accumulated dictionary.
func (b *Builder) Dict() *Dict {
	return b.dict
}

// NormalizeForm reduces a word form to a reverse-index key:
// lower case + replacement of "ё"/"Ё" with "е". This makes "ёж" findable
// by the query "еж" and vice versa, and ignores case.
func NormalizeForm(s string) string {
	// Convert to lower case and replace "ё"→"е" in a single pass.
	// We use []rune since the string is UTF-8.
	rs := []rune(s)
	for i, r := range rs {
		switch r {
		case 'Ё':
			rs[i] = 'е'
		case 'ё':
			rs[i] = 'е'
		default:
			// Plain ASCII tolower does not work for Cyrillic;
			// the А–Я / A–Z range is handled explicitly.
			switch {
			case r >= 'A' && r <= 'Z':
				rs[i] = r + ('a' - 'A')
			case r >= 'А' && r <= 'Я':
				rs[i] = r + ('а' - 'А')
			}
		}
	}
	return string(rs)
}
