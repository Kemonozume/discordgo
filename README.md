#discordgo

discordgo is a simple, small bot for discordapp.com

## Getting Started

Install discordgo
~~~  go
go get github.com/Kemonozume/discordgo
~~~ 

Start Using it
~~~ go
package main

import (
    "github.com/Kemonozume/discordgo"
    "fmt"
)

func handleMessage(message discordgo.DMessageCreate, bot *discordgo.DiscordBot) {
	fmt.Printf("%20v: %v\n", message.D.Author.Username, message.D.Content)
	for i,v := range message.D.Mentions {
		fmt.Printf("#%v %v(%v)\n", i, v.Username, v.ID)
	}
	if message.D.Content == "ping" {
		bot.SendMessage(discordgo.NewMessage("pong"), channelid)
	}
	if message.D.Content == "say hello" {
		bot.SendMessage(discordgo.DMessageSend{Content: "hello", Tts: true}, channelid)
	}
}

func main() {
 	bot := discordgo.NewDiscordBot()
 	bot.Login("email", "password")
 	bot.SetHandleFunction(handleMessage)
 	go bot.Start() 
 	
 	//wait for connection to be established
 	time.Sleep(4 * time.Second)
 	for i, channel := range bot.Channels {
		fmt.Printf("%v channel: %v(%v)", i, channel.Name, channel.ID)
	}
	for i, member := range bot.Members {
		fmt.Printf("%v member: %v(%v)", i, member.User.Name, member.User.ID)
	}
	time.Sleep(3 * time.Minute)
}
~~~

### Development
The Discord API is still in development. Functions may break at any time.  
In such an event, please contact me or submit a pull request.

The API is also available in these languages :
* [Java](https://github.com/nerd/Discord4J)
* [.NET](https://github.com/RogueException/Discord.Net)
* [C#](https://github.com/Luigifan/DiscordSharp)
* [Node.js](https://github.com/discord-js/discord.js) / [Alternative](https://github.com/izy521/discord.io)
* [Python](https://github.com/Rapptz/discord.py)
* [Ruby](https://github.com/meew0/discordrb)

### Pull requests
No one is perfect at programming and I am no exception. If you see something that can be improved, please feel free to submit a pull request! 

### Overview
- [x] Reading Messages
- [x] Sending Messages (tts, mentions)
- [ ] Documentation
- [ ] Edits 
- [ ] Typing notifications

and probably some more, i guess i added around 20% of the Unofficial Discord API
