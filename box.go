package valuebox

import (
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

func Resolve(target, path interface{}) (res interface{}, err error) {
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

	var newPath []string

	if target == nil {
		return nil, ErrorTargetIsNIL
	}

	if pathArray, ok := path.([]string); !ok {
		if pathString, ok := path.(string); !ok {
			return nil, ErrorResolveInvalidParams
		} else if len(pathString) == 0 {
			return target, nil
		} else {
			newPath = strings.Split(pathString, ".")
		}
	} else if len(pathArray) == 0 {
		return target, nil
	} else {
		newPath = pathArray
	}

	if targetMap, ok := target.(map[string]interface{}); ok {
		newTarget := targetMap[newPath[0]]
		return Resolve(newTarget, newPath[1:])
	} else if targetMap, ok := target.([]interface{}); ok {
		index, err := strconv.Atoi(newPath[0])

		if err != nil {
			return nil, ErrorInvalidArrayIndex(newPath[0])
		}

		newTarget := targetMap[index]
		return Resolve(newTarget, newPath[1:])
	} else {
		return nil, ErrorResolveInvalidFirstParam
	}
}

func (r *Box) Set(name string, value interface{}) {
	r.values[name] = value
}

func (r *Box) Get(valueName string) (interface{}, error) {
	path := strings.Split(valueName, ".")
	name := path[0]
	targetObject := r.values[name]

	if targetObject == nil {
		return nil, ErrorNoValueFound(valueName)
	}

	res, err := Resolve(targetObject, path[1:])

	if err != nil && err == ErrorTargetIsNIL {
		return nil, ErrorNoValueFound(valueName)
	}

	return res, err
}

func (r *Box) GetInt64(valueName string) (int64, error) {
	if resolvedValue, err := r.Get(valueName); err != nil {
		return 0, err
	} else if intValue, ok := resolvedValue.(int64); ok {
		return intValue, nil
	} else if intValue, ok := resolvedValue.(int32); ok {
		return int64(intValue), nil
	} else if intValue, ok := resolvedValue.(int16); ok {
		return int64(intValue), nil
	} else if intValue, ok := resolvedValue.(int8); ok {
		return int64(intValue), nil
	} else if intValue, ok := resolvedValue.(int); ok {
		return int64(intValue), nil
	} else {
		return 0, ErrorCantResolveToType{"int64", valueName}
	}
}

func (r *Box) GetFloat64(valueName string) (float64, error) {
	if resolvedValue, err := r.Get(valueName); err != nil {
		return 0, err
	} else if floatValue, ok := resolvedValue.(float64); ok {
		return floatValue, nil
	} else if floatValue, ok := resolvedValue.(float32); ok {
		return float64(floatValue), nil
	} else {
		return 0, ErrorCantResolveToType{"float64", valueName}
	}
}

func (r *Box) GetBool(valueName string) (bool, error) {
	if resolvedValue, err := r.Get(valueName); err != nil {
		return false, err
	} else if boolValue, ok := resolvedValue.(bool); !ok {
		return false, ErrorCantResolveToType{"bool", valueName}
	} else {
		return boolValue, nil
	}
}

func (r *Box) GetString(valueName string) (string, error) {
	if resolvedValue, err := r.Get(valueName); err != nil {
		return "", err
	} else if stringValue, ok := resolvedValue.(string); !ok {
		return "", ErrorCantResolveToType{"string", valueName}
	} else {
		return stringValue, nil
	}
}
