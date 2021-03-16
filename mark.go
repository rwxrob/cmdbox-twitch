package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rwxrob/auth-go"
	"github.com/rwxrob/cmdtab"
	"github.com/rwxrob/conf-go"
	"github.com/rwxrob/uniq-go"
)

func init() {
	x := cmdtab.New("mark")
	x.Usage = `[text]`

	x.Description = ` 
		Places a mark in the current active video and the later resulting
		VOD for use when editing or helping others find a location.

		The name of the mark is the current time to the second. 

		If an optional string argument is passed, it will be saved as
		a mnemonic helper for later lookup when using the marks to create
		highlight videos. This is helpful because Twitch does not allow long
		names in marks.`

	x.Method = func(args []string) error {

		description := uniq.IsoSecond()

		if len(args) > 0 {
			description += "\n" + strings.Join(args, " ")
			if len(description) > 140 {
				return errors.New("must be less than 125 characters (140 total)")
			}
		}

		// get the twitchid from config
		config, err := conf.New()
		if err != nil {
			return err
		}
		err = config.Load()
		if err != nil {
			return err
		}
		twitchid := config.Get("twitch.id")
		if twitchid == "" {
			return errors.New("twitch.id not found in configuration")
		}

		// fetch the twitch client and auth token data
		_, app, err := auth.Lookup("twitch")
		if err != nil {
			return err
		}

		// create the data to post
		jsn, err := json.Marshal(map[string]string{
			"user_id":     twitchid,
			"description": description,
		})
		if err != nil {
			return err
		}
		data := bytes.NewBuffer(jsn)

		// setup the post request
		ctx, cancel := context.WithCancel(context.Background())
		req, err := http.NewRequestWithContext(ctx, "POST",
			"https://api.twitch.tv/helix/streams/markers", data)
		if err != nil {
			return err
		}

		// add the headers
		req.Header.Add("Client-ID", app.ClientID)
		req.Header.Add("Authorization", "Bearer "+app.AccessToken)
		req.Header.Add("Content-Type", "application/json")

		// get a client
		client := new(http.Client)

		// start the timer
		go func() {
			time.Sleep(time.Second * 10)
			fmt.Println("Cancelled")
			cancel()
		}()

		// get a response
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		// display effective return status
		fmt.Println(resp.StatusCode)

		return nil
	}
}
