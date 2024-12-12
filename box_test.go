package valuebox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

const misconfiguredTest = "Misconfigured test case"

const (
	getBool    = "GetBool(%s)"
	getFloat64 = "GetFloat64(%s)"
	getInt64   = "GetInt64(%s)"
	getString  = "GetString(%s)"
)

type successTestCase struct {
	setPath       string
	setValue      interface{}
	getFunctionID string
	getPath       string
	expectedValue interface{}
	expectedError error
}

const jsonData = `
{
    "albumTitle": "The Colour and The Shape",
    "artist": "Foo Fighters",
    "price": {
		"regular": 12.35,
		"withMembership": 10.35
	},
    "trackCount": 13,
    "releaseDate": "1997-05-20T18:25:43.511Z",
    "soldOut": true,
	"genre": [
		"alternative rock",
		"post-grunge",
		"hard-rock",
		"grunge",
		"punk rock"
	],
	"singles": [
		{
			"title": "Monkey Wrench",
			"released": "1997-04-28T00:00:00"
		},
		{
			"title": "Everlong",
			"released": "1997-08-18T00:00:00"
		},
		{
			"title": "My Hero",
			"released": "1998-01-19T:00:00:00"
		}
	]
}
`

func runSuccessSetAndRetrieve(testCase successTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		var err error
		var ok bool
		// var i int64
		// var f float64

		b := New()
		b.Set(testCase.setPath, testCase.setValue)

		switch testCase.getFunctionID {
		case getBool:
			var v bool
			var r bool

			r, err = b.GetBool(testCase.getPath)

			if v, ok = testCase.expectedValue.(bool); !ok {
				t.Fatal(misconfiguredTest)
			}

			if v != r {
				t.Errorf("wanted \"%v\", got \"%v\"", v, r)
			}
		case getFloat64:
			var v float64
			var r float64

			r, err = b.GetFloat64(testCase.getPath)

			if f32, ok := testCase.expectedValue.(float32); ok {
				v = float64(f32)
			} else if v, ok = testCase.expectedValue.(float64); !ok {
				t.Fatal(misconfiguredTest)
			}

			if v != r {
				t.Errorf("wanted \"%f\", got \"%f\"", v, r)
			}
		case getInt64:
			var v int64
			var r int64

			r, err = b.GetInt64(testCase.getPath)

			if i, ok := testCase.expectedValue.(int); ok {
				v = int64(i)
			} else if i8, ok := testCase.expectedValue.(int8); ok {
				v = int64(i8)
			} else if i16, ok := testCase.expectedValue.(int16); ok {
				v = int64(i16)
			} else if i32, ok := testCase.expectedValue.(int32); ok {
				v = int64(i32)
			} else if v, ok = testCase.expectedValue.(int64); !ok {
				t.Fatal(misconfiguredTest)
			}

			if v != r {
				t.Errorf("wanted \"%d\", got \"%d\"", v, r)
			}
		case getString:
			var v string
			var r string

			r, err = b.GetString(testCase.getPath)

			if v, ok = testCase.expectedValue.(string); !ok {
				t.Fatal(misconfiguredTest)
			}

			if v != r {
				t.Errorf("wanted \"%s\", got \"%s\"", v, r)
			}
		default:
			t.Fatalf("Not supported functionID: \"%s\"", testCase.getFunctionID)
		}

		if err != testCase.expectedError {
			t.Errorf("wanted error \"%s\", got: \"%s\"", testCase.expectedError, err)
		}
	}
}

// Test new Box initialization.
func TestBoxInit(t *testing.T) {
	_ = New()
}

func TestBoxSet(t *testing.T) {
	b := New()

	b.Set("name", "Andrea")
}

