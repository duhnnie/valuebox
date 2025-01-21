package valuebox

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const misconfiguredTest = "Misconfigured test case"

const (
	get        = "Get(%s)"
	getBool    = "GetBool(%s)"
	getFloat64 = "GetFloat64(%s)"
	getString  = "GetString(%s)"
	getSlice   = "GetSlice(%s)"
	getMap     = "GetMap(%s)"
)

type getTestCase struct {
	box              *Box
	setPath          string
	setValue         []byte
	getFunctionID    string
	getPath          string
	expectedValue    interface{}
	expectedGetError error
	expectedSetError error
}

var jsonData = []byte(`
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
	],
	"nilValue": null
}
`)

func validateErrors(expectedError error, actualError error) func(t *testing.T) {
	return func(t *testing.T) {
		if actualError == nil && actualError == expectedError {
			return
		} else if (actualError == nil || expectedError == nil) && actualError != expectedError {
			t.Errorf("wanted error \"%s\", got: \"%s\"", expectedError, actualError)
		} else if reflect.TypeOf(actualError).Name() != reflect.TypeOf(expectedError).Name() ||
			actualError.Error() != expectedError.Error() {
			t.Errorf("wanted error \"%s\", got: \"%s\"", expectedError, actualError)
		}
	}
}

func runSuccessSetAndRetrieve(testCase getTestCase) func(t *testing.T) {
	return func(t *testing.T) {
		var getError error
		var ok bool

		b := testCase.box

		if testCase.setPath != "" {
			setError := b.Set(testCase.setPath, testCase.setValue)

			t.Run(
				fmt.Sprintf("%s/validateSetError", t.Name()),
				validateErrors(testCase.expectedSetError, setError),
			)
		}

		switch testCase.getFunctionID {
		case get:
			var v = testCase.expectedValue
			var r interface{}

			r, getError = b.Get(testCase.getPath)

			if v != r {
				t.Errorf("wanted \"%v\", got \"%v\"", v, r)
			}
		case getBool:
			var v bool
			var r bool

			r, getError = b.GetBool(testCase.getPath)

			if v, ok = testCase.expectedValue.(bool); !ok {
				panic(misconfiguredTest)
			}

			if v != r {
				t.Errorf("wanted \"%v\", got \"%v\"", v, r)
			}
		case getFloat64:
			var v float64
			var r float64

			r, getError = b.GetFloat64(testCase.getPath)

			if f32, ok := testCase.expectedValue.(float32); ok {
				v = float64(f32)
			} else if v, ok = testCase.expectedValue.(float64); !ok {
				panic(misconfiguredTest)
			}

			if v != r {
				t.Errorf("wanted \"%f\", got \"%f\"", v, r)
			}
		case getString:
			var v string
			var r string

			r, getError = b.GetString(testCase.getPath)

			if v, ok = testCase.expectedValue.(string); !ok {
				panic(misconfiguredTest)
			}

			if v != r {
				t.Errorf("wanted \"%s\", got \"%s\"", v, r)
			}
		default:
			t.Fatalf("Not supported functionID: \"%s\"", testCase.getFunctionID)
		}

		t.Run(
			fmt.Sprintf("%s/validateGetError", t.Name()),
			validateErrors(testCase.expectedGetError, getError),
		)
	}
}

