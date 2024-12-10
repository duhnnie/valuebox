package valuebox

import "strings"

type Box struct {
	values map[string]interface{}
}

func NewRepo() *Box {
	return &Box{
		values: make(map[string]interface{}),
	}
}

func Resolve(target, path interface{}) (interface{}, error) {
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

	if targetMap, ok := target.(map[string]interface{}); !ok {
		return nil, ErrorResolveInvalidFirstParam
	} else {
		newTarget := targetMap[newPath[0]]
		return Resolve(newTarget, newPath[1:])
	}
}

func (r *Box) Set(name string, value interface{}) {
	r.values[name] = value
}

func (r *Box) get(valueName string) (interface{}, error) {
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
	if resolvedValue, err := r.GetFloat64(valueName); err == nil {
		return int64(resolvedValue), nil
	} else if resolvedValue, err := r.get(valueName); err != nil {
		return 0, err
	} else if intValue, ok := resolvedValue.(int64); !ok {
		return 0, ErrorCantResolveToType{"int64", valueName}
	} else {
		return intValue, nil
	}
}

func (r *Box) GetFloat64(valueName string) (float64, error) {
	if resolvedValue, err := r.get(valueName); err != nil {
		return 0, err
	} else if floatValue, ok := resolvedValue.(float64); !ok {
		return 0, ErrorCantResolveToType{"float64", valueName}
	} else {
		return floatValue, nil
	}
}

func (r *Box) GetBool(valueName string) (bool, error) {
	if resolvedValue, err := r.get(valueName); err != nil {
		return false, err
	} else if boolValue, ok := resolvedValue.(bool); !ok {
		return false, ErrorCantResolveToType{"bool", valueName}
	} else {
		return boolValue, nil
	}
}

func (r *Box) GetString(valueName string) (string, error) {
	if resolvedValue, err := r.get(valueName); err != nil {
		return "", err
	} else if stringValue, ok := resolvedValue.(string); !ok {
		return "", ErrorCantResolveToType{"string", valueName}
	} else {
		return stringValue, nil
	}
}