func TestBoxSuccesSetAndRetrieveBasic(t *testing.T) {
	var data map[string]interface{}

	bf := bytes.NewBuffer([]byte(jsonData))
	json.NewDecoder(bf).Decode(&data)

	var testCases = []successTestCase{
		{"name", "Andrea", getString, "name", "Andrea", nil},
		{"someInt", 234, getInt64, "someInt", 234, nil},
		{"someInt8", int8(34), getInt64, "someInt8", int8(34), nil},
		{"someInt16", int16(234), getInt64, "someInt16", int16(234), nil},
		{"someInt32", int32(234), getInt64, "someInt32", int32(234), nil},
		{"someInt64", int64(234), getInt64, "someInt64", int64(234), nil},
		{"someFloat32", float32(234.45), getFloat64, "someFloat32", float32(234.45), nil},
		{"someFloat64", float64(234.45), getFloat64, "someFloat64", float64(234.45), nil},
		{"someBool", true, getBool, "someBool", true, nil},

		{"root", data, getString, "root.albumTitle", "The Colour and The Shape", nil},
		{"root", data, getString, "root.artist", "Foo Fighters", nil},
		{"root", data, getFloat64, "root.price.regular", 12.35, nil},
		{"root", data, getFloat64, "root.price.withMembership", 10.35, nil},
		{"root", data, getFloat64, "root.trackCount", float64(13), nil},
		{"root", data, getString, "root.releaseDate", "1997-05-20T18:25:43.511Z", nil},
		{"root", data, getBool, "root.soldOut", true, nil},
		{"root", data, getString, "root.genre.0", "alternative rock", nil},
		{"root", data, getString, "root.genre.1", "post-grunge", nil},
		{"root", data, getString, "root.genre.2", "hard-rock", nil},
		{"root", data, getString, "root.genre.3", "grunge", nil},
		{"root", data, getString, "root.genre.4", "punk rock", nil},
		{"root", data, getString, "root.singles.0.title", "Monkey Wrench", nil},
		{"root", data, getString, "root.singles.0.released", "1997-04-28T00:00:00", nil},
		{"root", data, getString, "root.singles.1.title", "Everlong", nil},
		{"root", data, getString, "root.singles.1.released", "1997-08-18T00:00:00", nil},
		{"root", data, getString, "root.singles.2.title", "My Hero", nil},
		{"root", data, getString, "root.singles.2.released", "1998-01-19T:00:00:00", nil},
	}

	for _, testCase := range testCases {
		getDescription := fmt.Sprintf(testCase.getFunctionID, testCase.getPath)
		t.Run(fmt.Sprintf("Set(%s,%T)|%s", testCase.setPath, testCase.setValue, getDescription), runSuccessSetAndRetrieve(testCase))
	}
}

func TestBoxErrors(t *testing.T) {
	var data map[string]interface{}

	bf := bytes.NewBuffer([]byte(jsonData))
	json.NewDecoder(bf).Decode(&data)

	var testCases = []successTestCase{
		// ErrorNoValueFound
		{"name", "Andrea", getString, "otherName", "", ErrorNoValueFound("otherName")},
		{"someInt", 234, getInt64, "otherInt", 0, ErrorNoValueFound("otherInt")},
		{"someFloat32", 765.3, getFloat64, "otherFloat", 0.0, ErrorNoValueFound("otherFloat")},
		{"someBool", true, getBool, "otherBool", false, ErrorNoValueFound("otherBool")},
		{"root", data, getString, "root_.albumTitle", "", ErrorNoValueFound("root_.albumTitle")},
		{"root", data, getString, "root.artist_", "", ErrorNoValueFound("root.artist_")},
		{"root", data, getFloat64, "root_.price.regular", 0.0, ErrorNoValueFound("root_.price.regular")},
		{"root", data, getFloat64, "root.price.withMembership_", 0.0, ErrorNoValueFound("root.price.withMembership_")},
		{"root", data, getBool, "root_.soldOut", false, ErrorNoValueFound("root_.soldOut")},
		{"root", data, getBool, "root.soldOut_", false, ErrorNoValueFound("root.soldOut_")},
		{"root", data, getString, "root_.genre.0", "", ErrorNoValueFound("root_.genre.0")},
		{"root", data, getString, "root.genre_.1", "", ErrorNoValueFound("root.genre_.1")},
		{"root", data, getString, "root_.singles.0.title", "", ErrorNoValueFound("root_.singles.0.title")},
		{"root", data, getString, "root.singles_.1.title", "", ErrorNoValueFound("root.singles_.1.title")},
		{"root", data, getString, "root.singles.2.released_", "", ErrorNoValueFound("root.singles.2.released_")},

		// ErrorCantResolveToType
		{"name", 12, getString, "name", "", ErrorCantResolveToType{"string", "name"}},
		{"someInt", "234", getInt64, "someInt", 0, ErrorCantResolveToType{"int64", "someInt"}},
		{"someFloat64", false, getFloat64, "someFloat64", 0.0, ErrorCantResolveToType{"float64", "someFloat64"}},
		{"someBool", 34.21, getBool, "someBool", false, ErrorCantResolveToType{"bool", "someBool"}},

		{"root", data, getFloat64, "root.artist", 0.0, ErrorCantResolveToType{"float64", "root.artist"}},
		{"root", data, getBool, "root.price.regular", false, ErrorCantResolveToType{"bool", "root.price.regular"}},
		{"root", data, getString, "root.soldOut", "", ErrorCantResolveToType{"string", "root.soldOut"}},
		{"root", data, getInt64, "root.genre.0", 0, ErrorCantResolveToType{"int64", "root.genre.0"}},
	}

	for _, testCase := range testCases {
		getDescription := fmt.Sprintf(testCase.getFunctionID, testCase.getPath)
		t.Run(fmt.Sprintf("Set(%s,%T)|%s", testCase.setPath, testCase.setValue, getDescription), runSuccessSetAndRetrieve(testCase))
	}
}
