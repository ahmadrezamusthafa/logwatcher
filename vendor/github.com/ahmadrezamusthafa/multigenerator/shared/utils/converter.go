package utils

import (
	"github.com/ahmadrezamusthafa/multigenerator/shared/consts"
	"strconv"
	"time"
)

func StringToBool(value string) bool {
	switch value {
	case "t", "true":
		return true
	default:
		return false
	}
}

func StringToFloat64(value string) float64 {
	var floatValue float64
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return floatValue
}

func StringToTime(value string) time.Time {
	var timeValue time.Time
	timeValue, err := time.Parse(consts.DateTimeFormat, value)
	if err != nil {
		return time.Time{}
	}
	return timeValue
}

func InterfacePtrToInt64(input interface{}) int64 {
	if val, ok := input.(*int64); ok {
		return *val
	}
	res := InterfacePtrToInt(input)
	return int64(res)
}

func InterfacePtrToInt(input interface{}) int {
	if val, ok := input.(*int); ok {
		return *val
	}
	return 0
}

func InterfaceToInt64(input interface{}) int64 {
	if val, ok := input.(int64); ok {
		return val
	}
	res := InterfaceToInt(input)
	return int64(res)
}

func InterfaceToInt(input interface{}) int {
	if val, ok := input.(int); ok {
		return val
	}
	return 0
}

func InterfacePtrToFloat64(input interface{}) float64 {
	if val, ok := input.(*float64); ok {
		return *val
	}
	res := InterfacePtrToFloat32(input)
	return float64(res)
}

func InterfacePtrToFloat32(input interface{}) float32 {
	if val, ok := input.(*float32); ok {
		return *val
	}
	return 0
}

func InterfaceToFloat64(input interface{}) float64 {
	if val, ok := input.(float64); ok {
		return val
	}
	res := InterfaceToFloat32(input)
	return float64(res)
}

func InterfaceToFloat32(input interface{}) float32 {
	if val, ok := input.(float32); ok {
		return val
	}
	return 0
}
