package main

import (
	"context"
	"fmt"
	"go-nordvpn/nordvpnapiv1"
	"go-nordvpn/pkg/api"
	"go-nordvpn/pkg/template"
	"go-nordvpn/pkg/utils"
	"log"
	"log/slog"
	"math"
	"math/rand"
	"net/http"
	"sort"
)

var maxLimit nordvpnapiv1.Limit = math.MaxInt16
var statusOnline nordvpnapiv1.ServerStatus = nordvpnapiv1.ServerStatusOnline

func logServer(num int, server nordvpnapiv1.Server) {
	slog.Info("\tServer", "#", num, "id", server.Id, "hostname", server.Hostname, "load", server.Load, "ip", server.Station)
	cities := api.ServerLocationArray(server.Locations).GetCitiesFromEnvVar()
	if cities == nil {
		found := map[int32]bool{}
		for _, location := range api.ServerLocationArray(server.Locations) {
			city := location.Country.City
			_, present := found[city.Id]
			if !present {
				cities = append(cities, city)
			}
		}
	}
	for _, city := range cities {
		slog.Info("\t\tCity", "id", city.Id, "name", city.Name)
	}
}

func main() {

	hc := http.Client{}
	c, err := nordvpnapiv1.NewClientWithResponses("https://api.nordvpn.com", nordvpnapiv1.WithHTTPClient(&hc))
	if err != nil {
		log.Panic("Failed to create client", err)
	}

	serverParams := nordvpnapiv1.GetServersParams{
		Limit:         &maxLimit,
		FiltersStatus: &statusOnline,
	}

	technology, err := api.GetTechnologyFromEnvVar(c)
	if err != nil {
		log.Panic("Failed to get technologies", err)
	}

	if technology == nil {
		log.Panicf("Required environment variable %s or %s was not set", api.TechnologyEnvVarName, api.ProtocolEnvVarName)
	} else {
		slog.Info("Technology", "id", technology.Id, "name", technology.Name, "internalIdentifier", technology.InternalIdentifier)
		serverParams.FiltersServersTechnologiesId = &technology.Id
	}
	protocol, port, err := api.GetProtocolAndPort(*technology)
	if err != nil {
		log.Fatal("Failed to get protocol and port", err)
	}

	country, err := api.GetCountryFromEnvVar(c)
	if err != nil {
		log.Panic("Failed to get countries", err)
	}

	if country == nil {
		slog.Warn(fmt.Sprintf("Optional environment variable %s not set", api.CountryEnvVarName))
	} else {
		slog.Info("Country", "id", country.Id, "name", country.Name, "code", country.Code)
		serverParams.FiltersCountryId = &country.Id
	}

	groups, err := api.GetGroupsFromEnvVar(c)
	if err != nil {
		log.Panic("failed to get groups", err)
	}

	if groups == nil {
		slog.Warn(fmt.Sprintf("Optional environment variable %s not set", api.GroupEnvVarName))
	} else {
		slog.Info("Groups", "count", len(*groups))
		var groupIds []int32
		for _, group := range *groups {
			groupIds = append(groupIds, group.Id)
			slog.Info("Group", "id", group.Id, "title", group.Title, "identifier", group.Identifier)
		}
		if len(groupIds) > 0 {
			serverParams.FiltersServersGroupsId = &groupIds
		}
	}

	resp, err := c.GetServersWithResponse(context.TODO(), &serverParams)

	if err != nil {
		log.Panic("Failed to get servers", err)
	}

	if resp.StatusCode() != http.StatusOK {
		log.Panicf("Failed to get servers, status code: %d", resp.StatusCode())
	}

	servers := api.ServerArray(*resp.JSON200)

	if len(servers) == 0 {
		log.Fatal("No servers matched all the criteria")
	} else {
		slog.Info("Servers fetched", "count", len(servers))
	}

	cityServers, err := servers.FilterByCityFromEnvVar()

	if err != nil {
		slog.Warn("Failed to filter servers by city", "error", err)
	} else {
		if cityServers == nil {
			slog.Warn(fmt.Sprintf("Optional environment variable %s not set", api.ServerCityEnvName))
		} else {
			slog.Info("Servers filtered by city", "count", len(*cityServers))
			servers = *cityServers
		}
	}

	maxSelected := 1
	randomTop, _ := utils.GetSingleEnvVar("RANDOM_TOP")
	if randomTop != nil {
		if randomTop.Type == "int32" {
			if randomTop.Int32Value > 0 {
				maxSelected = int(randomTop.Int32Value)
			} else {
				slog.Error("Zero or negative value", "RANDOM_TOP", randomTop.Int32Value)
				randomTop = nil
			}
		} else {
			slog.Error("No a numeric value", "RANDOM_TOP", randomTop.StringValue)
			randomTop = nil
		}
	}
	if randomTop == nil {
		maxSelected = 10
	}

	maxSelected = min(maxSelected, len(servers))
	if maxSelected == 0 {
		log.Fatal("No servers matched all the criteria")
	}

	sort.Sort(api.ServersByLoad(servers))

	maxDisplay := min(10, len(servers))
	slog.Info(fmt.Sprintf("Top %d of %d selected nordvpn servers:", maxDisplay, maxSelected))
	for i := 0; i < maxDisplay; i++ {
		logServer(i+1, servers[i])
	}

	selected := 0
	if randomTop != nil {
		selected = rand.Intn(maxSelected)
	}
	server := servers[selected]

	slog.Info("Selected server:")
	logServer(selected+1, server)
	slog.Info("\t\tConfig", "protocol", protocol, "port", port)
	err = template.WriteOVPNFile(server.Hostname, server.Station, protocol, port)
	if err != nil {
		log.Panic("Failed to create OpenVPN configuration", err)
	}
}
