package discordgo

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"bytes"
	"errors"
	"sync"

	"io/ioutil"
	"net/http/httputil"

	"github.com/Kemonozume/restcl"
	"github.com/gorilla/websocket"
)

type DiscordBot struct {
	Guilds            []Guild
	HeartbeatInterval int
	token             string
	gateway           string
	ct                *time.Ticker
	dialer            websocket.Dialer
	conn              *websocket.Conn
	mut               *sync.Mutex
	isRunning         bool
	fun               HandleMessage
	eventFuncs        map[string]EventFunction
	rest              *restcl.Rest
}

type EventFunction func(*DiscordBot)
type HandleMessage func(MessageResponse, *DiscordBot)

func NewDiscordBot() *DiscordBot {
	d := &DiscordBot{
		dialer: websocket.Dialer{Subprotocols: []string{""}, TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "discord.gg",
		}},
		eventFuncs: make(map[string]EventFunction),
		mut:        &sync.Mutex{},
	}

	rest := restcl.NewRest()
	rest.SetPrefix("https://discordapp.com/api").Use(d.modify)
	rest.Create("/auth/login").SetMethod("POST").Build("login")
	rest.Create("/gateway").SetMethod("GET").Build("gateway")
	rest.Create("/channels/{channelid}/messages").SetMethod("POST").Build("sendmessage")
	rest.Create("/guilds/{guildid}/members/{userid}").SetMethod("PATCH").Build("changerole")
	rest.Create("/channels/{channelid}").SetMethod("PATCH").Build("changechannelinfo")
	rest.Create("/guilds/{guildid}").SetMethod("PATCH").Build("changeserverinfo")
	d.rest = rest
	return d
}

func (d *DiscordBot) modify(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", d.token)
}

func (d *DiscordBot) GetGuildById(id string) (guild Guild, index int) {
	for idx, guild1 := range d.Guilds {
		if guild1.ID == id {
			guild = guild1
			index = idx
			break
		}
	}
	return
}

func (d *DiscordBot) GetChannelById(id string, guild Guild) (channel Channel, index int) {
	for idx, channel1 := range guild.Channels {
		if channel1.ID == id {
			channel = channel1
			index = idx
			break
		}
	}
	return
}

func (d *DiscordBot) updatePresence(msg dPUMessage) {
	guild, index := d.GetGuildById(msg.D.GuildID)
	if guild.ID != "" {
		var index2 int
		var presence Presence
		for idx, pres := range guild.Presences {
			if pres.User.ID == msg.D.User.ID {
				presence = pres
				presence.User = msg.D.User
				presence.GameID = msg.D.GameID
				presence.Status = msg.D.Status
				index2 = idx
				break
			}
		}
		guild.Presences[index2] = presence
		d.Guilds[index] = guild
	}
}

func (d *DiscordBot) updateMemberFromGuild(msg dGMUMessage) {
	guild, index := d.GetGuildById(msg.D.GuildID)
	if guild.ID != "" {
		var index2 int
		var umember Member
		for idx, member := range guild.Members {
			if member.User.ID == msg.D.User.ID {
				umember = member
				umember.Roles = msg.D.Roles
				index2 = idx
				break
			}
		}
		guild.Members[index2] = umember
		d.Guilds[index] = guild
	}
}

func (d *DiscordBot) removeMemberFromGuild(user User, guildid string) {
	guild, index := d.GetGuildById(guildid)
	if guild.ID != "" {
		var members []Member
		for _, member := range guild.Members {
			if !(member.User.ID == user.ID) {
				members = append(members, member)
			}
		}
		guild.Members = members
		d.Guilds[index] = guild
	}
}

func (d *DiscordBot) addMemberToGuild(msg dGMAMessage) {
	guild, index := d.GetGuildById(msg.D.GuildID)
	if guild.ID != "" {
		member := Member{}
		member.User = msg.D.User
		member.JoinedAt = msg.D.JoinedAt
		member.Deaf = false
		member.Roles = msg.D.Roles
		member.Mute = false
		guild.Members = append(guild.Members, member)
		d.Guilds[index] = guild
	}
}

func (d *DiscordBot) updateChannel(msg dCUMessage) {
	guild, gindex := d.GetGuildById(msg.D.GuildID)
	if guild.ID != "" {
		channel, cindex := d.GetChannelById(msg.D.ID, guild)
		if channel.Name != "" {
			channel.Name = msg.D.Name
			channel.Position = msg.D.Position
			channel.Topic = msg.D.Topic
			guild.Channels[cindex] = channel
			d.Guilds[gindex] = guild
		}
	}
}

