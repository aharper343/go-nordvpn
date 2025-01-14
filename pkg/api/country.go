package api

import (
	"context"
	"fmt"
	"go-nordvpn/nordvpnapiv1"
	"go-nordvpn/pkg/utils"
	"log/slog"
	"net/http"
)

const CountryEnvVarName = "COUNTRY"

type CountryArray []nordvpnapiv1.Country

func (c CountryArray) FilterByNameOrCode(nameOrCode string) CountryArray {
	var filteredCountries CountryArray
	for _, country := range c {
		slog.Debug("Country.FilterByNameOrCode", "id", country.Id, "Name", country.Name, "Code", country.Code, "With", nameOrCode)
		if utils.CaseInsensitiveCompareStrings(country.Name, nameOrCode) || utils.CaseInsensitiveCompareStrings(country.Code, nameOrCode) {
			filteredCountries = append(filteredCountries, country)
		}
	}
	return filteredCountries
}

func (c CountryArray) FilterById(id int32) CountryArray {
	var filteredCountries CountryArray
	for _, country := range c {
		slog.Debug("Country.FilterById", "id", country.Id, "Name", country.Name, "Code", country.Code, "With", id)
		if country.Id == id {
			filteredCountries = append(filteredCountries, country)
		}
	}
	return filteredCountries
}

func GetCountryEnvVarValue() *utils.StringOrInt32 {
	lookup, exists := utils.GetSingleEnvVar(CountryEnvVarName)
	if exists {
		return lookup
	}
	return nil
}

func GetCountry(c *nordvpnapiv1.ClientWithResponses, lookup utils.StringOrInt32) (*nordvpnapiv1.Country, error) {
	resp, err := c.GetCountriesWithResponse(context.TODO(), &nordvpnapiv1.GetCountriesParams{Limit: &maxLimit})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get countries, status code: %d", resp.StatusCode())
	}
	countries := CountryArray(*resp.JSON200)
	slog.Info("Countries retrieved", "count", len(countries))
	var matchedArray CountryArray
	var matchingField, matchingValue string
	if lookup.Type == "int32" {
		matchedArray = countries.FilterById(lookup.Int32Value)
		matchingField = "Id"
		matchingValue = string(lookup.Int32Value)
	} else {
		matchedArray = countries.FilterByNameOrCode(lookup.StringValue)
		matchingField = "Name or Code"
		matchingValue = lookup.StringValue
	}
	switch len(matchedArray) {
	case 0:
		return nil, fmt.Errorf("no country found with %s %s", matchingField, matchingValue)
	case 1:
		return &matchedArray[0], nil
	default:
		return nil, fmt.Errorf("%d countries found with %s %s", len(matchedArray), matchingField, matchingValue)
	}
}

func GetCountryFromEnvVar(c *nordvpnapiv1.ClientWithResponses) (*nordvpnapiv1.Country, error) {
	lookup := GetCountryEnvVarValue()
	if lookup != nil {
		return GetCountry(c, *lookup)
	}
	return nil, nil
}
