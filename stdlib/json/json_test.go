package json_test

import (
	gojson "encoding/json"
	"strings"
	"testing"

	"github.com/ganehag/tengo/v3"
	"github.com/ganehag/tengo/v3/require"
	"github.com/ganehag/tengo/v3/stdlib/json"
)

type ARR = []interface{}
type MAP = map[string]interface{}

// Compile-time check: encoding/json.Number must satisfy tengo's jsonNumber interface.
// If it ever stops doing so, the TestJSONNumber test will fail to compile here first.
var _ interface {
	Int64() (int64, error)
	Float64() (float64, error)
} = gojson.Number("")

func TestJSON(t *testing.T) {
	testJSONEncodeDecode(t, nil)

	testJSONEncodeDecode(t, 0)
	testJSONEncodeDecode(t, 1)
	testJSONEncodeDecode(t, -1)
	testJSONEncodeDecode(t, 1984)
	testJSONEncodeDecode(t, -1984)

	testJSONEncodeDecode(t, 0.0)
	testJSONEncodeDecode(t, 1.0)
	testJSONEncodeDecode(t, -1.0)
	testJSONEncodeDecode(t, 19.84)
	testJSONEncodeDecode(t, -19.84)

	testJSONEncodeDecode(t, "")
	testJSONEncodeDecode(t, "foo")
	testJSONEncodeDecode(t, "foo bar")
	testJSONEncodeDecode(t, "foo \"bar\"")
	// See: https://github.com/d5/tengo/issues/268
	testJSONEncodeDecode(t, "1\u001C04")
	testJSONEncodeDecode(t, "çığöşü")
	testJSONEncodeDecode(t, "ç1\u001C04IĞÖŞÜ")
	testJSONEncodeDecode(t, "错误测试")

	testJSONEncodeDecode(t, true)
	testJSONEncodeDecode(t, false)

	testJSONEncodeDecode(t, ARR{})
	testJSONEncodeDecode(t, ARR{0})
	testJSONEncodeDecode(t, ARR{false})
	testJSONEncodeDecode(t, ARR{1, 2, 3,
		"four", false})
	testJSONEncodeDecode(t, ARR{1, 2, 3,
		"four", false, MAP{"a": 0, "b": "bee", "bool": true}})

	testJSONEncodeDecode(t, MAP{})
	testJSONEncodeDecode(t, MAP{"a": 0})
	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee"})
	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee", "bool": true})

	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee",
		"arr": ARR{1, 2, 3, "four"}})
	testJSONEncodeDecode(t, MAP{"a": 0, "b": "bee",
		"arr": ARR{1, 2, 3, MAP{"a": false, "b": 109.4}}})

	testJSONEncodeDecode(t, MAP{"id1": 7075984636689534001, "id2": 7075984636689534002})
	testJSONEncodeDecode(t, ARR{1e3, 1E7})
}

func TestDecode(t *testing.T) {
	testDecodeError(t, `{`)
	testDecodeError(t, `}`)
	testDecodeError(t, `{}a`)
	testDecodeError(t, `{{}`)
	testDecodeError(t, `{}}`)
	testDecodeError(t, `[`)
	testDecodeError(t, `]`)
	testDecodeError(t, `[]a`)
	testDecodeError(t, `[[]`)
	testDecodeError(t, `[]]`)
	testDecodeError(t, `"`)
	testDecodeError(t, `"abc`)
	testDecodeError(t, `abc"`)
	testDecodeError(t, `.123`)
	testDecodeError(t, `123.`)
	testDecodeError(t, `1.2.3`)
	testDecodeError(t, `'a'`)
	testDecodeError(t, `true, false`)
	testDecodeError(t, `{"a:"b"}`)
	testDecodeError(t, `{a":"b"}`)
	testDecodeError(t, `{"a":"b":"c"}`)
}

// TestJSONNumber verifies that encoding/json.Number values — produced when
// a caller uses json.Decoder.UseNumber() — survive the Tengo round-trip.
func TestJSONNumber(t *testing.T) {
	input := `{"count":42,"price":19.99,"big":7075984636689534001}`

	var raw map[string]interface{}
	dec := gojson.NewDecoder(strings.NewReader(input))
	dec.UseNumber()
	require.NoError(t, dec.Decode(&raw))

	// Convert the Go map (containing json.Number values) to a Tengo object.
	obj, err := tengo.FromInterface(raw)
	require.NoError(t, err)

	// Encode with the Tengo JSON encoder.
	b, err := json.Encode(obj)
	require.NoError(t, err)

	// Decode back and verify the types and values.
	a, err := json.Decode(b)
	require.NoError(t, err)

	m := a.(*tengo.Map).Value
	require.Equal(t, tengo.Int{Value: 42}, m["count"])
	require.Equal(t, tengo.Float{Value: 19.99}, m["price"])
	require.Equal(t, tengo.Int{Value: 7075984636689534001}, m["big"])
}

func testDecodeError(t *testing.T, input string) {
	_, err := json.Decode([]byte(input))
	require.Error(t, err)
}

func testJSONEncodeDecode(t *testing.T, v interface{}) {
	o, err := tengo.FromInterface(v)
	require.NoError(t, err)

	b, err := json.Encode(o)
	require.NoError(t, err)

	a, err := json.Decode(b)
	require.NoError(t, err, string(b))

	vj, err := gojson.Marshal(v)
	require.NoError(t, err)

	aj, err := gojson.Marshal(tengo.ToInterface(a))
	require.NoError(t, err)

	require.Equal(t, vj, aj)
}
