package nordvpnapiv1

import (
	"go-nordvpn/pkg/utils"
)

func (s ServerStatus) Equals(t ServerStatus) bool {
	return utils.PrefixedCaseInsensitiveCompareStrings(string(s), string(t))
}

func (s ServerIPType) Equals(t ServerIPType) bool {
	return utils.CaseInsensitiveCompareStrings(string(s), string(t))
}
