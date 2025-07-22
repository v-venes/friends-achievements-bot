package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/v-venes/friends-achievements-bot/pkg/models"
)

const STEAM_BASE_URL = "http://api.steampowered.com"

type SteamClient struct {
	steamKey   string
	httpClient *http.Client
}

type NewSteamClientParams struct {
	SteamKey string
}

func NewSteamClient(params NewSteamClientParams) *SteamClient {
	httpClient := &http.Client{}
	return &SteamClient{
		steamKey:   params.SteamKey,
		httpClient: httpClient,
	}
}

func (s *SteamClient) GetPlayerSummary(steamid string) (*models.Player, error) {
	resp, err := s.httpClient.Get(fmt.Sprintf("%s/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", STEAM_BASE_URL, s.steamKey, steamid))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	playerSummaryResponse := &models.PlayerSummaryResponse{}

	err = json.NewDecoder(resp.Body).Decode(&playerSummaryResponse)
	if err != nil {
		return nil, err
	}

	if len(playerSummaryResponse.Response.Players) == 0 {
		return nil, errors.New("No SteamID found!")
	}

	return &playerSummaryResponse.Response.Players[0], nil
}
