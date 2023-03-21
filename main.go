package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type EmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url"`
}

// EmbedAuthor for Embed Author Structure
type EmbedAuthor struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

// EmbedField for Embed Field Structure
type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

// Embed is for Embed Structure
type Embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	URL         string       `json:"url"`
	Color       int          `json:"color"`
	Content     []string     `json:"content"`
	Footer      EmbedFooter  `json:"footer"`
	Author      EmbedAuthor  `json:"author"`
	Fields      []EmbedField `json:"fields"`
	Timestamp   string       `json:"timestamp"`
}

type Payload struct {
	Wait      bool    `json:"wait"`
	Content   string  `json:"content"`
	Username  string  `json:"username"`
	AvatarURL string  `json:"avatar_url"`
	TTS       bool    `json:"tts"`
	Embeds    []Embed `json:"embeds"`
}

func Color() int {
	if os.Getenv("PLUGIN_COLOR") != "" {
		color := strings.Replace(os.Getenv("PLUGIN_COLOR"), "#", "", -1)
		if s, err := strconv.ParseInt(color, 16, 32); err == nil {
			return int(s)
		}
	}

	switch os.Getenv("DRONE_BUILD_STATUS") {
	case "success":
		// green
		return 0x1ac600
	case "failure", "error", "killed":
		// red
		return 0xff3232
	default:
		// yellow
		return 0xffd930
	}
}

const username = "Drone"
const avatarURL = "https://oyster.ignimgs.com/mediawiki/apis.ign.com/starcraft-2/3/36/Drone.jpg"

func main() {
	// var description string
	// switch os.Getenv("DRONE_BUILD_EVENT") {
	// case "push":
	// 	description = fmt.Sprintf("**%s** pushed to `%s`.", os.Getenv("DRONE_COMMIT_AUTHOR"), os.Getenv("DRONE_COMMIT_BRANCH"))
	// case "pull_request":
	// 	description = fmt.Sprintf("**%s** opened a pull request from `%s` to `%s`", os.Getenv("DRONE_COMMIT_AUTHOR"), os.Getenv("DRONE_SOURCE_BRANCH"), os.Getenv("DRONE_COMMIT_BRANCH"))
	// case "tag":
	// 	description = fmt.Sprintf("**%s** created tag `%s`.", os.Getenv("DRONE_COMMIT_AUTHOR"), os.Getenv("DRONE_TAG"))
	// }
	log.Print("Hello!")
	var emoji string
	title := "%s Build %s #**%s**"
	fields := []EmbedField{
		{
			Name:   "Branch",
			Value:  os.Getenv("DRONE_COMMIT_BRANCH"),
			Inline: true,
		},
		{
			Name:   "Status",
			Value:  os.Getenv("DRONE_BUILD_STATUS"),
			Inline: true,
		},
		{
			Name:   "Changes",
			Value:  os.Getenv("DRONE_COMMIT_LINK"),
			Inline: true,
		},
	}
	if os.Getenv("DRONE_BUILD_STATUS") == "failure" {
		emoji = ":red_circle:"
		failed_steps := EmbedField{
			Name:   "Failed steps",
			Value:  os.Getenv("DRONE_FAILED_STEPS"),
			Inline: false,
		}

		fields = append(fields, failed_steps)
	} else {
		emoji = ":green_circle:"
	}
	embed := Embed{
		Title:       fmt.Sprintf(title, emoji, os.Getenv("DRONE_REPO"), os.Getenv("DRONE_BUILD_NUMBER")),
		Description: os.Getenv("DRONE_COMMIT_MESSAGE"),
		Fields:      fields,
		URL:         os.Getenv("DRONE_BUILD_LINK"),
		Color:       Color(),
		Author: EmbedAuthor{
			Name:    os.Getenv("DRONE_COMMIT_AUTHOR"),
			IconURL: os.Getenv("DRONE_COMMIT_AUTHOR_AVATAR"),
		},
	}
	payload := &Payload{
		Username:  username,
		AvatarURL: avatarURL,
		Embeds:    []Embed{embed},
	}

	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(payload)
	req, _ := http.NewRequest("POST", os.Getenv("PLUGIN_WEBHOOK"), payloadBuf)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, e := client.Do(req)
	if e != nil {
		log.Print(e)
		os.Exit(1)
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		log.Print(res.StatusCode)
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Print(bodyString)
		os.Exit(1)
	}
	defer res.Body.Close()
	log.Print("Bye!")
}
