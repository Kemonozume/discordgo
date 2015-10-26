package discordgo

import (
	"encoding/json"
	"fmt"
	"time"
)

//GATEWAY RESPONSE STRUCT
type dGatewayResponse struct {
	URL string `json:"url"`
}

//HANDSHAKE REQUEST STRUCT
type dHandshake struct {
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
		V                 int                `json:"v"`
		User              dROurUser          `json:"user"`
		SessionID         string             `json:"session_id"`
		ReadState         []dRReadState      `json:"read_state"`
		PrivateChannels   []dRPrivateChannel `json:"private_channels"`
		HeartbeatInterval int                `json:"heartbeat_interval"`
		Guilds            []DRGuild          `json:"guilds"`
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

type dRReadState struct {
	MentionCount  int    `json:"mention_count"`
	LastMessageID string `json:"last_message_id"`
	ID            string `json:"id"`
}

type dRPrivateChannel struct {
	Recipient     dRRecipient `json:"recipient"`
	LastMessageID string      `json:"last_message_id"`
	IsPrivate     bool        `json:"is_private"`
	ID            string      `json:"id"`
}

type dRRecipient struct {
	Username      string `json:"username"`
	ID            string `json:"id"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

type DRGuild struct {
	VoiceStates  []interface{} `json:"voice_states"`
	Roles        []DRRole      `json:"roles"`
	Region       string        `json:"region"`
	Presences    []dRPresence  `json:"presences"`
	OwnerID      string        `json:"owner_id"`
	Name         string        `json:"name"`
	Members      []DRMember    `json:"members"`
	JoinedAt     time.Time     `json:"joined_at"`
	ID           string        `json:"id"`
	Icon         string        `json:"icon"`
	Channels     []DRChannel   `json:"channels"`
	AfkTimeout   int           `json:"afk_timeout"`
	AfkChannelID interface{}   `json:"afk_channel_id"`
}

type DRRole struct {
	Position    int    `json:"position"`
	Permissions int    `json:"permissions"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	Hoist       bool   `json:"hoist"`
	Color       int    `json:"color"`
}

type dRPresence struct {
	User   DRUser      `json:"user"`
	Status string      `json:"status"`
	GameID interface{} `json:"game_id"`
}

type DRUser struct {
	Username      string      `json:"username"`
	ID            string      `json:"id"`
	Discriminator json.Number `json:"discriminator,Number"`
	Avatar        string      `json:"avatar"`
}

func (d DRUser) Mention() string {
	return fmt.Sprintf("<@%v>", d.ID)
}

type DRMember struct {
	User     DRUser    `json:"user"`
	Roles    []string  `json:"roles"`
	Mute     bool      `json:"mute"`
	JoinedAt time.Time `json:"joined_at"`
	Deaf     bool      `json:"deaf"`
}

type DRChannel struct {
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
type DMessageCreate struct {
	Op int    `json:"op"`
	S  int    `json:"s"`
	T  string `json:"t"`
	D  struct {
		Attachments     []interface{} `json:"attachments"`
		Author          DRUser        `json:"author"`
		ChannelID       string        `json:"channel_id"`
		Content         string        `json:"content"`
		EditedTimestamp interface{}   `json:"edited_timestamp"`
		Embeds          []interface{} `json:"embeds"`
		ID              string        `json:"id"`
		MentionEveryone bool          `json:"mention_everyone"`
		Mentions        []DRUser      `json:"mentions"`
		Nonce           string        `json:"nonce"`
		Timestamp       string        `json:"timestamp"`
		Tts             bool          `json:"tts"`
	} `json:"d"`
}

//Message SEnd
type DMessageSend struct {
	Content  string   `json:"content"`
	Mentions []string `json:"mentions"`
	Tts      bool     `json:"tts"`
}

func (d *DMessageSend) AddMention(user DRUser) {
	d.Mentions = append(d.Mentions, user.ID)
}

func NewMessage(content string) DMessageSend {
	return DMessageSend{Content: content}
}

//Login Message
type dLoginMessage struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Guild Member Remove message
type dGMRMessage struct {
	T  string `json:"t"`
	S  int    `json:"s"`
	Op int    `json:"op"`
	D  struct {
		User    DRUser `json:"user"`
		GuildID string `json:"guild_id"`
	} `json:"d"`
}

//Guild Member Added message
type dGMAMessage struct {
	T  string `json:"t"`
	S  int    `json:"s"`
	Op int    `json:"op"`
	D  struct {
		User     DRUser    `json:"user"`
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
		User    DRUser   `json:"user"`
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
		User    DRUser      `json:"user"`
		Status  string      `json:"status"`
		Roles   []string    `json:"roles"`
		GuildID string      `json:"guild_id"`
		GameID  interface{} `json:"game_id"`
	} `json:"d"`
}
