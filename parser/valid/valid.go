package valid

import (
	"reflect"
	"time"
)

// Uint check Uint
func Uint(fieldType reflect.Type, v interface{}) (bool, interface{}) {
	var returnValue interface{}
	valid := false

	switch fieldType.Kind() {
	case reflect.Ptr:
		switch fieldType.Elem().Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if value, ok := v.(float64); ok {
				realValue := uint64(value)
				returnValue = &realValue
				valid = true
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value, ok := v.(float64); ok {
			returnValue = uint64(value)
			valid = true
		}
	}
	return valid, returnValue
}

// Int check Int
func Int(fieldType reflect.Type, v interface{}) (bool, interface{}) {
	var returnValue interface{}
	valid := false

	switch fieldType.Kind() {
	case reflect.Ptr:
		{
			switch fieldType.Elem().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if value, ok := v.(float64); ok {
					realValue := int64(value)
					returnValue = &realValue
					valid = true
				}
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value, ok := v.(float64); ok {
			returnValue = int64(value)
			valid = true
		}
	}
	return valid, returnValue
}

// Float check Float
func Float(fieldType reflect.Type, v interface{}) (bool, interface{}) {
	var returnValue interface{}
	valid := false
	switch fieldType.Kind() {
	case reflect.Ptr:
		switch fieldType.Elem().Kind() {
		case reflect.Float32, reflect.Float64:
			if value, ok := v.(float64); ok {
				realValue := value
				returnValue = &realValue
				valid = true
			}
		}
	case reflect.Float32, reflect.Float64:
		if value, ok := v.(float64); ok {
			returnValue = value
			valid = true
		}
	}
	return valid, returnValue
}

// Bool check bool
func Bool(fieldType reflect.Type, v interface{}) (bool, interface{}) {
	var returnValue interface{}
	valid := false
	switch fieldType.Kind() {
	case reflect.Ptr:
		switch fieldType.Elem().Kind() {
		case reflect.Float32, reflect.Float64:
			if value, ok := v.(bool); ok {
				realValue := value
				returnValue = &realValue
				valid = true
			}
		}
	case reflect.Bool:
		if value, ok := v.(float64); ok {
			returnValue = value
			valid = true
		}
	}
	return valid, returnValue
}

// Time check Time
func Time(
	fieldType reflect.Type,
	v interface{},
	layout string,
) (bool, interface{}) {
	var returnValue interface{}
	valid := false

	switch fieldType.Kind() {
	case reflect.Ptr:
		switch fieldType.Elem().Kind() {
		case reflect.Struct:
			if fieldType.Elem() == reflect.TypeOf(time.Time{}) {
				if value, ok := v.(string); ok {
					// if realValue, err := time.Parse(time.RFC3339, value); err == nil {
					// 	returnValue = &realValue
					// 	valid = true
					// } else if realValue, err := time.Parse("2006-01-02T15:04:05.999Z", value); err == nil {
					// 	returnValue = &realValue
					// 	valid = true
					// }
					if realValue, err := time.Parse(layout, value); err == nil {
						returnValue = &realValue
						valid = true
					} else {
						panic(" time format error ")
					}
				}
			}
		}
	case reflect.Struct:
		{
			if fieldType == reflect.PtrTo(reflect.TypeOf(time.Time{})) {
				if value, ok := v.(string); ok {
					// if realValue, err := time.Parse(time.RFC3339, value); err == nil {
					// 	returnValue = realValue
					// 	isMatchValue = true
					// } else if realValue, err := time.Parse("2006-01-02T15:04:05.999Z", value); err == nil {
					// 	returnValue = realValue
					// 	isMatchValue = true
					// }
					if realValue, err := time.Parse(layout, value); err == nil {
						returnValue = realValue
						valid = true
					} else {
						panic(" time format error ")
					}
				}
			}
		}
	}
	return valid, returnValue
}

// String check String
func String(
	fieldType reflect.Type,
	v interface{},
) (bool, interface{}) {
	var returnValue interface{}
	valid := false

	switch fieldType.Kind() {
	case reflect.Ptr:
		switch fieldType.Elem().Kind() {
		case reflect.String:
			if value, ok := v.(string); ok {
				realValue := value
				returnValue = &realValue
				valid = true
			}
		}
	case reflect.String:
		if value, ok := v.(string); ok {
			returnValue = value
			valid = true
		}
	}
	return valid, returnValue
}

func MatchType(fieldType reflect.Type, v interface{}, timeLayout string) (bool, interface{}) {
	var returnValue interface{}
	valid := false

	switch fieldType.Kind() {
	case reflect.Ptr:
		{
			switch fieldType.Elem().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return Int(fieldType, v)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return Uint(fieldType, v)
			case reflect.Float32, reflect.Float64:
				return Float(fieldType, v)
			case reflect.Bool:
				return Bool(fieldType, v)
			case reflect.String:
				return String(fieldType, v)
			case reflect.Struct:
				if fieldType.Elem() == reflect.TypeOf(time.Time{}) {
					return Time(fieldType, v, timeLayout)
				}
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int(fieldType, v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Uint(fieldType, v)
	case reflect.Float32, reflect.Float64:
		return Float(fieldType, v)
	case reflect.Bool:
		return Bool(fieldType, v)
	case reflect.String:
		return String(fieldType, v)
	case reflect.Struct:
		{
			if fieldType == reflect.PtrTo(reflect.TypeOf(time.Time{})) {
				return Time(fieldType, v, timeLayout)
			}
		}

	}

	return valid, returnValue
}
