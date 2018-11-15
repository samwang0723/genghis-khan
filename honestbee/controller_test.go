package honestbee

import (
	"fmt"
	"testing"
)

// go test -v github.com/samwang0723/genghis-khan/honestbee
func TestGetServices(t *testing.T) {
	location := &Location{
		Latitude:  14.5367633,
		Longitude: 121.009545,
	}
	services, _ := GetServices("PH", location)
	if len(*services) == 0 {
		t.Error("Cannot fetch country services")
	}
	msg := fmt.Sprintf("Successfully retrieve %d services", len(*services))
	t.Log(msg)
}

func TestGetBrands(t *testing.T) {
	location := &Location{
		Latitude:  25.047571,
		Longitude: 121.577812,
	}
	brands, err := GetBrands("TW", "1", "groceries", location)
	if err != nil {
		t.Error(err.Error())
	}
	msg := fmt.Sprintf("Successfully retrieve %d brands", len(brands.Brands))
	t.Log(msg)
}

func TestGetDepartments(t *testing.T) {
	location := &Location{
		Latitude:  25.047571,
		Longitude: 121.577812,
	}
	departments, err := GetDepartments("11150", location)
	if err != nil {
		t.Error(err.Error())
	}
	msg := fmt.Sprintf("Successfully retrieve %d departments", len(departments.Departments))
	t.Log(msg)
}

func TestGetProducts(t *testing.T) {
	products, err := GetProducts("47306")
	if err != nil {
		t.Error(err.Error())
	}
	msg := fmt.Sprintf("Successfully retrieve %d products", len(*products.Products))
	t.Log(msg)
}

func TestSearchProducts(t *testing.T) {
	products, err := SearchProducts("11150", "Apple")
	if err != nil {
		t.Error(err.Error())
	}
	msg := fmt.Sprintf("Successfully retrieve %d products", len(*products.Products))
	t.Log(msg)
}