func TestBoxSuccesSetAndRetrieveBasic(t *testing.T) {
	box := New()
	box.Set("root", jsonData)

	var testCases = []getTestCase{
		{box, "", nil, getString, "root.albumTitle", "The Colour and The Shape", nil, nil},
		{box, "", nil, getString, "root.artist", "Foo Fighters", nil, nil},
		{box, "", nil, getFloat64, "root.price.regular", 12.35, nil, nil},
		{box, "", nil, getFloat64, "root.price.withMembership", 10.35, nil, nil},
		{box, "", nil, getFloat64, "root.trackCount", float64(13), nil, nil},
		{box, "", nil, getString, "root.releaseDate", "1997-05-20T18:25:43.511Z", nil, nil},
		{box, "", nil, getBool, "root.soldOut", true, nil, nil},
		{box, "", nil, getString, "root.genre.0", "alternative rock", nil, nil},
		{box, "", nil, getString, "root.genre.1", "post-grunge", nil, nil},
		{box, "", nil, getString, "root.genre.2", "hard-rock", nil, nil},
		{box, "", nil, getString, "root.genre.3", "grunge", nil, nil},
		{box, "", nil, getString, "root.genre.4", "punk rock", nil, nil},
		{box, "", nil, getString, "root.singles.0.title", "Monkey Wrench", nil, nil},
		{box, "", nil, getString, "root.singles.0.released", "1997-04-28T00:00:00", nil, nil},
		{box, "", nil, getString, "root.singles.1.title", "Everlong", nil, nil},
		{box, "", nil, getString, "root.singles.1.released", "1997-08-18T00:00:00", nil, nil},
		{box, "", nil, getString, "root.singles.2.title", "My Hero", nil, nil},
		{box, "", nil, getString, "root.singles.2.released", "1998-01-19T:00:00:00", nil, nil},
		{box, "", nil, get, "root.nilValue", nil, nil, nil},

		// Replacing values
		{box, "root.albumTitle", []byte("\"otherTitle\""), getString, "root.albumTitle", "otherTitle", nil, nil},
		{box, "root.artist", []byte("2"), getFloat64, "root.artist", 2.0, nil, nil},
		{box, "root.singles.2", []byte("true"), getBool, "root.singles.2", true, nil, nil},
	}

	for _, testCase := range testCases {
		getDescription := fmt.Sprintf(testCase.getFunctionID, testCase.getPath)
		t.Run(fmt.Sprintf("Set(%s,%T)|%s", testCase.setPath, testCase.setValue, getDescription), runSuccessSetAndRetrieve(testCase))
	}
}

func TestBoxErrors(t *testing.T) {
	box := New()
	box.Set("root", jsonData)

	var testCases = []getTestCase{
		// ErrorCodeNoValueFound
		{box, "", nil, getString, "root_.albumTitle", "", &ResolveError{ErrorCodeNoValueFound, "root_", nil}, nil},
		{box, "", nil, getString, "root.artist_", "", &ResolveError{ErrorCodeNoValueFound, "root.artist_", nil}, nil},
		{box, "", nil, getFloat64, "root_.price.regular", 0.0, &ResolveError{ErrorCodeNoValueFound, "root_", nil}, nil},
		{box, "", nil, getFloat64, "root.price.withMembership_", 0.0, &ResolveError{ErrorCodeNoValueFound, "root.price.withMembership_", nil}, nil},
		{box, "", nil, getBool, "root_.soldOut", false, &ResolveError{ErrorCodeNoValueFound, "root_", nil}, nil},
		{box, "", nil, getBool, "root.soldOut_", false, &ResolveError{ErrorCodeNoValueFound, "root.soldOut_", nil}, nil},
		{box, "", nil, getString, "root_.genre.0", "", &ResolveError{ErrorCodeNoValueFound, "root_", nil}, nil},
		{box, "", nil, getString, "root.genre_.1", "", &ResolveError{ErrorCodeNoValueFound, "root.genre_", nil}, nil},
		{box, "", nil, getString, "root_.singles.0.title", "", &ResolveError{ErrorCodeNoValueFound, "root_", nil}, nil},
		{box, "", nil, getString, "root.singles_.1.title", "", &ResolveError{ErrorCodeNoValueFound, "root.singles_", nil}, nil},
		{box, "", nil, getString, "root.singles.2.released_", "", &ResolveError{ErrorCodeNoValueFound, "root.singles.2.released_", nil}, nil},

		// ErrorCodeInvalidArrayIndex
		{box, "", nil, getString, "root.genre.c", "", &ResolveError{ErrorCodeNonNumericArrayIndex, "root.genre.c", nil}, nil},
		{box, "", nil, get, "root.singles.nonNumeric", nil, &ResolveError{ErrorCodeNonNumericArrayIndex, "root.singles.nonNumeric", nil}, nil},

		// TypeResolvingError
		{box, "", nil, getString, "root.price.regular", "", &TypeResolvingError{"string", "root.price.regular"}, nil},
		{box, "", nil, getFloat64, "root.soldOut", 0.0, &TypeResolvingError{"float64", "root.soldOut"}, nil},
		{box, "", nil, getBool, "root.genre.0", false, &TypeResolvingError{"bool", "root.genre.0"}, nil},

		// Replacing values
		{box, "root.albumTitle.0", []byte("\"otherTitle\""), getString, "root.albumTitle.0", "", &ResolveError{ErrorCodeNoValueFound, "root.albumTitle.0", nil}, &ResolveError{ErrorCodeNotAMapOrSlice, "root.albumTitle", nil}},
		{box, "root.artist.isBand", []byte("true"), getBool, "root.artist.isBand", false, &ResolveError{ErrorCodeNoValueFound, "root.artist.isBand", nil}, &ResolveError{ErrorCodeNotAMapOrSlice, "root.artist", nil}},
	}

	for _, testCase := range testCases {
		getDescription := fmt.Sprintf(testCase.getFunctionID, testCase.getPath)
		t.Run(fmt.Sprintf("Set(%s,%T)|%s", testCase.setPath, testCase.setValue, getDescription), runSuccessSetAndRetrieve(testCase))
	}
}

