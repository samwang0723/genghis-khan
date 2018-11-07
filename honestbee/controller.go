package honestbee

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type Service struct {
	ServiceType string `json:"service_type"`
	Avaliable   bool   `json:"available"`
}

type Meta struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalCount  int `json:"total_count"`
}

type Brands struct {
	Brands  []Brand `json:"brands"`
	Meta    Meta    `json:"meta"`
	ID      string  `json:"id"`
	Variant string  `json:"variant"`
}

type Brand struct {
	ID                       int    `json:"id"`
	Name                     string `json:"name"`
	About                    string `json:"about"`
	Description              string `json:"description"`
	FreeDeliveryEligible     bool   `json:"freeDeliveryEligible"`
	BrandColor               string `json:"brandColor"`
	DefaultConciergeFee      string `json:"defaultConciergeFee"`
	DefaultDeliveryFee       string `json:"defaultDeliveryFee"`
	MinimumOrderFreeDelivery string `json:"minimumOrderFreeDelivery"`
	MinimumSpendExtraFee     string `json:"minimumSpendExtraFee"`
	ServiceType              string `json:"serviceType"`
	Slug                     string `json:"slug"`
	StoreID                  int    `json:"storeId"`
	ImageURL                 string `json:"imageUrl"`
	Currency                 string `json:"currency"`
	SameStorePrice           bool   `json:"sameStorePrice"`
	ProductsCount            int    `json:"productsCount"`
	Closed                   bool   `json:"closed"`
	OpensAt                  string `json:"opensAt"`
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

	services := make([]Service, 0)
	if err := jsoniter.NewDecoder(resp.Body).Decode(&services); err != nil {
		return nil, err
	}
	return &services, nil
}

func GetBrands(countryCode string, service string, latitude float32, longitude float32) (*Brands, error) {
	client := http.Client{}
	url := fmt.Sprintf("https://core.honestbee.com/api/brands?countryCode=%s&page=1&page_size=6&serviceType=%s&latitude=%f&longitude=%f", countryCode, service, latitude, longitude)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/vnd.honestbee+json;version=2")
	req.Header.Add("Accept-Language", "zh-TW")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	brands := Brands{}
	if err := jsoniter.NewDecoder(resp.Body).Decode(&brands); err != nil {
		return nil, err
	}
	return &brands, nil
}
