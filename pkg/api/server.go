package api

import (
	"fmt"
	"go-nordvpn/nordvpnapiv1"
	"go-nordvpn/pkg/utils"
	"log/slog"
)

const ServerCityEnvName = "CITY"

type ServerGroupArray []nordvpnapiv1.ServerGroup

type ServerLocationArray []nordvpnapiv1.ServerLocation

type ServerCityArray []nordvpnapiv1.ServerCity

func (l ServerLocationArray) GetCitiesByCityName(name string) ServerCityArray {
	var filtered ServerCityArray
	found := map[int32]bool{}
	for _, serverLocation := range l {
		city := serverLocation.Country.City
		if utils.PrefixedCaseInsensitiveCompareStrings(city.Name, name) {
			_, present := found[city.Id]
			if !present {
				found[city.Id] = true
				filtered = append(filtered, city)
			}
		}
	}
	return filtered
}

func (l ServerLocationArray) GetCitiesByCityId(id int32) ServerCityArray {
	var filtered ServerCityArray
	found := map[int32]bool{}
	for _, serverLocation := range l {
		city := serverLocation.Country.City
		if city.Id == id {
			_, present := found[city.Id]
			if !present {
				found[city.Id] = true
				filtered = append(filtered, city)
			}
		}
	}
	return filtered
}

func (l ServerLocationArray) GetCities(lookup utils.StringOrInt32Array) ServerCityArray {
	var filtered ServerCityArray
	found := map[int32]bool{}
	for _, value := range lookup.ToStringArray() {
		for _, city := range l.GetCitiesByCityName(value) {
			_, present := found[city.Id]
			if !present {
				filtered = append(filtered, city)
			}
		}
	}
	for _, value := range lookup.ToInt32Array() {
		for _, city := range l.GetCitiesByCityId(value) {
			_, present := found[city.Id]
			if !present {
				filtered = append(filtered, city)
			}
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	return filtered
}

func (l ServerLocationArray) GetCitiesFromEnvVar() ServerCityArray {
	cities := GetCityEnvVarValue()
	if cities != nil {
		return l.GetCities(*cities)
	}
	return nil
}

func (l ServerLocationArray) CountByCityName(name string) int {
	var count = 0
	for _, serverLocation := range l {
		if utils.PrefixedCaseInsensitiveCompareStrings(serverLocation.Country.City.Name, name) {
			count++
		}
	}
	return count
}

func (l ServerLocationArray) CountByCityId(id int32) int {
	var count = 0
	for _, serverLocation := range l {
		if serverLocation.Country.City.Id == id {
			count++
		}
	}
	return count
}

type ServerArray []nordvpnapiv1.Server

func (s ServerArray) FilterByCityName(name string) ServerArray {
	var filteredServers ServerArray
	for _, server := range s {
		v := ServerLocationArray(server.Locations)
		if v.CountByCityName(name) > 0 {
			filteredServers = append(filteredServers, server)
		}
	}
	return filteredServers
}

func (s ServerArray) FilterByCityId(id int32) ServerArray {
	var filteredServers ServerArray
	for _, server := range s {
		v := ServerLocationArray(server.Locations)
		if v.CountByCityId(id) > 0 {
			filteredServers = append(filteredServers, server)
		}
	}
	return filteredServers
}

type id2ServerMap map[int32]nordvpnapiv1.Server

func addServersToMap(matchedMap *id2ServerMap, serverArray ServerArray) int {
	count := 0
	for _, server := range serverArray {
		(*matchedMap)[server.Id] = server
		count++
	}
	return count
}

func (s ServerArray) FilterByCity(lookup utils.StringOrInt32Array) (*ServerArray, error) {
	matchedMap := make(id2ServerMap)
	stringArray := lookup.ToStringArray()
	failed := 0
	for _, value := range stringArray {
		if addServersToMap(&matchedMap, s.FilterByCityName(value)) == 0 {
			slog.Warn("No servers found", "Location[].Country.City.Name", value)
			failed++
		}
	}
	int32Array := lookup.ToInt32Array()
	for _, value := range int32Array {
		if addServersToMap(&matchedMap, s.FilterByCityId(value)) == 0 {
			slog.Warn("No servers found", "Location[].Country.City.Id", value)
		}
	}
	if len(matchedMap) == 0 {
		return nil, fmt.Errorf("no servers found")
	}
	if failed > 0 {

		slog.Warn(fmt.Sprintf("%d city not found", failed))
	}
	var matchedArray ServerArray
	for _, server := range matchedMap {
		matchedArray = append(matchedArray, server)
	}
	return &matchedArray, nil
}

func (s ServerArray) FilterByCityFromEnvVar() (*ServerArray, error) {
	cities := GetCityEnvVarValue()
	if cities != nil {
		return s.FilterByCity(*cities)
	}
	return nil, nil
}

func GetCityEnvVarValue() *utils.StringOrInt32Array {
	lookup, exists := utils.GetMultiEnvVar(ServerCityEnvName)
	if exists {
		return lookup
	}
	return nil
}

type ServersByLoad ServerArray

func (a ServersByLoad) Len() int           { return len(a) }
func (a ServersByLoad) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ServersByLoad) Less(i, j int) bool { return a[i].Load < a[j].Load }
