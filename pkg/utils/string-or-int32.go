package utils

type StringOrInt32 struct {
	Type        string
	Int32Value  int32
	StringValue string
}

type StringOrInt32Array []StringOrInt32

func (stringOrInt32Array StringOrInt32Array) ToStringArray() []string {
	var matchedArray []string
	for _, stringOrInt32 := range stringOrInt32Array {
		if stringOrInt32.Type == "string" {
			matchedArray = append(matchedArray, stringOrInt32.StringValue)
		}
	}
	if len(matchedArray) == 0 {
		return nil
	}
	return matchedArray
}

func (stringOrInt32Array StringOrInt32Array) ToInt32Array() []int32 {
	var matchedArray []int32
	for _, stringOrInt32 := range stringOrInt32Array {
		if stringOrInt32.Type == "int32" {
			matchedArray = append(matchedArray, stringOrInt32.Int32Value)
		}
	}
	if len(matchedArray) == 0 {
		return nil
	}
	return matchedArray
}
