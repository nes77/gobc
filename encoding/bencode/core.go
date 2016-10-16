package bencode

import (
    "fmt"
    "bytes"
    "math/big"
    "sort"
)

type Bencodable interface {
    Bencode() []byte
}

type Integer *big.Int
type String string
type List []Bencodable
type Dict map[string]Bencodable

func (i Integer) Bencode() []byte {
    s := fmt.Sprintf("i%se", i)

    return []byte(s)
}

func (s String) Bencode() []byte {
    return []byte(fmt.Sprintf("%d:%s", len(s), s))
}

func (l List) Bencode() []byte {
    var buf bytes.Buffer

    buf.WriteRune('l')

    for e := range l {
        buf.Write(e.Bencode())
    }

    buf.WriteRune('e')

    return buf.Bytes()
}

func (d Dict) Bencode() []byte {
    var buf bytes.Buffer

    buf.WriteRune('d')

    var keys []string

    for k := range d {
        keys = append(keys, k)
    }

    sort.Strings(keys)

    for k := range keys {
        buf.Write(String(k).Bencode())
        buf.Write(d[k].Bencode())
    }

    buf.WriteRune('e')

    return buf.Bytes()
}

