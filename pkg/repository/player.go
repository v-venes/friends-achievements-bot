package repository

import (
	"context"
	"time"

	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	PLAYER_COLLECTION = "players"
	DEFAULT_DATABASE  = "achievements"
)

type PlayerRepository struct {
	MongoClient *mongo.Client
}

type PlayerModel struct {
	SteamID                string    `bson:"steamid"`
	ProfileVisibilityState int       `bson:"profile_visibility_state"`
	Name                   string    `bson:"name"`
	ProfileURL             string    `bson:"profile_url"`
	Avatar                 string    `bson:"avatar"`
	LastLogoff             time.Time `bson:"last_logoff,omitempty"`
	CurrentStatus          int       `bson:"current_status,omitempty"`
	RealName               string    `bson:"real_name,omitempty"`
	AccountCreatedAt       time.Time `bson:"account_created_at"`
	AccountCountryCode     string    `bson:"account_coutry_code,omitempty"`
	CreatedAt              time.Time `bson:"created_at"`
	UpdatedAt              time.Time `bson:"updated_at"`
}

func NewPlayerFromSteam(resp *steamclient.Player) *PlayerModel {
	return &PlayerModel{
		SteamID:                resp.SteamID,
		ProfileVisibilityState: resp.CommunityVisibilityState,
		Name:                   resp.PersonaName,
		ProfileURL:             resp.ProfileURL,
		Avatar:                 resp.AvatarFull,
		LastLogoff:             time.Unix(resp.LastLogoff, 0),
		CurrentStatus:          resp.PersonaState,
		RealName:               resp.RealName,
		AccountCreatedAt:       time.Unix(resp.TimeCreated, 0),
		AccountCountryCode:     resp.LocCountryCode,
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
}

func (pr *PlayerRepository) CreatePlayer(player PlayerModel) error {
	_, err := pr.MongoClient.Database(DEFAULT_DATABASE).Collection(PLAYER_COLLECTION).InsertOne(context.TODO(), player)
	if err != nil {
		return err
	}

	return nil
}
