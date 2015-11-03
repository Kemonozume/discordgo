package discordgo

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	EVENT_GUILD_MEMBER_REMOVE = "GUILD_MEMBER_REMOVE"
	EVENT_GUILD_MEMBER_ADD    = "GUILD_MEMBER_ADD"
	EVENT_GUILD_MEMBER_UPDATE = "GUILD_MEMBER_UPDATE"
	EVENT_PRESENCE_UPDATE     = "PRESENCE_UPDATE"
	EVENT_MESSAGE_CREATE      = "MESSAGE_CREATE"
	EVENT_READY               = "READY"
	EVENT_CHANNEL_UPDATE      = "CHANNEL_UPDATE"
	EVENT_GUILD_UPDATE        = "GUILD_UPDATE"
)

//GATEWAY RESPONSE STRUCT
type gatewayResponse struct {
	URL string `json:"url"`
}

//HANDSHAKE REQUEST STRUCT
type handshake struct {
	Op int `json:"op"`
	D  dHD `json:"d"`
}

type dHD struct {
	Token      string       `json:"token"`
	Properties dHProperties `json:"properties"`
	V          int          `json:"v"`
}

type dHProperties struct {
	Os              string `json:"$os"`
	Browser         string `json:"$browser"`
	Device          string `json:"$device"`
	Referrer        string `json:"$referrer"`
	ReferringDomain string `json:"$referring_domain"`
}

//READYMESSAGE STRUCTS
type dReadyMessage struct {
	T  string `json:"t"`
	S  int    `json:"s"`
	Op int    `json:"op"`
	D  struct {
		V                 int              `json:"v"`
		User              dROurUser        `json:"user"`
		SessionID         string           `json:"session_id"`
		ReadState         []readState      `json:"read_state"`
		PrivateChannels   []privateChannel `json:"private_channels"`
		HeartbeatInterval int              `json:"heartbeat_interval"`
		Guilds            []Guild          `json:"guilds"`
	} `json:"d"`
}

