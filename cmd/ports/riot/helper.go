package riot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func buildNewGetAccountByRiotIdRequest(
	ctx context.Context,
	server string,
	params AccountByRiotIdRequestParams,
) (*http.Request, error) {
	var err error
	server = fmt.Sprintf(server, "americas")

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("riot/account/v1/accounts/by-riot-id/%s/%s", params.Name, params.Tagline)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req.WithContext(ctx), nil
}

func parseAccountByRiotIdResponse(
	response *http.Response,
) (*AccountByRiotIdResponse, error) {
	bodyBytes, err := io.ReadAll(response.Body)
	defer func() { _ = response.Body.Close() }()
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		var dest Status
		err = json.Unmarshal(bodyBytes, &dest)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(fmt.Sprintf("HTTP StatusCode: %d - %s", dest.StatusCode, dest.Message))
	}

	var dest AccountByRiotIdResponse
	err = json.Unmarshal(bodyBytes, &dest)

	if err != nil {
		return nil, err
	}
	return &dest, nil
}

func buildNewGetQueueEntriesByPlayerUuidRequest(
	ctx context.Context,
	server string,
	params QueueEntriesByPlayerUuidParams,
) (*http.Request, error) {
	var err error
	server = fmt.Sprintf(server, "na1")

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("lol/league/v4/entries/by-puuid/%s", params.PlayerUuid)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req.WithContext(ctx), nil
}

func parseQueueEntriesByPlayerUuidResponse(response *http.Response) ([]*QueueResponse, error) {
	bodyBytes, err := io.ReadAll(response.Body)
	defer func() { _ = response.Body.Close() }()
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		var dest Status
		err = json.Unmarshal(bodyBytes, &dest)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(fmt.Sprintf("HTTP StatusCode: %d - %s", dest.StatusCode, dest.Message))
	}

	var dest []*QueueResponse
	err = json.Unmarshal(bodyBytes, &dest)

	if err != nil {
		return nil, err
	}
	return dest, nil
}