func (d *DiscordBot) updateGuild(msg dGUMessage) {
	guild, gindex := d.GetGuildById(msg.D.ID)
	if guild.ID != "" {
		guild.Roles = msg.D.Roles
		guild.Region = msg.D.Region
		guild.OwnerID = msg.D.OwnerID
		guild.Name = msg.D.Name
		guild.JoinedAt = msg.D.JoinedAt
		guild.Icon = msg.D.Icon
		guild.AfkTimeout = msg.D.AfkTimeout
		guild.AfkChannelID = msg.D.AfkChannelID
		d.Guilds[gindex] = guild
	}
}

func (d *DiscordBot) GetMemberByName(name string) Member {
	for _, guild := range d.Guilds {
		for _, member := range guild.Members {
			if member.User.Username == name {
				return member
			}
		}
	}
	return Member{}
}

func (d *DiscordBot) SetHandleFunction(f HandleMessage) {
	d.fun = f
}

func (d *DiscordBot) AddCallBack(event string, f EventFunction) {
	//TODO check if event exists and panic if it doesnt
	d.eventFuncs[event] = f
}

func (d *DiscordBot) Login(email string, password string) error {
	login := loginMessage{
		Email:    email,
		Password: password,
	}
	by, err := json.Marshal(login)
	if err != nil {
		return err
	}

	var v map[string]interface{}
	d.rest.Get("login").SetBody(bytes.NewReader(by)).Exec(&v)
	token, exists := v["token"]
	if exists {
		d.token = token.(string)
		return nil
	} else {
		return errors.New("token not found, login information wrong?")
	}
}

func (d *DiscordBot) getGateway() (err error) {
	var v map[string]interface{}
	_, err = d.rest.Get("gateway").Exec(&v)
	if err != nil {
		return
	}
	gate, exists := v["url"]
	if exists {
		d.gateway = gate.(string)
		return nil
	} else {
		return errors.New("gateway not found")
	}
}

func (d *DiscordBot) handshake() (err error) {
	a := handshake{
		Op: 2,
		D: dHD{
			Token: d.token,
			V:     2,
			Properties: dHProperties{
				Os:              "discordgo",
				Browser:         "discordgo",
				Device:          "discordgo",
				Referrer:        "",
				ReferringDomain: "",
			},
		},
	}

	by, err := json.Marshal(a)
	if err != nil {
		return err
	}
	d.conn.WriteMessage(websocket.TextMessage, by)
	return nil
}