func TestComplexReplacements(t *testing.T) {
	box := New()
	box.Set("root", jsonData)

	mySlice := []string{"rock alternativo", "rock de los 90s"}

	mySliceInBytes, err := json.Marshal(mySlice)
	if err != nil {
		panic(err)
	}

	box.Set("root.genre", mySliceInBytes)
	otherSlice, err := box.GetSlice("root.genre")
	if err != nil {
		t.Fatal(err)
	}

	if len(otherSlice) != len(mySlice) {
		t.Fatalf("returned slice doesn't have the same len, expected: %d, got: %d", len(mySlice), len(otherSlice))
	}

	for index, el := range mySlice {
		if el != otherSlice[index].(string) {
			t.Fatalf("Not equal")
		}
	}

	// TODO: test replacement of objects
}

func TestToJSON(t *testing.T) {
	expected := `{"myJSON":{"age":20,"grades":{"english":9,"maths":10,"science":8},"hobbies":["music","programming","outdoors"],"isStudent":true,"name":"Pepito","summary":"Hi, my name is Pepito, I'm 21, my hobbies are music, programming, outdoors, my average grades is 9.00"}}`

	b := New()

	j := []byte(`{
		"name": "Pepito",
		"age": 20
	}`)

	b.Set("myJSON", j)
	b.Set("myJSON.name", []byte("Pepito Watson"))
	b.Set("myJSON.isStudent", []byte("true"))
	b.Set("myJSON.hobbies", []byte(`["music", "programming", "outdoors"]`))
	b.Set("myJSON.grades", []byte(`{"maths": 10, "english": 9, "science": 8}`))

	name, _ := b.GetString("myJSON.name")
	age, _ := b.GetFloat64("myJSON.age")
	hobbies, _ := b.GetStringSlice("myJSON.hobbies")
	grades, _ := b.GetFloat64Map("myJSON.grades")

	gradesTotal := 0.0

	for _, v := range grades {
		gradesTotal += v
	}

	gradesAvg := gradesTotal / float64(len(grades))

	summary := fmt.Sprintf(
		"\"Hi, my name is %s, I'm %d, my hobbies are %s, my average grades is %0.2f\"",
		name,
		int(age)+1,
		strings.Join(hobbies, ", "),
		gradesAvg,
	)

	_ = b.Set("myJSON.summary", []byte(summary))
	j, _ = b.ToJSON()

	if string(j) != expected {
		t.Fatalf("Expected: %s, got: %s", expected, j)
	}
}

func TestNewWithValues(t *testing.T) {
	initialValues := map[string]interface{}{
		"key1": "value1",
		"key2": 2.0,
		"key3": true,
		"key4": []interface{}{"elem1", "elem2"},
		"key5": map[string]interface{}{"subkey1": "subvalue1"},
	}

	box := NewWithValues(initialValues)

	tests := []struct {
		path     string
		expected interface{}
	}{
		{"key1", "value1"},
		{"key2", 2.0},
		{"key3", true},
		{"key4.0", "elem1"},
		{"key4.1", "elem2"},
		{"key5.subkey1", "subvalue1"},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			value, err := box.Get(test.path)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(value, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, value)
			}
		})
	}
}
