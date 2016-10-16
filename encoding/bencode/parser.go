package bencode

import (
    "io"
    "bytes"
    "bufio"
    "fmt"
    "unicode"
    "strconv"
    "math/big"
)

type token int

const (
    badtok token = iota
    end
    dict
    integer
    list
    str
    eof
)

func (t token) String() string {
    switch t {
    case badtok: return "badtok"
    case end: return "end"
    case dict: return "dict"
    case integer: return "integer"
    case str: return "str"
    case eof: return "eof"
    }

    panic("Invalid token")
}

const UnexpectedEnd error = fmt.Errorf("Unexpected end.")
const UnexpectedEOF error = fmt.Errorf("Unexpected EOF.")
const BadDictKey error = fmt.Errorf("Bad Dict Key. Expected string.")

type dataToken struct {
    token token
    data interface{}
}

func (d dataToken) String() string {
    return fmt.Sprintf("Token[%s -> %s]", d.token, d.data)
}

type lexer interface {
    next() (dataToken, error)
}

type lexerImpl struct {
    r bufio.Reader
}

func newLexer(r io.Reader) lexer {
    return &lexerImpl{r:bufio.NewReader(r)}
}

func (l *lexerImpl) next() (dataToken, error) {
    c, err := l.r.ReadByte()

    if err != nil && err != io.EOF {
        return dataToken{}, err
    } else if err == io.EOF {
        return dataToken{token:eof}
    }

    switch c {
    case 'd':
        return dataToken{token:dict}
    case 'l':
        return dataToken{token:list}
    case 'e':
        return dataToken{token:end}
    case 'i':
        return l.parseInt()
    case unicode.IsDigit(rune(c)):
        l.r.UnreadByte()
        return l.parseString()
    }

    return dataToken{}, fmt.Errorf("Invalid input char: %c", c)
}

func (l *lexerImpl) parseInt() (dataToken, error) {
    d, err := l.r.ReadString('e')

    if err != nil {
        return dataToken{}, err
    }

    d = d[:len(d)-1]

    b := big.NewInt(0)
    b, ok := b.SetString(d, 10)

    if ok {
        return dataToken{token:integer,data:b}
    }

    return dataToken{}, fmt.Errorf("Invalid integer.")
}

func (l *lexerImpl) parseString() (dataToken, error) {
    le, err := l.r.ReadString(':')

    if err != nil {
        return dataToken{}, err
    }

    le = le[:len(le)-1]

    sz, err := strconv.ParseUint(le, 10, 64)

    if err != nil {
        return dataToken{}, err
    }

    out := make([]byte, sz, sz)
    l.r.Read(out)

    return dataToken{token:str, data:string(out)}
}

type parser struct {
    lex lexer
}

func (p *parser) parse() (Bencodable, error) {
    lex := p.lex

    t, err := lex.next()

    if err != nil {
        return nil, err
    }

    switch t.token {
    case dict:
        return p.parseDict(), nil
    case list:
        return p.parseList(), nil
    case integer:
        return Integer(t.data.(*big.Int)), nil
    case str:
        return String(t.data.(string)), nil
    case end:
        return nil, UnexpectedEnd
    case eof:
        return nil, UnexpectedEOF
    }

    panic("Should never be here!")
}

func (p *parser) parseList() (List, error) {
    var out List

    n, err := p.parse()

    for err != nil {
        out = append(out, n)
        n, err = p.parse()
    }

    if err.Error() == UnexpectedEnd {
        return out, nil
    }

    return nil, err
}

func (p *parser) parseDict() (Dict, error) {
    out := make(map[string]Bencodable)

    n, err := p.parse()

    if err == UnexpectedEnd {
        return out, nil
    }

    n2, err2 := p.parse()

    for err != nil && err2 != nil {
        switch n := n.(type) {
        default:
            return out, BadDictKey

        case String:
            out[string(n)] = n2
        }

        n, err = p.parse()

        if err == UnexpectedEnd {
            return out, nil
        }

        n2, err2 = p.parse()
    }

    if err != nil {
        return nil, err
    } else if err2 != nil {
        return nil, err2
    }

    panic("Should never reach here!")

}

func Decode(r io.Reader) (Bencodable, error) {
    lex := newLexer(r)

    o, err := parser{lex:lex}.parse()

    if err != nil && o == nil {
        return nil, err
    }

    return o, nil
}

func DecodeBytes(b []byte) (Bencodable, error) {
    return Decode(bytes.NewReader(b))
}