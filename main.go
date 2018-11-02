package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	FACEBOOK_API = "https://graph.facebook.com/v2.6/me/messages?access_token=%s"
	IMAGE        = "http://37.media.tumblr.com/e705e901302b5925ffb2bcf3cacb5bcd/tumblr_n6vxziSQD11slv6upo3_500.gif"
)

type Callback struct {
	Object string `json:"object,omitempty"`
	Entry  []struct {
		ID        string      `json:"id,omitempty"`
		Time      int         `json:"time,omitempty"`
		Messaging []Messaging `json:"messaging,omitempty"`
	} `json:"entry,omitempty"`
}

type Messaging struct {
	Sender    User    `json:"sender,omitempty"`
	Recipient User    `json:"recipient,omitempty"`
	Timestamp int     `json:"timestamp,omitempty"`
	Message   Message `json:"message,omitempty"`
}

type User struct {
	ID string `json:"id,omitempty"`
}

// https://developers.facebook.com/docs/messenger-platform/send-messages/quick-replies#locations
// "recipient":{
// 	"id":"<PSID>"
// },
// "message":{
// 	"text": "Here is a quick reply!",
// 	"quick_replies":[
// 		{
// 			"content_type":"text",
// 			"title":"Search",
// 			"payload":"<POSTBACK_PAYLOAD>",
// 			"image_url":"http://example.com/img/red.png"
// 		},
// 		{
// 			"content_type":"location"
// 		}
// 	]
// }

type QuickReply struct {
	ContentType string `json:"content_type,omitempty"`
	Payload     string `json:"payload,omitempty"`
}

type Message struct {
	MID          string        `json:"mid,omitempty"`
	Text         string        `json:"text,omitempty"`
	QuickReplies *[]QuickReply `json:"quick_replies,omitempty"`
	Attachments  *[]Attachment `json:"attachments,omitempty"`
	Attachment   *Attachment   `json:"attachment,omitempty"`
}

type Attachment struct {
	Type    string  `json:"type,omitempty"`
	Payload Payload `json:"payload,omitempty"`
}

type Response struct {
	Recipient User    `json:"recipient,omitempty"`
	Message   Message `json:"message,omitempty"`
}

type Coordinates struct {
	Lat  float32 `json:"lat,omitempty"`
	Long float32 `json:"long,omitempty"`
}

type Payload struct {
	URL         string       `json:"url,omitempty"`
	Coordinates *Coordinates `json:"coordinates,omitempty"`
}

func VerificationEndpoint(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("hub.challenge")
	token := r.URL.Query().Get("hub.verify_token")

	if token == os.Getenv("VERIFY_TOKEN") {
		w.WriteHeader(200)
		w.Write([]byte(challenge))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Error, wrong validation token"))
	}
}

func ProcessMessage(event Messaging) {
	if event.Message.Attachment != nil {
		coordinates := event.Message.Attachment.Payload.Coordinates
		log.Printf("User's location %f, %f", coordinates.Lat, coordinates.Long)
	}

	client := &http.Client{}
	var replies []QuickReply
	replies = append(replies, QuickReply{
		ContentType: "location",
	})
	response := Response{
		Recipient: User{
			ID: event.Sender.ID,
		},
		Message: Message{
			Text:         "Please tell me your location",
			QuickReplies: &replies,

			// Attachment: &Attachment{
			// 	Type: "image",
			// 	Payload: Payload{
			// 		URL: IMAGE,
			// 	},
			// },
		},
	}
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&response)
	url := fmt.Sprintf(FACEBOOK_API, os.Getenv("PAGE_ACCESS_TOKEN"))
	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func MessagesEndpoint(w http.ResponseWriter, r *http.Request) {
	var callback Callback
	json.NewDecoder(r.Body).Decode(&callback)
	if callback.Object == "page" {
		for _, entry := range callback.Entry {
			for _, event := range entry.Messaging {
				ProcessMessage(event)
			}
		}
		w.WriteHeader(200)
		w.Write([]byte("Got your message"))
	} else {
		w.WriteHeader(404)
		w.Write([]byte("Message not supported"))
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := mux.NewRouter()
	r.HandleFunc("/webhook", VerificationEndpoint).Methods("GET")
	r.HandleFunc("/webhook", MessagesEndpoint).Methods("POST")
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
