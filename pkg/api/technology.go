package api

import (
	"context"
	"fmt"
	"go-nordvpn/nordvpnapiv1"
	"go-nordvpn/pkg/utils"
	"log/slog"
	"net/http"
	"strings"
)

const TechnologyEnvVarName = "TECHNOLOGY"
const ProtocolEnvVarName = "PROTOCOL"

type TechnologyArray []nordvpnapiv1.Technology

func (t TechnologyArray) FilterByNameOrInternalIdentifier(nameOrInternalIdentifier string) TechnologyArray {
	var filteredTechnologies TechnologyArray
	for _, technology := range t {
		slog.Debug("Technology.FilterByNameOrInternalIdentifier", "id", technology.Id, "Name", technology.Name, "InternalIdentifier", technology.InternalIdentifier, "With", nameOrInternalIdentifier)
		if utils.CaseInsensitiveCompareStrings(technology.Name, nameOrInternalIdentifier) || utils.CaseInsensitiveCompareStrings(technology.InternalIdentifier, nameOrInternalIdentifier) {
			filteredTechnologies = append(filteredTechnologies, technology)
		}
	}
	return filteredTechnologies
}

func (t TechnologyArray) FilterById(id int32) TechnologyArray {
	var filteredTechnologies TechnologyArray
	for _, technology := range t {
		slog.Debug("Technology.FilterById", "id", technology.Id, "Name", technology.Name, "InternalIdentifier", technology.InternalIdentifier, "With", id)
		if technology.Id == id {
			filteredTechnologies = append(filteredTechnologies, technology)
		}
	}
	return filteredTechnologies
}

func GetTechnologyEnvVarValue() *utils.StringOrInt32 {
	lookup, exists := utils.GetSingleEnvVar(TechnologyEnvVarName)
	if !exists {
		lookup, exists = utils.GetSingleEnvVar(ProtocolEnvVarName)
	}
	if exists {
		if lookup.Type == "string" {
			lookup.StringValue = strings.Replace(lookup.StringValue, "_", "-", -1)
		}
		return lookup
	}
	return nil
}

func GetTechnology(c *nordvpnapiv1.ClientWithResponses, lookup utils.StringOrInt32) (*nordvpnapiv1.Technology, error) {
	resp, err := c.GetTechnologiesWithResponse(context.TODO(), &nordvpnapiv1.GetTechnologiesParams{Limit: &maxLimit})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get technologies, status code: %d", resp.StatusCode())
	}
	technologies := TechnologyArray(*resp.JSON200)
	slog.Info("Technologies retrieved", "count", len(technologies))
	var matchedArray TechnologyArray
	var matchingField, matchingValue string
	if lookup.Type == "int32" {
		matchedArray = technologies.FilterById(lookup.Int32Value)
		matchingField = "Id"
		matchingValue = string(lookup.Int32Value)
	} else {
		matchedArray = technologies.FilterByNameOrInternalIdentifier(lookup.StringValue)
		matchingField = "Name or InternalIdentifier"
		matchingValue = lookup.StringValue
	}
	switch len(matchedArray) {
	case 0:
		return nil, fmt.Errorf("no technology found with %s=%s", matchingField, matchingValue)
	case 1:
		return &matchedArray[0], nil
	default:
		return nil, fmt.Errorf("%d technologies found with %s=%s", len(matchedArray), matchingField, matchingValue)
	}
}

func GetTechnologyFromEnvVar(c *nordvpnapiv1.ClientWithResponses) (*nordvpnapiv1.Technology, error) {
	lookup := GetTechnologyEnvVarValue()
	if lookup != nil {
		return GetTechnology(c, *lookup)
	}
	return nil, nil
}

func GetProtocolAndPort(technology nordvpnapiv1.Technology) (string, int, error) {
	if utils.CaseInsensitiveCompareStrings(technology.InternalIdentifier, "openvpn-udp") {
		return "udp", 1194, nil
	}
	if utils.CaseInsensitiveCompareStrings(technology.InternalIdentifier, "openvpn-tcp") {
		return "tcp", 443, nil
	}
	return "", -1, fmt.Errorf("unsupported technology of %s", technology.InternalIdentifier)
}
