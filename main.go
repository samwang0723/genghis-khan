package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/samwang0723/genghis-khan/facebook"
	"github.com/samwang0723/genghis-khan/honestbee"
)

const (
	FACEBOOK_API = "https://graph.facebook.com/v2.6/me/messages?access_token=%s"
	IMAGE        = "http://37.media.tumblr.com/e705e901302b5925ffb2bcf3cacb5bcd/tumblr_n6vxziSQD11slv6upo3_500.gif"
)

var latitude float32
var longitude float32
var currentQueryStoreID string

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

func postbackHandling(event facebook.Messaging) *facebook.Response {
	data := strings.Split(event.PostBack.Payload, ":")
	switch data[0] {
	case honestbee.BRANDS:
		brands, err := honestbee.GetBrands("TW", data[2], data[1], latitude, longitude)
		if err != nil {
			str := fmt.Sprintf("No brand served in your location: %s", err.Error())
			return facebook.ComposeText(event.Sender.ID, str)
		}
		return facebook.ComposeBrandList(event, *brands)
	case honestbee.DEPARTMENTS:
		departments, err := honestbee.GetDepartments(data[1], latitude, longitude)
		if err != nil {
			str := fmt.Sprintf("No departments found: %s", err.Error())
			return facebook.ComposeText(event.Sender.ID, str)
		}
		return facebook.ComposeDepartmentList(event.Sender.ID, *departments)
	case honestbee.PRODUCTS:
		products, err := honestbee.GetProducts(data[1])
		if err != nil {
			str := fmt.Sprintf("No products found: %s", err.Error())
			return facebook.ComposeText(event.Sender.ID, str)
		}
		return facebook.ComposeProductList(event.Sender.ID, *products)
	case honestbee.SEARCH:
		currentQueryStoreID = data[1]
		str := fmt.Sprintf("We've selected store %s, please type search keywords", currentQueryStoreID)
		return facebook.ComposeText(event.Sender.ID, str)
	}
	return nil
}

func keywordFilters(event facebook.Messaging) *facebook.Response {
	if event.PostBack != nil {
		return postbackHandling(event)
	} else if event.AccountLinking != nil {
		return facebook.ComposeText(event.Sender.ID, event.AccountLinking.AuthorizationCode)
	}

	switch event.Message.Text {
	case "get_location":
		return facebook.ComposeLocation(event)
	case "login":
		return facebook.Login(event.Sender.ID)
	}

	if len(currentQueryStoreID) != 0 {
		products, err := honestbee.SearchProducts(currentQueryStoreID, event.Message.Text)
		if err != nil {
			str := fmt.Sprintf("No products found: %s", err.Error())
			return facebook.ComposeText(event.Sender.ID, str)
		}
		return facebook.ComposeProductList(event.Sender.ID, *products)
	}

	coordinates := facebook.ParseLocation(event)
	if coordinates != nil {
		services, err := honestbee.GetServices("TW", coordinates.Lat, coordinates.Long)
		if err != nil {
			str := fmt.Sprintf("Cannot read services: %s", err.Error())
			return facebook.ComposeText(event.Sender.ID, str)
		}
		latitude = coordinates.Lat
		longitude = coordinates.Long
		return facebook.ComposeServicesButton(event.Sender.ID, services)
	}

	return nil
}

func Respond(body *bytes.Buffer) {
	client := &http.Client{}
	url := fmt.Sprintf(FACEBOOK_API, os.Getenv("PAGE_ACCESS_TOKEN"))
	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Println(err.Error())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	err = facebook.CheckFacebookError(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
}

func ProcessMessage(event facebook.Messaging) {
	typing := facebook.SenderTypingAction(event)
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&typing)
	Respond(body)

	time.Sleep(2 * time.Second)
	response := keywordFilters(event)
	if response == nil {
		return
	}
	body.Truncate(body.Len())
	json.NewEncoder(body).Encode(&response)
	Respond(body)
}

func MessagesEndpoint(w http.ResponseWriter, r *http.Request) {
	var callback facebook.Callback
	json.NewDecoder(r.Body).Decode(&callback)
	log.Printf("%+v\n", callback)
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
