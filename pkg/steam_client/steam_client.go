package steamclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	STEAM_BASE_URL       = "http://api.steampowered.com"
	STEAM_STORE_BASE_URL = "https://store.steampowered.com/"
)

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

func (s *SteamClient) GetPlayerSummary(steamid string) (*Player, error) {
	resp, err := s.httpClient.Get(fmt.Sprintf("%s/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s", STEAM_BASE_URL, s.steamKey, steamid))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	playerSummaryResponse := &ApiResponse[PlayerList]{}

	err = json.NewDecoder(resp.Body).Decode(&playerSummaryResponse)
	if err != nil {
		return nil, err
	}

	if len(playerSummaryResponse.Response.Players) == 0 {
		return nil, errors.New("No SteamID found!")
	}

	return &playerSummaryResponse.Response.Players[0], nil
}

func (s *SteamClient) GetRecentlyPlayedGames(steamid string) (*RecentlyPlayedGames, error) {
	resp, err := s.httpClient.Get(fmt.Sprintf("%s/IPlayerService/GetRecentlyPlayedGames/v0001/?key=%s&steamid=%s&format=json", STEAM_BASE_URL, s.steamKey, steamid))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	recentlyPlayedGamesResponse := &ApiResponse[RecentlyPlayedGames]{}

	err = json.NewDecoder(resp.Body).Decode(&recentlyPlayedGamesResponse)
	if err != nil {
		return nil, err
	}

	return &recentlyPlayedGamesResponse.Response, nil
}

func (s *SteamClient) GetGameStats(steamid string, appid int) (*PlayerGameStats, error) {
	resp, err := s.httpClient.Get(fmt.Sprintf("%s/ISteamUserStats/GetUserStatsForGame/v0002/?key=%s&steamid=%s&appid=%d", STEAM_BASE_URL, s.steamKey, steamid, appid))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	playerStatsReponse := &PlayerGameStats{}

	err = json.NewDecoder(resp.Body).Decode(&playerStatsReponse)
	if err != nil {
		return nil, err
	}

	return playerStatsReponse, nil
}

func (s *SteamClient) GetGameDetails(appid int) (*GameDetailsData, error) {
	resp, err := s.httpClient.Get(fmt.Sprintf("%s/api/appdetails?appids=%d&cc=br&l=portuguese", STEAM_STORE_BASE_URL, appid))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var gameDetailsReponse map[int]GameDetails

	err = json.NewDecoder(resp.Body).Decode(&gameDetailsReponse)
	if err != nil {
		return nil, err
	}

	gameDetailsData := gameDetailsReponse[appid].Data

	return &gameDetailsData, nil
}

func (s *SteamClient) GetAllGameAchievements(appid int) (*AllGameAchievements, error) {
	resp, err := s.httpClient.Get(fmt.Sprintf("%s/ISteamUserStats/GetSchemaForGame/v2/?key=%s&appid=%d", STEAM_BASE_URL, s.steamKey, appid))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	gameDetailsReponse := &AllGameAchievementsReponse{}

	err = json.NewDecoder(resp.Body).Decode(&gameDetailsReponse)
	if err != nil {
		return nil, err
	}

	gameDetailsData := gameDetailsReponse.Game
	gameDetailsData.AppID = appid

	return &gameDetailsData, nil
}
