package discordgo

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"bytes"
	"errors"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	baseURL    = "https://discordapp.com/api"
	loginURL   = baseURL + "/auth/login"
	gatewayURL = baseURL + "/gateway"
)

type HandleMessage func(DMessageCreate, *DiscordBot)

type DiscordBot struct {
	User              dROurUser
	Members           []DRMember
	Roles             []DRRole
	Channels          []DRChannel
	Guilds            []DRGuild
	HeartbeatInterval int
	token             string
	gateway           string
	ct                *time.Ticker
	dialer            websocket.Dialer
	conn              *websocket.Conn
	mut               *sync.Mutex
	isRunning         bool
	fun               HandleMessage
}

func NewDiscordBot() *DiscordBot {
	return &DiscordBot{
		dialer: websocket.Dialer{Subprotocols: []string{""}, TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         "discord.gg",
		}},
		mut: &sync.Mutex{},
	}
}

func (d *DiscordBot) GetUserByName(name string) DRMember {
	for _, v := range d.Members {
		if v.User.Username == name {
			return v
		}
	}
	return DRMember{}
}

func (d *DiscordBot) SetHandleFunction(f HandleMessage) {
	d.fun = f
}

func (d *DiscordBot) Login(email string, password string) error {
	login := dLoginMessage{
		Email:    email,
		Password: password,
	}
	by, err := json.Marshal(login)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", loginURL, bytes.NewReader(by))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	var v map[string]interface{}
	err = transformToJson(resp, &v)
	if err != nil {
		return err
	}
	token, exists := v["token"]
	if exists {
		d.token = token.(string)
		return nil
	} else {
		return errors.New("token not found, login information wrong?")
	}
}

func (d *DiscordBot) getGateway() (err error) {
	req, err := http.NewRequest("GET", gatewayURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("origin", "https://discordapp.com")
	req.Header.Add("authorization", d.token)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	var v map[string]interface{}
	err = transformToJson(resp, &v)
	if err != nil {
		return err
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
	a := dHandshake{
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

		switch code {
		case "MESSAGE_CREATE":
			var MessageCreate DMessageCreate
			err := json.Unmarshal(message, &MessageCreate)
			checkErr(err)
			if d.fun != nil {
				d.fun(MessageCreate, d)
			}
		case "READY":
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

			d.User = ReadyMessage.D.User
			for _, v := range ReadyMessage.D.Guilds {
				d.Guilds = append(d.Guilds, v)
				d.Members = append(d.Members, v.Members...)
			}
		}

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

var client *http.Client = &http.Client{
	Timeout: time.Duration(30 * time.Second),
	Transport: &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		DisableKeepAlives: true,
	},
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (d *DiscordBot) SendMessage(message DMessageSend, channelid string) (by []byte, err error) {
	bmessage, err := json.Marshal(message)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", "https://discordapp.com/api/channels/"+channelid+"/messages", bytes.NewReader(bmessage))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", d.token)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	by, err = ioutil.ReadAll(resp.Body)
	return
}

func (d *DiscordBot) ChangeRolesForUser(user DRMember, guildid string) (by []byte, err error) {
	ma := map[string]interface{}{
		"roles": user.Roles,
	}
	bmessage, err := json.Marshal(ma)
	if err != nil {
		return
	}
	req, err := http.NewRequest("PATCH", "https://discordapp.com/api/guilds/"+guildid+"/members/"+user.User.ID, bytes.NewReader(bmessage))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("authorization", d.token)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	by, err = ioutil.ReadAll(resp.Body)
	return
}

func transformToJson(resp *http.Response, c interface{}) (err error) {
	defer resp.Body.Close()
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(by, c)
	return
}
