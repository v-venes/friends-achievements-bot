package repository

import (
	"context"
	"time"

	steamclient "github.com/v-venes/friends-achievements-bot/pkg/steam_client"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	DEFAULT_DATABASE  = "achievements"
	PLAYER_COLLECTION = "players"
)

type PlayerRepository struct {
	MongoClient *mongo.Client
}

type PlayerModel struct {
	PlayerID               string    `bson:"player_id"`
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
		PlayerID:               resp.SteamID,
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
	playerCollection := pr.MongoClient.Database(DEFAULT_DATABASE).Collection(PLAYER_COLLECTION)

	updateDoc, err := bson.Marshal(player)
	if err != nil {
		return err
	}

	var document bson.M
	err = bson.Unmarshal(updateDoc, &document)
	if err != nil {
		return err
	}

	update := bson.M{"$set": document}
	filter := bson.M{"player_id": player.PlayerID}
	opts := options.UpdateOne().SetUpsert(true)
	_, err = playerCollection.UpdateOne(context.TODO(), filter, update, opts)
	return err
}
