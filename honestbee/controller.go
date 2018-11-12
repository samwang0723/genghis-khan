package honestbee

import (
	"bytes"
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

type Category struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	ImageURL      string `json:"imageUrl"`
	ProductsCount int    `json:"productsCount"`
}

type Department struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	ImageURL      string      `json:"imageUrl"`
	ProductsCount int         `json:"productsCount"`
	Categories    *[]Category `json:"categories"`
}

type Departments struct {
	Departments []Department `json:"departments"`
}

type Products struct {
	Products *[]Product `json:"products"`
}

type Product struct {
	ID               int    `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	ImageURL         string `json:"imageUrl"`
	PreviewImageURL  string `json:"previewImageUrl"`
	ImageURLBasename string `json:"imageUrlBasename"`
	Currency         string `json:"currency"`
	MaxQuantity      string `json:"maxQuantity"`
	Slug             string `json:"slug"`
	UnitType         string `json:"unitType"`
	SoldBy           string `json:"soldBy"`
	AmountPerUnit    string `json:"amountPerUnit"`
	Size             string `json:"size"`
	Status           string `json:"status"`
	Price            string `json:"price"`
	NormalPrice      string `json:"normalPrice"`
	PackingSize      string `json:"packingSize"`
	Alcohol          bool   `json:"alcohol"`
}

type SearchQuery struct {
	Page     int    `json:"page,omitempty"`
	PageSize int    `json:"pageSize,omitempty"`
	Platform string `json:"platform,omitempty"`
	Q        string `json:"q,omitempty"`
	UserID   string `json:"userId,omitempty"`
	UUID     string `json:"uuid,omitempty"`
}

const BRANDS = "brands"
const DEPARTMENTS = "departments"
const PRODUCTS = "products"
const PRODUCT = "product"
const SEARCH = "search"
const LOGIN_URL = "https://tranquil-anglerfish.glitch.me/login"

func SearchProducts(storeID string, query string) (*Products, error) {
	client := http.Client{}
	url := fmt.Sprintf("https://core.honestbee.com/api/stores/%s", storeID)

	queryJson := SearchQuery{
		Page:     1,
		PageSize: 4,
		Platform: "iOS",
		Q:        query,
		UserID:   "",
		UUID:     "508786e0-57b8-4252-87d6-13295a81733a",
	}
	data, err := jsoniter.Marshal(queryJson)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", url, bytes.NewBuffer(data))
	req.Header.Add("Accept", "application/vnd.honestbee+json;version=2")
	req.Header.Add("Accept-Language", "zh-TW")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	products := Products{}
	if err := jsoniter.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}
	return &products, nil
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

func GetBrands(countryCode string, page string, service string, latitude float32, longitude float32) (*Brands, error) {
	client := http.Client{}
	url := fmt.Sprintf("https://core.honestbee.com/api/brands?countryCode=%s&page=%s&page_size=3&serviceType=%s&latitude=%f&longitude=%f", countryCode, page, service, latitude, longitude)
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

func GetDepartments(storeID string, latitude float32, longitude float32) (*Departments, error) {
	client := http.Client{}
	url := fmt.Sprintf("https://core.honestbee.com/api/stores/%s/directory?latitude=%f&longitude=%f", storeID, latitude, longitude)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/vnd.honestbee+json;version=2")
	req.Header.Add("Accept-Language", "zh-TW")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	departments := Departments{}
	if err := jsoniter.NewDecoder(resp.Body).Decode(&departments); err != nil {
		return nil, err
	}
	return &departments, nil
}

func GetProducts(departmentID string) (*Products, error) {
	client := http.Client{}
	url := fmt.Sprintf("https://core.honestbee.com/api/departments/%s?page=1&pageSize=10&sort=ranking", departmentID)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/vnd.honestbee+json;version=2")
	req.Header.Add("Accept-Language", "zh-TW")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	products := Products{}
	if err := jsoniter.NewDecoder(resp.Body).Decode(&products); err != nil {
		return nil, err
	}
	return &products, nil
}
