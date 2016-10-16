package bencode

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
)

// Bencodable represents types that can be Bencoded
type Bencodable interface {
	Bencode() []byte
}

// Integer represents a bencodable Integer
type Integer struct {
	i *big.Int
}

// NewInt64 creates a new Integer from an int64
func NewInt64(i int64) Integer {
	return Integer{i: big.NewInt(i)}
}

// NewUint64 creates a new Integer from a uint64
func NewUint64(i uint64) Integer {
	b := big.NewInt(0)

	b.SetUint64(i)

	return Integer{i: b}
}

// NewIntFromString creates a new Integer from a string with the given base
func NewIntFromString(s string, base int) (Integer, error) {
	b := big.NewInt(0)

	b, ok := b.SetString(s, base)

	if ok {
		return Integer{i: b}, nil
	}

	return Integer{}, nil
}

// Int returns the integer represented by an Integer
func (i Integer) Int() *big.Int {
	o := big.NewInt(0)

	o.Set(i.i)

	return o
}

// String represents a bencodeable string
type String string

// List represents a list of bencodable elements
type List []Bencodable

// Dict represents a dictionary of bencodable strings to bencodeable elements
type Dict map[string]Bencodable

// Bencode converts an Integer into its bencoded form
func (i Integer) Bencode() []byte {
	s := fmt.Sprintf("i%se", i.i)

	return []byte(s)
}

func (i Integer) String() string {
	return i.i.String()
}

// Bencode converts a String into its bencoded form
func (s String) Bencode() []byte {
	return []byte(fmt.Sprintf("%d:%s", len(s), s))
}

func (s String) String() string {
	return string(s)
}

// Bencode converts a List into its bencoded form
func (l List) Bencode() []byte {
	var buf bytes.Buffer

	buf.WriteRune('l')

	for _, e := range l {
		buf.Write(e.Bencode())
	}

	buf.WriteRune('e')

	return buf.Bytes()
}

func (l List) String() string {
	return string(l.Bencode())
}

// Bencode converts a Dict into its bencoded form
func (d Dict) Bencode() []byte {
	var buf bytes.Buffer

	buf.WriteRune('d')

	var keys []string

	for k := range d {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		buf.Write(String(k).Bencode())
		buf.Write(d[k].Bencode())
	}

	buf.WriteRune('e')

	return buf.Bytes()
}

func (d Dict) String() string {
	return string(d.Bencode())
}
