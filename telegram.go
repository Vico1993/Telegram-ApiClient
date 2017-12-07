package telegram

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Client is a structure to use my telegram method
type Client struct {
	ChatID string `json:"chat_id"`
	BotAPI string `json:"bot_api"`
}

type sendStruct struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"result"`
}

type returnStruct struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

// SendTelegramMessage is use to send message from a bot to an chat Room
func (cl *Client) SendTelegramMessage(msg string, notification bool) bool {
	var URL *url.URL
	URL, err := url.Parse("https://api.telegram.org/bot" + cl.BotAPI + "/sendMessage")
	if err != nil {
		fmt.Println("Erreur Parsing URL :", err)
		panic("Please check this error")
	}

	parameters := url.Values{}
	parameters.Add("chat_id", cl.ChatID)
	parameters.Add("parse_mode", "markdown")
	parameters.Add("text", msg)
	if !notification {
		parameters.Add("disable_notification", "true")
	}
	URL.RawQuery = parameters.Encode()

	// Build the request
	req, err := http.NewRequest("GET", URL.String(), nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("erreur ReadAll: ", err)
	}
	var result = string(body)
	var newrecord sendStruct

	json.NewDecoder(strings.NewReader(result)).Decode(&newrecord)
	if newrecord.Ok {
		return true
	}

	return false
}

// GetTelegramCommand BETA...
// GetTelegramCommand is use to get command that you send to your bot
func GetTelegramCommand(body []byte) string {
	data := &returnStruct{}
	jsonErr := json.Unmarshal([]byte(string(body)), data)
	if jsonErr != nil {
		log.Fatal("Error json Unmarshal : ", jsonErr)
	}

	if data.Message.From.IsBot == true {
		log.Fatal("C'est un bot...")
	}

	return data.Message.Text
}
