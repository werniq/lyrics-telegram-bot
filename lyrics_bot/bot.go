package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
)

var (
	botApiKey = os.Getenv("BOT_API")
)

var (
	errorLog = log.New(os.Stdout, "ERROR\t", log.Lshortfile|log.Ltime|log.Lmsgprefix)
)

type Response struct {
	Hits []Hits `json:"hits"`
}

type Result struct {
	Type         string `json:"_type"`
	Artist       string `json:"artist_names"`
	Title        string `json:"full_title"`
	ThumbnailImg string `json:"header_image_thumbnail_url"`
	Url          string `json:"url"`
}

type Hits struct {
	Result Result `json:"result"`
}

func main() {
	bot, err := tgbotapi.NewBotAPI("6078677160:AAErExpO5bGeZya2axfJ3yNzmBbpGcIM9Cc")
	if err != nil {
		errorLog.Println(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s:", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		//update.Message.Text
		url := fmt.Sprintf("https://genius-song-lyrics1.p.rapidapi.com/search/?q=%s&per_page=10&page=1", update.Message.Text)
		fmt.Println("===================================")
		fmt.Println(update.Message.Text)
		fmt.Println("===================================")
		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("X-RapidAPI-Key", "0e16774897msh0d7f6465454053ap11a57ejsn3d7990cdb478")
		req.Header.Add("X-RapidAPI-Host", "genius-song-lyrics1.p.rapidapi.com")

		resp, _ := http.DefaultClient.Do(req)

		defer resp.Body.Close()

		var response Response
		json.NewDecoder(resp.Body).Decode(&response)
		//json.Unmarshal(resp.Body, &response)

		msg.Text = fmt.Sprintf(`
			Song Title:  %s
			Artist:     %s
			Title: 	    %s
			Image:      %s
			Genius URL: %s			
		`,
			response.Hits[0].Result.Title,
			response.Hits[0].Result.Artist,
			response.Hits[0].Result.ThumbnailImg,
			response.Hits[0].Result.Url)

		if _, err := bot.Send(msg); err != nil {
			errorLog.Println(err)
		}
	}
}