//BOT USER
type dROurUser struct {
	Verified      bool   `json:"verified"`
	Username      string `json:"username"`
	ID            string `json:"id"`
	Email         string `json:"email"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

type readState struct {
	MentionCount  int    `json:"mention_count"`
	LastMessageID string `json:"last_message_id"`
	ID            string `json:"id"`
}

type privateChannel struct {
	Recipient     recipient `json:"recipient"`
	LastMessageID string    `json:"last_message_id"`
	IsPrivate     bool      `json:"is_private"`
	ID            string    `json:"id"`
}

type recipient struct {
	Username      string `json:"username"`
	ID            string `json:"id"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

//Guild struct (contains members, member Status and channels
type Guild struct {
	VoiceStates  []interface{} `json:"voice_states"`
	Roles        []Role        `json:"roles"`
	Region       string        `json:"region"`
	Presences    []Presence    `json:"presences"`
	OwnerID      string        `json:"owner_id"`
	Name         string        `json:"name"`
	Members      []Member      `json:"members"`
	JoinedAt     time.Time     `json:"joined_at"`
	ID           string        `json:"id"`
	Icon         string        `json:"icon"`
	Channels     []Channel     `json:"channels"`
	AfkTimeout   int           `json:"afk_timeout"`
	AfkChannelID interface{}   `json:"afk_channel_id"`
}

type Role struct {
	Position    int    `json:"position"`
	Permissions int    `json:"permissions"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Hoist       bool   `json:"hoist"`
	Color       int    `json:"color"`
}

type Presence struct {
	User   User        `json:"user"`
	Status string      `json:"status"`
	GameID interface{} `json:"game_id"`
}

type User struct {
	Username      string      `json:"username"`
	ID            string      `json:"id"`
	Discriminator json.Number `json:"discriminator,Number"`
	Avatar        string      `json:"avatar"`
}

func (d User) Mention() string {
	return fmt.Sprintf("<@%v>", d.ID)
}

type Member struct {
	User     User      `json:"user"`
	Roles    []string  `json:"roles"`
	Mute     bool      `json:"mute"`
	JoinedAt time.Time `json:"joined_at"`
	Deaf     bool      `json:"deaf"`
}

type Channel struct {
	Type                 string                   `json:"type"`
	Topic                string                   `json:"topic"`
	Position             int                      `json:"position"`
	PermissionOverwrites []dRPermissionOverwrites `json:"permission_overwrites"`
	Name                 string                   `json:"name"`
	LastMessageID        string                   `json:"last_message_id"`
	ID                   string                   `json:"id"`
}

type dRPermissionOverwrites struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Deny  int    `json:"deny"`
	Allow int    `json:"allow"`
}

//MESSAGE_CREATE
type MessageResponse struct {
	Op int    `json:"op"`
	S  int    `json:"s"`
	T  string `json:"t"`
	D  struct {
		Attachments     []interface{} `json:"attachments"`
		Author          User          `json:"author"`
		ChannelID       string        `json:"channel_id"`
		Content         string        `json:"content"`
		EditedTimestamp interface{}   `json:"edited_timestamp"`
		Embeds          []interface{} `json:"embeds"`
		ID              string        `json:"id"`
		MentionEveryone bool          `json:"mention_everyone"`
		Mentions        []User        `json:"mentions"`
		Nonce           string        `json:"nonce"`
		Timestamp       string        `json:"timestamp"`
		Tts             bool          `json:"tts"`
	} `json:"d"`
}

//Message_Send
type MessageRequest struct {
	Content  string   `json:"content"`
	Mentions []string `json:"mentions"`
	Tts      bool     `json:"tts"`
}

func (d *MessageRequest) AddMention(user User) {
	d.Mentions = append(d.Mentions, user.ID)
}

func NewMessage(content string) MessageRequest {
	return MessageRequest{Content: content}
}

//Login Message
type loginMessage struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Guild Member Remove message
type dGMRMessage struct {
	T  string `json:"t"`
	S  int    `json:"s"`
	Op int    `json:"op"`
	D  struct {
		User    User   `json:"user"`
		GuildID string `json:"guild_id"`
	} `json:"d"`
}

//Guild Member Added message
type dGMAMessage struct {
	T  string `json:"t"`
	S  int    `json:"s"`
	Op int    `json:"op"`
	D  struct {
		User     User      `json:"user"`
		Roles    []string  `json:"roles"`
		JoinedAt time.Time `json:"joined_at"`
		GuildID  string    `json:"guild_id"`
	} `json:"d"`
}

//Guild Member Update message
type dGMUMessage struct {
	T  string `json:"t"`
	S  int    `json:"s"`
	Op int    `json:"op"`
	D  struct {
		User    User     `json:"user"`
		Roles   []string `json:"roles"`
		GuildID string   `json:"guild_id"`
	} `json:"d"`
}

//Presence Update Message
type dPUMessage struct {
	T  string `json:"t"`
	S  int    `json:"s"`
	Op int    `json:"op"`
	D  struct {
		User    User        `json:"user"`
		Status  string      `json:"status"`
		Roles   []string    `json:"roles"`
		GuildID string      `json:"guild_id"`
		GameID  interface{} `json:"game_id"`
	} `json:"d"`
}

//Channel Update Message
type dCUMessage struct {
	T  string `json:"t"`
	S  int    `json:"s"`
	Op int    `json:"op"`
	D  struct {
		Type                 string                   `json:"type"`
		Topic                string                   `json:"topic"`
		Position             int                      `json:"position"`
		PermissionOverwrites []dRPermissionOverwrites `json:"permission_overwrites"`
		Name                 string                   `json:"name"`
		LastMessageID        string                   `json:"last_message_id"`
		IsPrivate            bool                     `json:"is_private"`
		ID                   string                   `json:"id"`
		GuildID              string                   `json:"guild_id"`
	} `json:"d"`
}

//Guild Update Message
type dGUMessage struct {
	T  string `json:"t"`
	S  int    `json:"s"`
	Op int    `json:"op"`
	D  struct {
		Roles          []Role      `json:"roles"`
		Region         string      `json:"region"`
		OwnerID        string      `json:"owner_id"`
		Name           string      `json:"name"`
		JoinedAt       time.Time   `json:"joined_at"`
		ID             string      `json:"id"`
		Icon           string      `json:"icon"`
		EmbedEnabled   bool        `json:"embed_enabled"`
		EmbedChannelID interface{} `json:"embed_channel_id"`
		AfkTimeout     int         `json:"afk_timeout"`
		AfkChannelID   interface{} `json:"afk_channel_id"`
	} `json:"d"`
}

//Channel Update Request
type ChannelUpdateRequest struct {
	Name     string `json:"name"`
	Position int    `json:"position"`
	Topic    string `json:"topic"`
}

//Server Update Request
type ServerUpdateRequest struct {
	Name         string      `json:"name"`
	Icon         string      `json:"icon"`
	Region       string      `json:"region"`
	AfkChannelID interface{} `json:"afk_channel_id"`
	AfkTimeout   int         `json:"afk_timeout"`
}
