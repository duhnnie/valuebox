package valuebox

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

type jsonTypes interface {
	float64 | string | bool | interface{}
}

type jsonArrayTypes interface {
	[]float64 | []string | []bool | []interface{}
}

type jsonObjectTypes interface {
	map[string]interface{} | map[string]float64 | map[string]bool | map[string]string
}

type allJsonTypes interface {
	jsonTypes | jsonArrayTypes | jsonObjectTypes
}

type Box struct {
	values map[string]interface{}
}

func New() *Box {
	return &Box{
		values: make(map[string]interface{}),
	}
}

func deleteme(path string) func(interface{}, string, ErrorCode) (interface{}, string, ErrorCode) {
	return func(res interface{}, errPath string, errCode ErrorCode) (interface{}, string, ErrorCode) {
		if errCode != "" && path != "" {
			errPath = strings.Join([]string{path, errPath}, ".")
		}

		return res, errPath, errCode
	}
}

func resolve(target interface{}, path []string) (res interface{}, errPath string, errCode ErrorCode) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*ResolveError); ok {
				e.Path = strings.Join(append(path[:1], e.Path), ".")
				panic(e)
			} else if e, ok := r.(error); ok {
				err := &ResolveError{ErrorCodeOther, path[0], e}
				panic(err)
			} else {
				panic(r)
			}
		}
	}()

	if len(path) == 0 {
		return target, "", ""
	}

	var currentPath = path[0]

	if targetMap, ok := target.(map[string]interface{}); ok {
		if newTarget, exists := targetMap[path[0]]; !exists {
			return nil, currentPath, ErrorCodeNoValueFound
		} else {
			return deleteme(currentPath)(resolve(newTarget, path[1:]))
		}
	} else if targetMap, ok := target.([]interface{}); ok {
		index, err := strconv.Atoi(path[0])

		if err != nil {
			return nil, currentPath, ErrorCodeNonNumericArrayIndex
		}

		newTarget := targetMap[index]
		return deleteme(currentPath)(resolve(newTarget, path[1:]))
	} else {
		return nil, currentPath, ErrorCodeNoValueFound
	}
}

func internalGet(b *Box, valueName string) (res interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				res = nil
				err = e
			} else {
				panic(r)
			}
		}
	}()

	path := strings.Split(valueName, ".")
	name := path[0]
	targetObject, exists := b.values[name]

	if !exists {
		return nil, &ResolveError{ErrorCodeNoValueFound, name, nil}
	}

	res, errPath, errCode := resolve(targetObject, path[1:])

	if errCode != "" {
		return nil, &ResolveError{errCode, strings.Join(append(path[:1], errPath), "."), nil}
	}

	return res, nil
}

func internalGenericGet[T jsonTypes](b *Box, valueName string) (res T, err error) {
	var value T

	if resolvedValue, err := internalGet(b, valueName); err != nil {
		return value, err
	} else if value, ok := resolvedValue.(T); !ok {
		return value, &TypeResolvingError{reflect.TypeOf(value).String(), valueName}
	} else {
		return value, nil
	}
}

func toConcreteSlice[T jsonTypes](b *Box, valueName string) ([]T, error) {
	if resolvedValue, err := internalGet(b, valueName); err != nil {
		return nil, err
	} else if slice, ok := resolvedValue.([]interface{}); !ok {
		return nil, &TypeResolvingError{reflect.TypeOf(slice).String(), valueName}
	} else {
		var concreteSlice []T

		for i, v := range slice {
			if value, ok := v.(T); !ok {
				return nil, &TypeResolvingError{reflect.TypeOf(value).String(), strings.Join([]string{valueName, strconv.Itoa(i)}, ".")}
			} else {
				concreteSlice = append(concreteSlice, value)
			}
		}

		return concreteSlice, nil
	}
}

func toConcreteMap[T jsonTypes](b *Box, valueName string) (map[string]T, error) {
	if resolvedValue, err := internalGet(b, valueName); err != nil {
		return nil, err
	} else if m, ok := resolvedValue.(map[string]interface{}); !ok {
		return nil, &TypeResolvingError{reflect.TypeOf(m).String(), valueName}
	} else {
		var concreteMap map[string]T = make(map[string]T)

		for k, v := range m {
			if value, ok := v.(T); !ok {
				return nil, &TypeResolvingError{reflect.TypeOf(value).String(), strings.Join([]string{valueName, k}, ".")}
			} else {
				concreteMap[k] = value
			}
		}

		return concreteMap, nil
	}
}

func (b *Box) setToMap(m map[string]interface{}, key string, data []byte) error {
	var value interface{}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	m[key] = value
	return nil
}

func (b *Box) setToSlice(s []interface{}, index int, data []byte) error {
	var value interface{}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	s[index] = value
	return nil
}

func (b *Box) set(path []string, key string, data []byte) error {
	stringPath := strings.Join(path, ".")
	parent, err := b.Get(stringPath)

	if err != nil {
		return err
	}

	var index int

	if m, ok := parent.(map[string]interface{}); ok {
		err = b.setToMap(m, key, data)
	} else if s, ok := parent.([]interface{}); !ok {
		return &ResolveError{ErrorCodeNotAMapOrSlice, stringPath, nil}
	} else if index, err = strconv.Atoi(key); err == nil {
		err = b.setToSlice(s, index, data)
	}

	return err
}

func (b *Box) Set(name string, data []byte) error {
	path := strings.Split(name, ".")
	parentPath := path[:len(path)-1]

	if len(parentPath) > 0 {
		key := path[len(path)-1]
		return b.set(parentPath, key, data)
	}

	var value interface{}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	b.values[name] = value
	return nil
}

func (b *Box) Get(valueName string) (res interface{}, err error) {
	return internalGet(b, valueName)
}

func (b *Box) GetFloat64(valueName string) (float64, error) {
	return internalGenericGet[float64](b, valueName)
}

func (b *Box) GetBool(valueName string) (bool, error) {
	return internalGenericGet[bool](b, valueName)
}

func (b *Box) GetString(valueName string) (string, error) {
	return internalGenericGet[string](b, valueName)
}

func (b *Box) GetSlice(valueName string) ([]interface{}, error) {
	return internalGenericGet[[]interface{}](b, valueName)
}

func (b *Box) GetFloat64Slice(valueName string) ([]float64, error) {
	return toConcreteSlice[float64](b, valueName)
}

func (b *Box) GetStringSlice(valueName string) ([]string, error) {
	return toConcreteSlice[string](b, valueName)
}

func (b *Box) GetBoolSlice(valueName string) ([]bool, error) {
	return toConcreteSlice[bool](b, valueName)
}

func (b *Box) GetMap(valueName string) (map[string]interface{}, error) {
	return toConcreteMap[interface{}](b, valueName)
}

func (b *Box) GetBoolMap(valueName string) (map[string]bool, error) {
	return toConcreteMap[bool](b, valueName)
}

func (b *Box) GetFloat64Map(valueName string) (map[string]float64, error) {
	return toConcreteMap[float64](b, valueName)
}

func (b *Box) GetStringMap(valueName string) (map[string]string, error) {
	return toConcreteMap[string](b, valueName)
}

func (b *Box) ValueToJSON(valueName string) ([]byte, error) {
	if value, err := b.Get(valueName); err != nil {
		return nil, err
	} else {
		return json.Marshal(value)
	}
}

func (b *Box) ToJSON() ([]byte, error) {
	return json.Marshal(b.values)
}