func (d *DiscordBot) Start() (ok bool) {
	if d.token == "" {
		panic("Not logged in")
	}

	err := d.getGateway()
	if err != nil {
		panic(err.Error())
	}

	d.conn, _, err = d.dialer.Dial(d.gateway, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer d.conn.Close()

	err = d.handshake()
	if err != nil {
		log.Fatal("handshake: ", err)
	}

	d.isRunning = true

	for {
		//Read Message
		_, message, err := d.conn.ReadMessage()
		if err != nil {
			log.Println("read error: ", err)
			d.stopHeartBeat()
			break
		}

		//transform message to get a look at the T variable that specifies what kind of message we get
		var obj map[string]interface{}
		err = json.Unmarshal(message, &obj)
		if err != nil {
			log.Println("read: ", err)
			d.stopHeartBeat()
			break
		}

		code, ok := obj["t"].(string)
		if !ok {
			log.Println("t doesnt exist")
			log.Println(message)
			d.stopHeartBeat()
			break
		}

		go d.handleMessage(code, message)

		d.mut.Lock()
		if !d.isRunning {
			ok = true
			d.mut.Unlock()
			break
		} else {
			d.mut.Unlock()
		}
	}
	d.stopHeartBeat()
	return
}

func (d *DiscordBot) handleMessage(code string, message []byte) {
	switch code {
	case EVENT_GUILD_MEMBER_REMOVE:
		var GMRemove dGMRMessage
		err := json.Unmarshal(message, &GMRemove)
		checkErr(err)
		d.removeMemberFromGuild(GMRemove.D.User, GMRemove.D.GuildID)
		f, exists := d.eventFuncs[EVENT_GUILD_MEMBER_REMOVE]
		if exists {
			f(d)
		}

	case EVENT_GUILD_MEMBER_ADD:
		var GMAdd dGMAMessage
		err := json.Unmarshal(message, &GMAdd)
		checkErr(err)
		d.addMemberToGuild(GMAdd)
		f, exists := d.eventFuncs[EVENT_GUILD_MEMBER_ADD]
		if exists {
			f(d)
		}

	case EVENT_GUILD_MEMBER_UPDATE:
		var GMUpdate dGMUMessage
		err := json.Unmarshal(message, &GMUpdate)
		checkErr(err)
		d.updateMemberFromGuild(GMUpdate)
		f, exists := d.eventFuncs[EVENT_GUILD_MEMBER_UPDATE]
		if exists {
			f(d)
		}

	case EVENT_PRESENCE_UPDATE:
		var PUpdate dPUMessage
		err := json.Unmarshal(message, &PUpdate)
		checkErr(err)
		d.updatePresence(PUpdate)
		f, exists := d.eventFuncs[EVENT_PRESENCE_UPDATE]
		if exists {
			f(d)
		}

	case EVENT_MESSAGE_CREATE:
		var MessageCreate MessageResponse
		err := json.Unmarshal(message, &MessageCreate)
		checkErr(err)
		if d.fun != nil {
			d.fun(MessageCreate, d)
		}
		f, exists := d.eventFuncs[EVENT_MESSAGE_CREATE]
		if exists {
			f(d)
		}

	case EVENT_CHANNEL_UPDATE:
		var ChannelUpdate dCUMessage
		err := json.Unmarshal(message, &ChannelUpdate)
		checkErr(err)
		d.updateChannel(ChannelUpdate)
		f, exists := d.eventFuncs[EVENT_CHANNEL_UPDATE]
		if exists {
			f(d)
		}

	case EVENT_GUILD_UPDATE:
		var GuildUpdate dGUMessage
		err := json.Unmarshal(message, &GuildUpdate)
		checkErr(err)
		d.updateGuild(GuildUpdate)
		f, exists := d.eventFuncs[EVENT_GUILD_UPDATE]
		if exists {
			f(d)
		}

	case EVENT_READY:
		var ReadyMessage dReadyMessage
		err := json.Unmarshal(message, &ReadyMessage)
		checkErr(err)
		d.ct = time.NewTicker(time.Duration(ReadyMessage.D.HeartbeatInterval) * time.Millisecond)
		go func() {
			for range d.ct.C {
				a := map[string]interface{}{
					"op": 1,
					"d":  makeTimestamp(),
				}
				by, err := json.Marshal(a)
				if err != nil {
					panic(err.Error())
				}
				d.conn.WriteMessage(websocket.TextMessage, by)
			}
		}()
		for _, v := range ReadyMessage.D.Guilds {
			d.Guilds = append(d.Guilds, v)
		}
		f, exists := d.eventFuncs[EVENT_READY]
		if exists {
			f(d)
		}
	}
}

func (d *DiscordBot) stopHeartBeat() {
	if d.ct != nil {
		d.ct.Stop()
	}
}

func (d *DiscordBot) Stop() {
	d.mut.Lock()
	d.isRunning = false
	d.mut.Unlock()
}

func checkErr(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (d *DiscordBot) SendMessage(message MessageRequest, channelid string) (err error) {
	bmessage, err := json.Marshal(message)
	if err != nil {
		return
	}

	_, err = d.rest.Get("sendmessage").SetParams("channelid", channelid).SetBody(bytes.NewReader(bmessage)).Exec(nil)
	return
}

func (d *DiscordBot) ChangeRolesForUser(user Member, guildid string) (err error) {
	ma := map[string]interface{}{
		"roles": user.Roles,
	}
	bmessage, err := json.Marshal(ma)
	if err != nil {
		return
	}
	_, err = d.rest.Get("changerole").SetParams("guildid", guildid, "userid", user.User.ID).
		SetBody(bytes.NewReader(bmessage)).Exec(nil)

	return
}

func (d *DiscordBot) ChangeChannelInformation(channelupdate ChannelUpdateRequest, channelid string) (err error) {
	bmessage, err := json.Marshal(channelupdate)
	if err != nil {
		return
	}
	_, err = d.rest.Get("changechannelinfo").
		SetParams("channelid", channelid).
		SetBody(bytes.NewReader(bmessage)).
		Exec(nil)
	return
}

func (d *DiscordBot) ChangeServerInformation(serverupdate ServerUpdateRequest, guildid string) (err error) {
	bmessage, err := json.Marshal(serverupdate)
	if err != nil {
		return
	}
	_, err = d.rest.Get("changeserverinfo").
		SetParams("guildid", guildid).
		SetBody(bytes.NewReader(bmessage)).
		Exec(nil)
	return
}

func dumpRequest(req *http.Request, name string) {
	dump1, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		panic(err.Error())
	}
	ioutil.WriteFile(name, dump1, 0777)
}

func dumpResponse(resp *http.Response, name string) {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err.Error())
	}
	ioutil.WriteFile(name, dump, 0777)
}
