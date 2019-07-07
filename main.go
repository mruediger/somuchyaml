package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app       = kingpin.New("somuchyaml", "A mattermost bot.").DefaultEnvars()
	server    = app.Flag("server", "server url").Required().String()
	websocket = app.Flag("websocket", "the websocket url used for listening").Required().String()
	username  = app.Flag("username", "username to connect to the server").Required().String()
	password  = app.Flag("password", "password to connect to the server").Required().String()

	goldenfile = app.Flag("goldenfile", "Enable golden file output.").Bool()

	lastMatch = time.Unix(0, 0)
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	client := model.NewAPIv4Client(*server)

	if _, resp := client.Login(*username, *password); resp.Error != nil {
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

	var handler func(*model.Post, *model.Client4)

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
					handler(post, client)
				}
			}
		}
	}()

	// block forever
	select {}
}

func matchMessage(message string) bool {
	yamlRe := regexp.MustCompile("^(.*[[:space:]]+|)((?i)yaml)([[:space:]]+.*|)$")
	return yamlRe.MatchString(message)
}

func matchFrequency(lastMatch time.Time, now time.Time) bool {
	return now.Sub(lastMatch) < 5*time.Second
}

func sendYaml(originalPost *model.Post, client *model.Client4) {
	if matchMessage(originalPost.Message) {
		if matchFrequency(lastMatch, time.Now()) {
			filename := "yaml.jpg"

			file, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}

			fileinfo, resp := client.UploadFile(file, originalPost.ChannelId, filename)

			if resp.Error != nil {
				log.Println(resp.Error)
				os.Exit(1)
			}

			post := &model.Post{}
			post.ChannelId = originalPost.ChannelId
			post.Filenames = []string{"filename"}
			post.FileIds = []string{fileinfo.FileInfos[0].Id}

			if _, resp := client.CreatePost(post); resp.Error != nil {
				log.Println(err)
				os.Exit(1)
			}
		}
		lastMatch = time.Now()
	}

}

func responseOutput(post *model.Post, client *model.Client4) {
	fmt.Println(post.ToUnsanitizedJson())
}
