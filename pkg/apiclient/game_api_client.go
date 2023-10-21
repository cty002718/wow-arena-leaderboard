package apiclient

import (
	"encoding/json"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type Season struct {
	ID int `json:"id"`
}

type Server struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
}

type Character struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Server Server `json:"realm"`
}

type Faction struct {
	Type string `json:"type"`
}

type Stat struct {
	Won  int `json:"won"`
	Lost int `json:"lost"`
}

type Entry struct {
	Character Character `json:"character"`
	Faction   Faction   `json:"faction"`
	Rank      int       `json:"rank"`
	Rating    int       `json:"rating"`
	Stat      Stat      `json:"season_match_statistics"`
}

type LeaderboardResponse struct {
	Bracket string  `json:"name"`
	Season  Season  `json:"season"`
	Entries []Entry `json:"entries"`
}

type SeasonResponse struct {
	ID        int   `json:"id"`
	StartedAt int64 `json:"season_start_timestamp"`
	EndedAt   int64 `json:"season_end_timestamp"`
}

type TokenPriceResponse struct {
	LastUpdatedTimestamp int64 `json:"last_updated_timestamp"`
	Price                int   `json:"price"`
}

type ServerType struct {
	Type string `json:"type"`
}

type ServerInfoResponse struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Slug       string     `json:"slug"`
	ServerType ServerType `json:"type"`
}

type IGameApiClient interface {
	GetSeason(seasonId int) (*SeasonResponse, error)
	GetLeaderboard(seasonId int, bracket string) (*LeaderboardResponse, error)
	GetTokenPrice() (*TokenPriceResponse, error)
	GetServerInfo(realmSlug string) (*ServerInfoResponse, error)
}

type GameApiClient struct {
	Client      *resty.Client
	ApiHost     string
	AccessToken string
}

func NewGameApiClient(accessToken string) IGameApiClient {
	return &GameApiClient{
		Client:      resty.New(),
		ApiHost:     os.Getenv("API_HOST"),
		AccessToken: accessToken,
	}
}

func (g *GameApiClient) GetSeason(seasonId int) (*SeasonResponse, error) {
	seasonUrl := g.ApiHost + "/data/wow/pvp-region/0/pvp-season/{pvpSeasonId}"

	resp, err := g.Client.R().
		SetAuthToken(g.AccessToken).
		SetPathParams(map[string]string{
			"pvpSeasonId": cast.ToString(seasonId),
		}).
		SetQueryParams(map[string]string{
			"namespace": "dynamic-classic-tw",
			"locale":    "zh_TW",
		}).
		Get(seasonUrl)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to call season endpoint")
	}

	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("Season endpoint request failed with status code %d", resp.StatusCode())
	}

	response := SeasonResponse{}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (g *GameApiClient) GetLeaderboard(seasonId int, bracket string) (*LeaderboardResponse, error) {
	leaderboardUrl := g.ApiHost + "/data/wow/pvp-region/0/pvp-season/{pvpSeasonId}/pvp-leaderboard/{pvpBracket}"

	resp, err := g.Client.R().
		SetAuthToken(g.AccessToken).
		SetPathParams(map[string]string{
			"pvpSeasonId": cast.ToString(seasonId),
			"pvpBracket":  bracket,
		}).
		SetQueryParams(map[string]string{
			"namespace": "dynamic-classic-tw",
			"locale":    "zh_TW",
		}).
		Get(leaderboardUrl)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to call leaderboard endpoint")
	}

	if resp.StatusCode() != 200 {
		return nil, errors.Errorf("Leaderboard endpoint request failed with status code %d", resp.StatusCode())
	}

	response := LeaderboardResponse{}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (g *GameApiClient) GetTokenPrice() (*TokenPriceResponse, error) {
	tokenPriceUrl := g.ApiHost + "/data/wow/token/"

	resp, err := g.Client.R().
		SetAuthToken(g.AccessToken).
		SetQueryParams(map[string]string{
			"namespace": "dynamic-classic-tw",
			"locale":    "zh_TW",
		}).
		Get(tokenPriceUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to call token price endpoint")
	}

	response := TokenPriceResponse{}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (g *GameApiClient) GetServerInfo(realmSlug string) (*ServerInfoResponse, error) {
	serverInfoUrl := g.ApiHost + "/data/wow/realm/{realmSlug}"

	resp, err := g.Client.R().
		SetAuthToken(g.AccessToken).
		SetPathParams(map[string]string{
			"realmSlug": realmSlug,
		}).
		SetQueryParams(map[string]string{
			"namespace": "dynamic-classic-tw",
			"locale":    "zh_TW",
		}).
		Get(serverInfoUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to call server info endpoint")
	}

	response := ServerInfoResponse{}
	err = json.Unmarshal(resp.Body(), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
