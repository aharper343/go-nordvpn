package api

import (
	"context"
	"fmt"
	"go-nordvpn/nordvpnapiv1"
	"go-nordvpn/pkg/utils"
	"log/slog"
	"net/http"
)

const GroupEnvVarName = "GROUP"

type GroupArray []nordvpnapiv1.Group
type id2GroupMap map[int32]nordvpnapiv1.Group

func (groupArray GroupArray) FilterByTitleOrIdentifier(titleOrIdentifier string) GroupArray {
	var filteredGroups GroupArray
	for _, group := range groupArray {
		if utils.CaseInsensitiveCompareStrings(group.Title, titleOrIdentifier) || utils.CaseInsensitiveCompareStrings(group.Identifier, titleOrIdentifier) {
			filteredGroups = append(filteredGroups, group)
		}
	}
	return filteredGroups
}

func (groupArray GroupArray) FilterById(id int32) GroupArray {
	var filteredGroups GroupArray
	for _, group := range groupArray {
		if group.Id == id {
			filteredGroups = append(filteredGroups, group)
		}
	}
	return filteredGroups
}

func addGroupsToMap(matchedMap *id2GroupMap, groupArray GroupArray) int {
	count := 0
	for _, group := range groupArray {
		(*matchedMap)[group.Id] = group
		count++
	}
	return count
}

func GetGroupEnvVarValue() *utils.StringOrInt32Array {
	lookup, exists := utils.GetMultiEnvVar(GroupEnvVarName)
	if exists {
		return lookup
	}
	return nil
}

func GetGroups(c *nordvpnapiv1.ClientWithResponses, lookup utils.StringOrInt32Array) (*GroupArray, error) {
	resp, err := c.GetGroupsWithResponse(context.TODO(), &nordvpnapiv1.GetGroupsParams{Limit: &maxLimit})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get groups, status code: %d", resp.StatusCode())
	}
	groups := GroupArray(*resp.JSON200)
	slog.Info("Groups retrieved", "count", len(groups))
	matchedMap := make(id2GroupMap)
	stringArray := lookup.ToStringArray()
	failed := 0
	for _, value := range stringArray {
		if addGroupsToMap(&matchedMap, groups.FilterByTitleOrIdentifier(value)) == 0 {
			slog.Warn("No groups found", "title of identifier", value)
			failed++
		}
	}
	int32Array := lookup.ToInt32Array()
	for _, value := range int32Array {
		if addGroupsToMap(&matchedMap, groups.FilterById(value)) == 0 {
			slog.Warn("No groups found", "id", value)
			failed++
		}
	}
	if len(matchedMap) == 0 {
		return nil, fmt.Errorf("no groups found")
	}
	if failed > 0 {
		return nil, fmt.Errorf("%d groups not found", failed)
	}
	var matchedArray GroupArray
	for _, group := range matchedMap {
		matchedArray = append(matchedArray, group)
	}
	return &matchedArray, nil
}

func GetGroupsFromEnvVar(c *nordvpnapiv1.ClientWithResponses) (*GroupArray, error) {
	lookup := GetGroupEnvVarValue()
	if lookup != nil {
		return GetGroups(c, *lookup)
	}
	return nil, nil
}
