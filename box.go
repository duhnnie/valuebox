package valuebox

import (
	"encoding/json"
	"strconv"
	"strings"
)

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

func (r *Box) setToMap(m map[string]interface{}, key string, data []byte) error {
	var value interface{}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	m[key] = value
	return nil
}

func (r *Box) setToSlice(s []interface{}, index int, data []byte) error {
	var value interface{}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	s[index] = value
	return nil
}

func (r *Box) set(path []string, key string, data []byte) error {
	stringPath := strings.Join(path, ".")
	parent, err := r.Get(stringPath)

	if err != nil {
		return err
	}

	var index int

	if m, ok := parent.(map[string]interface{}); ok {
		err = r.setToMap(m, key, data)
	} else if s, ok := parent.([]interface{}); !ok {
		return &ResolveError{ErrorCodeNotAMapOrSlice, stringPath, nil}
	} else if index, err = strconv.Atoi(key); err == nil {
		err = r.setToSlice(s, index, data)
	}

	return err
}

func (r *Box) Set(name string, data []byte) error {
	path := strings.Split(name, ".")
	parentPath := path[:len(path)-1]

	if len(parentPath) > 0 {
		key := path[len(path)-1]
		return r.set(parentPath, key, data)
	}

	var value interface{}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	r.values[name] = value
	return nil
}

func (r *Box) Get(valueName string) (res interface{}, err error) {
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
	targetObject, exists := r.values[name]

	if !exists {
		return nil, &ResolveError{ErrorCodeNoValueFound, name, nil}
	}

	res, errPath, errCode := resolve(targetObject, path[1:])

	if errCode != "" {
		return nil, &ResolveError{errCode, strings.Join(append(path[:1], errPath), "."), nil}
	}

	return res, nil
}

func (r *Box) GetFloat64(valueName string) (float64, error) {
	if resolvedValue, err := r.Get(valueName); err != nil {
		return 0, err
	} else if floatValue, ok := resolvedValue.(float64); ok {
		return floatValue, nil
	} else if floatValue, ok := resolvedValue.(float32); ok {
		return float64(floatValue), nil
	} else {
		return 0, &TypeResolvingError{"float64", valueName}
	}
}

func (r *Box) GetBool(valueName string) (bool, error) {
	if resolvedValue, err := r.Get(valueName); err != nil {
		return false, err
	} else if boolValue, ok := resolvedValue.(bool); !ok {
		return false, &TypeResolvingError{"bool", valueName}
	} else {
		return boolValue, nil
	}
}

func (r *Box) GetString(valueName string) (string, error) {
	if resolvedValue, err := r.Get(valueName); err != nil {
		return "", err
	} else if stringValue, ok := resolvedValue.(string); !ok {
		return "", &TypeResolvingError{"string", valueName}
	} else {
		return stringValue, nil
	}
}

func (r *Box) GetSlice(valueName string) ([]interface{}, error) {
	if resolvedValue, err := r.Get(valueName); err != nil {
		return nil, err
	} else if sliceValue, ok := resolvedValue.([]interface{}); !ok {
		return nil, &TypeResolvingError{"[]interface{}", valueName}
	} else {
		return sliceValue, nil
	}
}

func (r *Box) GetMap(valueName string) (map[string]interface{}, error) {
	if resolvedValue, err := r.Get(valueName); err != nil {
		return nil, err
	} else if mapValue, ok := resolvedValue.(map[string]interface{}); !ok {
		return nil, &TypeResolvingError{"map[string]interface{}", valueName}
	} else {
		return mapValue, nil
	}
}
