package facebook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/samwang0723/genghis-khan/honestbee"
	"github.com/samwang0723/genghis-khan/utils"
)

const (
	FACEBOOK_API = "https://graph.facebook.com/v2.6/me/messages?access_token=%s"
)

func getCurrentViewingStoreID(SenderID string) string {
	key := fmt.Sprintf("viewing_store_%s", SenderID)
	val, err := utils.RedisClient().Get(key).Result()
	if err == nil {
		return val
	}
	return ""
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

func postbackHandling(event Messaging) *Response {
	location := getLocation(event.Sender.ID)
	data := strings.Split(event.PostBack.Payload, ":")
	switch data[0] {
	case honestbee.BRANDS:
		brands, err := honestbee.GetBrands("TW", data[2], data[1], location)
		if err != nil {
			str := fmt.Sprintf("No brand served in your location: %s", err.Error())
			return ShowText(event.Sender.ID, str)
		}
		return ListBrands(event, *brands)
	case honestbee.DEPARTMENTS:
		departments, err := honestbee.GetDepartments(data[1], location)
		if err != nil {
			str := fmt.Sprintf("No departments found: %s", err.Error())
			return ShowText(event.Sender.ID, str)
		}
		return ListDepartments(event.Sender.ID, *departments)
	case honestbee.PRODUCTS:
		products, err := honestbee.GetProducts(data[1])
		if err != nil {
			str := fmt.Sprintf("No products found: %s", err.Error())
			return ShowText(event.Sender.ID, str)
		}
		return ListProducts(event.Sender.ID, *products)
	case honestbee.SEARCH:
		key := fmt.Sprintf("viewing_store_%s", event.Sender.ID)
		err := utils.RedisClient().Set(key, data[1], 0).Err()
		if err != nil {
			str := fmt.Sprintf("Viewing storeID store error: %s", err.Error())
			return ShowText(event.Sender.ID, str)
		}
		str := fmt.Sprintf("We've selected store %s, please type search keywords", data[1])
		return ShowText(event.Sender.ID, str)
	}
	return nil
}

func keywordFilters(event Messaging) *Response {
	if event.PostBack != nil {
		return postbackHandling(event)
	} else if event.AccountLinking != nil {
		key := fmt.Sprintf("login_%s", event.Sender.ID)
		err := utils.RedisClient().Set(key, event.AccountLinking.AuthorizationCode, 0).Err()
		if err != nil {
			str := fmt.Sprintf("Session store error: %s", err.Error())
			return ShowText(event.Sender.ID, str)
		}
		return ShowText(event.Sender.ID, event.AccountLinking.AuthorizationCode)
	}

	switch event.Message.Text {
	case "get_location":
		return AskLocation(event)
	case "login":
		return Login(event.Sender.ID)
	}

	viewingStoreID := getCurrentViewingStoreID(event.Sender.ID)
	if len(viewingStoreID) != 0 {
		products, err := honestbee.SearchProducts(viewingStoreID, event.Message.Text)
		if err != nil {
			str := fmt.Sprintf("No products found: %s", err.Error())
			return ShowText(event.Sender.ID, str)
		}
		return ListProducts(event.Sender.ID, *products)
	}

	coordinates := ParseLocation(event)
	if coordinates != nil {
		location := &honestbee.Location{
			Latitude:  coordinates.Lat,
			Longitude: coordinates.Long,
		}
		services, err := honestbee.GetServices("TW", location)
		if err != nil {
			str := fmt.Sprintf("Cannot read services: %s", err.Error())
			return ShowText(event.Sender.ID, str)
		}
		key := fmt.Sprintf("location_%s", event.Sender.ID)
		json := fmt.Sprintf(`{"latitude":%f,"longitude":%f}`, coordinates.Lat, coordinates.Long)
		err = utils.RedisClient().Set(key, json, 0).Err()
		if err != nil {
			str := fmt.Sprintf("Location store error: %s", err.Error())
			return ShowText(event.Sender.ID, str)
		}
		return ListServices(event.Sender.ID, services)
	}

	return nil
}

func respond(body *bytes.Buffer) {
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
	err = CheckFacebookError(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
}

func ProcessMessage(event Messaging) {
	// sender typing animation
	typing := SenderTypingAction(event)
	bodyBuffer := new(bytes.Buffer)
	json.NewEncoder(bodyBuffer).Encode(&typing)
	respond(bodyBuffer)

	time.Sleep(2 * time.Second)

	// process keyword parsing and analysis to respond
	response := keywordFilters(event)
	if response == nil {
		return
	}
	bodyBuffer.Truncate(bodyBuffer.Len())
	json.NewEncoder(bodyBuffer).Encode(&response)
	respond(bodyBuffer)
}
