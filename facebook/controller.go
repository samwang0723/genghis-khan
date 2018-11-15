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

	"github.com/samwang0723/genghis-khan/honestbee"
)

const (
	facebookAPI = "https://graph.facebook.com/v2.6/me/messages?access_token=%s"
)

func postbackHandling(event *Messaging) *Response {
	location := event.location()
	data := strings.Split(event.PostBack.Payload, ":")
	switch data[0] {
	case honestbee.BRANDS:
		brands, err := honestbee.GetBrands("TW", data[2], data[1], location)
		if err != nil {
			str := fmt.Sprintf("No brand served in your location: %s", err.Error())
			return event.showText(str)
		}
		return event.listBrands(*brands)
	case honestbee.DEPARTMENTS:
		departments, err := honestbee.GetDepartments(data[1], location)
		if err != nil {
			str := fmt.Sprintf("No departments found: %s", err.Error())
			return event.showText(str)
		}
		return event.listDepartments(*departments)
	case honestbee.PRODUCTS:
		products, err := honestbee.GetProducts(data[1])
		if err != nil {
			str := fmt.Sprintf("No products found: %s", err.Error())
			return event.showText(str)
		}
		return event.listProducts(*products)
	case honestbee.SEARCH:
		err := event.saveViewingStoreID(data[1])
		if err != nil {
			str := fmt.Sprintf("Viewing storeID store error: %s", err.Error())
			return event.showText(str)
		}
		str := fmt.Sprintf("We've selected store %s, please type search keywords", data[1])
		return event.showText(str)
	}
	return nil
}

func keywordFilters(event *Messaging) *Response {
	if event.PostBack != nil {
		return postbackHandling(event)
	} else if event.AccountLinking != nil {
		accessToken, err := event.saveAccessToken()
		if err != nil {
			return event.showText(err.Error())
		}
		return event.showText(accessToken)
	}

	switch event.Message.Text {
	case "get_location":
		return event.askLocation()
	case "login":
		return event.login()
	}

	viewingStoreID := event.viewingStoreID()
	if len(viewingStoreID) != 0 {
		products, err := honestbee.SearchProducts(viewingStoreID, event.Message.Text)
		if err != nil {
			str := fmt.Sprintf("No products found: %s", err.Error())
			return event.showText(str)
		}
		return event.listProducts(*products)
	}

	coordinates := event.parseLocation()
	if coordinates != nil {
		location := &honestbee.Location{
			Latitude:  coordinates.Lat,
			Longitude: coordinates.Long,
		}
		services, err := honestbee.GetServices("TW", location)
		if err != nil {
			str := fmt.Sprintf("Cannot read services: %s", err.Error())
			return event.showText(str)
		}
		err = event.saveLocation(location)
		if err != nil {
			str := fmt.Sprintf("Location store error: %s", err.Error())
			return event.showText(str)
		}
		return event.listServices(services)
	}

	return nil
}

func respond(body *bytes.Buffer) {
	client := &http.Client{}
	url := fmt.Sprintf(facebookAPI, os.Getenv("PAGE_ACCESS_TOKEN"))
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
	err = readFacebookError(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
}

// Process - handling incoming messages
func (event *Messaging) Process() {
	// sender typing animation
	typing := event.senderTypingAction()
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
