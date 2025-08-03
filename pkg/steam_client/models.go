package steamclient

type ApiResponse[T any] struct {
	Response T `json:"response"`
}

type PlayerList struct {
	Players []Player `json:"players"`
}

type Player struct {
	SteamID                  string `json:"steamid"`
	CommunityVisibilityState int    `json:"communityvisibilitystate"`
	ProfileState             int    `json:"profilestate"`
	PersonaName              string `json:"personaname"`
	CommentPermission        int    `json:"commentpermission"`
	ProfileURL               string `json:"profileurl"`
	Avatar                   string `json:"avatar"`
	AvatarMedium             string `json:"avatarmedium"`
	AvatarFull               string `json:"avatarfull"`
	AvatarHash               string `json:"avatarhash"`
	LastLogoff               int64  `json:"lastlogoff"`
	PersonaState             int    `json:"personastate"`
	RealName                 string `json:"realname"`
	PrimaryClanID            string `json:"primaryclanid"`
	TimeCreated              int64  `json:"timecreated"`
	PersonaStateFlags        int    `json:"personastateflags"`
	LocCountryCode           string `json:"loccountrycode"`
	LocStateCode             string `json:"locstatecode"`
	LocCityID                int    `json:"loccityid"`
}

type RecentlyPlayedGames struct {
	TotalCount int                  `json:"total_count"`
	Games      []RecentlyPlayedGame `json:"games"`
}

type RecentlyPlayedGame struct {
	AppID            int    `json:"appid"`
	Name             string `json:"name"`
	PlaytimeTwoWeeks int    `json:"playtime_2weeks"`
	PlaytimeForever  int    `json:"playtime_forever"`
	ImgIconUrl       string `json:"img_icon_url"`
}

type PlayerGameStats struct {
	PlayerStats GameStats `json:"playerstats"`
}

type GameStats struct {
	SteamID      string                  `json:"steamID"`
	Achievements []GameStatsAchievements `json:"achievements"`
}

type GameStatsAchievements struct {
	Name     string `json:"name"`
	Achieved int    `json:"achieved"`
}

type GameDetailsResponse map[string]GameDetails

type GameDetails struct {
	Success bool            `json:"success"`
	Data    GameDetailsData `json:"data"`
}

type GameDetailsData struct {
	Type             string `json:"type"`
	Name             string `json:"name"`
	AppID            int    `json:"steam_appid"`
	ShortDescription string `json:"short_description"`
	HeaderImage      string `json:"header_image"`
}

type AllGameAchievementsReponse struct {
	Game AllGameAchievements `json:"game"`
}

type AllGameAchievements struct {
	AppID             int               `json:"appid"`
	AvaiableGameStats AvaiableGameStats `json:"availableGameStats"`
}

type AvaiableGameStats struct {
	Achievements []GameAchievement `json:"achievements"`
}

type GameAchievement struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}
