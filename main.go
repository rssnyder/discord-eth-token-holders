package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	frequency   int
	setNickname bool
	activity    string
	header      string
	token       string
	address     string
)

func init() {

	flag.IntVar(&frequency, "frequency", 5, "seconds between gas price cycles")
	flag.BoolVar(&setNickname, "setNickname", false, "wether to set nickname of bot")
	flag.StringVar(&activity, "activity", "", "text for activity")
	flag.StringVar(&header, "header", "", "text for nickname")
	flag.StringVar(&token, "token", "", "discord bot token")
	flag.StringVar(&address, "address", "", "address of the token contract")

	flag.Parse()
}

const (
	holdersUrl = "https://eth-token-holders.cloud.rileysnyder.org/%s"
)

func main() {

	// create a new discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// show as online
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening discord connection,", err)
		return
	}

	// set activity as desc
	if setNickname {
		err = dg.UpdateGameStatus(0, activity)
		if err != nil {
			fmt.Printf("Unable to set activity: \n", err)
		} else {
			fmt.Println("Set activity")
		}
	}

	// get guides for bot
	guilds, err := dg.UserGuilds(100, "", "")
	if err != nil {
		fmt.Println("Error getting guilds: ", err)
		setNickname = false
	}

	ticker := time.NewTicker(time.Duration(frequency) * time.Second)
	var nickname string

	for {

		select {
		case <-ticker.C:

			data := getHolders(address)

			nickname = fmt.Sprintf("%s%s", header, data)

			if setNickname {

				for _, g := range guilds {

					err = dg.GuildMemberNickname(g.ID, "@me", nickname)
					if err != nil {
						fmt.Printf("Error updating nickname: %s\n", err)
						continue
					} else {
						fmt.Printf("Set nickname in %s: %s\n", g.Name, nickname)
					}
				}
			} else {

				err = dg.UpdateGameStatus(0, nickname)
				if err != nil {
					fmt.Printf("Unable to set activity: %s\n", err)
				} else {
					fmt.Printf("Set activity: %s\n", nickname)
				}
			}
		}
	}
}

func getHolders(contract string) string {
	var holders string

	reqURL := fmt.Sprintf(holdersUrl, contract)
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return holders
	}

	req.Header.Add("User-Agent", "Mozilla/5.0")
	req.Header.Add("accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return holders
	}

	results, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return holders
	}

	holders = string(results)

	return holders
}
