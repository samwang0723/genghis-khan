package honestbee

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

//https://core-staging.honestbee.com/api/countries/PH/available_services?latitude=14.5367633&longitude=121.009545
//https://core.honestbee.com/api/brands?countryCode=TW&page=1&page_size=48&serviceType=groceries&latitude=14.5367633&longitude=121.009545
// /api/stores/17/directory?zoneId=1&latitude=1.3481059&longitude=103.8638325
// /api/zones/2/brands?tag=橄欖油,西式料理&service_type=food

type Service struct {
	ServiceType string `json:"service_type"`
	Avaliable   bool   `json:"available"`
}

func GetServices(countryCode string, latitude float32, longitude float32) (*[]Service, error) {
	client := http.Client{}
	url := fmt.Sprintf("https://core.honestbee.com/api/countries/%s/available_services?latitude=%f&longitude=%f", countryCode, latitude, longitude)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/vnd.honestbee+json;version=1")
	req.Header.Add("Accept-Language", "zh-TW")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var services []Service
	if err := jsoniter.NewDecoder(resp.Body).Decode(services); err != nil {
		return nil, err
	}
	return &services, nil
}
