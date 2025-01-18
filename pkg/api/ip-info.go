package api

import (
	"context"
	"fmt"
	"go-nordvpn/nordvpnwebapiv1"
	"net/http"
)

func GetIPInfo(webAPIClient *nordvpnwebapiv1.ClientWithResponses) (*nordvpnwebapiv1.IPInfo, error) {
	resp, err := webAPIClient.GetIPInfoWithResponse(context.TODO())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get countries, status code: %d", resp.StatusCode())
	}
	return resp.JSON200, nil
}
