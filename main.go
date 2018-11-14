package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/samwang0723/genghis-khan/utils"

	"github.com/gorilla/mux"
	"github.com/samwang0723/genghis-khan/facebook"
	"github.com/samwang0723/genghis-khan/honestbee"
)

const (
	FACEBOOK_API = "https://graph.facebook.com/v2.6/me/messages?access_token=%s"
)

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

func getLocation(SenderID string) *honestbee.Location {
	key := fmt.Sprintf("location_%s", SenderID)
	val, err := utils.RedisClient().Get(key).Result()
	if err == nil {
		location := new(honestbee.Location)
		jsoniter.UnmarshalFromString(val, &location)
		return location
	}
	return nil
}

func postbackHandling(event facebook.Messaging) *facebook.Response {
	location := getLocation(event.Sender.ID)
	data := strings.Split(event.PostBack.Payload, ":")
	switch data[0] {
	case honestbee.BRANDS:
		brands, err := honestbee.GetBrands("TW", data[2], data[1], location)
		if err != nil {
			str := fmt.Sprintf("No brand served in your location: %s", err.Error())
			return facebook.ComposeText(event.Sender.ID, str)
		}
		return facebook.ComposeBrandList(event, *brands)
	case honestbee.DEPARTMENTS:
		departments, err := honestbee.GetDepartments(data[1], location)
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
		key := fmt.Sprintf("login_%s", event.Sender.ID)
		err := utils.RedisClient().Set(key, event.AccountLinking.AuthorizationCode, 0).Err()
		if err != nil {
			str := fmt.Sprintf("Session store error: %s", err.Error())
			return facebook.ComposeText(event.Sender.ID, str)
		}
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
		location := &honestbee.Location{
			Latitude:  coordinates.Lat,
			Longitude: coordinates.Long,
		}
		services, err := honestbee.GetServices("TW", location)
		if err != nil {
			str := fmt.Sprintf("Cannot read services: %s", err.Error())
			return facebook.ComposeText(event.Sender.ID, str)
		}
		key := fmt.Sprintf("location_%s", event.Sender.ID)
		json := fmt.Sprintf(`{"latitude":%f,"longitude":%f}`, coordinates.Lat, coordinates.Long)
		err = utils.RedisClient().Set(key, json, 0).Err()
		if err != nil {
			str := fmt.Sprintf("Location store error: %s", err.Error())
			return facebook.ComposeText(event.Sender.ID, str)
		}
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

	router := mux.NewRouter()
	router.HandleFunc("/webhook", VerificationEndpoint).Methods("GET")
	router.HandleFunc("/webhook", MessagesEndpoint).Methods("POST")

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		utils.RedisClient().Close() //FIXME: What if still have pending tasks
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
