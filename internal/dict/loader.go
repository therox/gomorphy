package dict

import (
	"fmt"
	"os"
)

// LoadFile loads a .bin dictionary file into memory.
// The whole file is read via Decode (format: magic + gob stream).
func LoadFile(path string) (*Dict, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening dictionary %q: %w", path, err)
	}
	defer f.Close()

	d, err := Decode(f)
	if err != nil {
		return nil, fmt.Errorf("parsing dictionary %q: %w", path, err)
	}
	return d, nil
}

// SaveFile saves the dictionary to a .bin file (creating or overwriting it).
func SaveFile(d *Dict, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating dictionary file %q: %w", path, err)
	}
	defer f.Close()

	if err := d.Encode(f); err != nil {
		return fmt.Errorf("writing dictionary %q: %w", path, err)
	}
	return nil
}

// Lookup returns the interpretations of the word form word from the
// reverse index. Case and the letter "ё" are normalized before lookup.
// If the form is not found, nil is returned without an error.
func (d *Dict) Lookup(word string) []IndexEntry {
	if d == nil || d.FormIndex == nil {
		return nil
	}
	return d.FormIndex[NormalizeForm(word)]
}

// TagString returns the tag string by index, or an empty string when the
// index is out of range.
func (d *Dict) TagString(t Tag) string {
	if int(t) >= len(d.Tags) {
		return ""
	}
	return d.Tags[t]
}
