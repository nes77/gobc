package bencode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	i := "i1337e"

	di, err := DecodeString(i)

	assert.Nil(t, err)

	assert.IsType(t, Integer{}, di)

	bi := di.(Integer).Int()

	assert.NotNil(t, bi)

	assert.Equal(t, int64(1337), bi.Int64())

	s := "4:butt"

	ds, err := DecodeString(s)

	assert.Nil(t, err)

	assert.IsType(t, String(""), ds)

	assert.Equal(t, string(ds.(String)), "butt")

	l := "l4:butti4000ee"

	dl, err := DecodeString(l)

	assert.Nil(t, err)

	assert.IsType(t, List(nil), dl)

	ls := dl.(List)

	assert.Len(t, ls, 2)

	assert.IsType(t, String(""), ls[0])
	assert.IsType(t, Integer{}, ls[1])

	assert.Equal(t, "butt", string(ls[0].(String)))
	assert.Equal(t, int64(4000), ls[1].(Integer).Int().Int64())

	d := "d4:butti4000ee"

	dd, err := DecodeString(d)

	assert.Nil(t, err)
	assert.IsType(t, Dict(nil), dd)
	assert.IsType(t, Integer{}, dd.(Dict)["butt"])

	assert.Equal(t, int64(4000), dd.(Dict)["butt"].(Integer).Int().Int64())

}

func TestEncode(t *testing.T) {
	d := "d4:butti4000e5:butt2i700ee"

	dd, _ := DecodeString(d)

	enc := dd.Bencode()

	assert.Equal(t, d, string(enc))

	l := "l4:butti4000eledeledee"

	dl, _ := DecodeString(l)

	enc = dl.Bencode()

	assert.Equal(t, l, string(enc))

	l = "llllldeeeeee"

	dl, _ = DecodeString(l)

	enc = dl.Bencode()

	assert.Equal(t, l, string(enc))

}

func BenchmarkDecode(b *testing.B) {
	d := "d4:butti4000e5:butt2i700ee"
	var err error

	for i := 0; i < b.N; i++ {
		_, err = DecodeString(d)
	}

	if err != nil {
		b.Fail()
	}

}

func BenchmarkEncode(b *testing.B) {
	d := "d4:butti4000e5:butt2i700ee"
	dd, _ := DecodeString(d)
	var enc []byte

	for i := 0; i < b.N; i++ {
		enc = dd.Bencode()
	}

	if enc == nil {
		b.Fail()
	}
}
