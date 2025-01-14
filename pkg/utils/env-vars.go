package utils

import (
	"os"
	"strconv"
	"strings"
)

func toStringOrInt32(stringValue string) *StringOrInt32 {
	stringValue = strings.Trim(stringValue, CutSet)
	int32Value, err := strconv.Atoi(stringValue)
	if stringValue == "" {
		return nil
	}
	if err == nil {
		return &StringOrInt32{Type: "int32", Int32Value: int32(int32Value)}
	} else {
		return &StringOrInt32{Type: "string", StringValue: stringValue}
	}
}

func GetSingleEnvVar(envName string) (*StringOrInt32, bool) {
	envValue, exists := os.LookupEnv(envName)
	if exists {
		stringOrInt32 := toStringOrInt32(envValue)
		if stringOrInt32 != nil {
			return stringOrInt32, true
		}
	}
	return nil, false
}

func GetMultiEnvVar(envName string) (*StringOrInt32Array, bool) {
	var stringOrInt32Array StringOrInt32Array
	envValue, exists := os.LookupEnv(envName)
	if exists {
		envParts := strings.Split(envValue, ";")
		for _, value := range envParts {
			stringOrInt32 := toStringOrInt32(value)
			if stringOrInt32 != nil {
				stringOrInt32Array = append(stringOrInt32Array, *stringOrInt32)
			}
		}
	}
	if len(stringOrInt32Array) > 0 {
		return &stringOrInt32Array, true
	}
	return nil, false
}
