package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app       = kingpin.New("somuchyaml", "A mattermost bot.")
	server    = app.Flag("server", "server url").Required().String()
	websocket = app.Flag("websocket", "the websocket url used for listening").Required().String()
	username  = app.Flag("username", "username to connect to the server").Required().String()
	password  = app.Flag("password", "password to connect to the server").Required().String()

	team    = app.Flag("team", "the team on the server").Default("darksystem").String()
	channel = app.Flag("channel", "the channel the bot will listen").Default("town-square").String()

	goldenfile = app.Flag("goldenfile", "Enable golden file output.").Bool()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	client := model.NewAPIv4Client(*server)

	if _, resp := client.Login(*username, *password); resp.Error != nil {
		log.Println(resp.Error)
		os.Exit(1)
	}

	team, resp := client.GetTeamByName(*team, "")
	if resp.Error != nil {
		log.Println(resp.Error)
		os.Exit(1)
	}

	channel, resp := client.GetChannelByName(*channel, team.Id, "")
	if resp.Error != nil {
		log.Println(resp.Error)
		os.Exit(1)
	}

	ws, err := model.NewWebSocketClient4(*websocket, client.AuthToken)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if ws != nil {
				ws.Close()
			}
			os.Exit(0)
		}
	}()

	ws.Listen()

	var handler func(*model.Post, *model.Client4, *model.Channel)

	if *goldenfile {
		handler = responseOutput
	} else {
		handler = sendYaml
	}

	go func() {
		for {
			select {
			case event := <-ws.EventChannel:
				if s, ok := event.Data["post"].(string); ok {
					post := model.PostFromJson(strings.NewReader(s))
					handler(post, client, channel)
				}
			}
		}
	}()

	// block forever
	select {}
}

func matchMessage(message string) bool {
	yamlRe := regexp.MustCompile("^(.*[[:^alpha:]]+|)y|Ya|Am|Al|L([[:^alpha:]]+.*|)$")
	return yamlRe.MatchString(message)
}

func sendYaml(post *model.Post, client *model.Client4, channel *model.Channel) {
	if matchMessage(post.Message) {
		filename := "yaml.jpg"

		file, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		fileinfo, resp := client.UploadFile(file, channel.Id, filename)

		if resp.Error != nil {
			log.Println(resp.Error)
			os.Exit(1)
		}

		post := &model.Post{}
		post.ChannelId = channel.Id
		post.Filenames = []string{"filename"}
		post.FileIds = []string{fileinfo.FileInfos[0].Id}

		if _, resp := client.CreatePost(post); resp.Error != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
}

func responseOutput(post *model.Post, client *model.Client4, channel *model.Channel) {
	fmt.Println(post.ToUnsanitizedJson())
}
