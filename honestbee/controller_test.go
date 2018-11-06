package honestbee

import (
	"fmt"
	"testing"
)

// go test -v github.com/samwang0723/genghis-khan/honestbee
func TestGetServices(t *testing.T) {
	services, _ := GetServices("PH", 14.5367633, 121.009545)
	if len(*services) == 0 {
		t.Error("Cannot fetch country services")
	}
	msg := fmt.Sprintf("Successfully retrieve %d services", len(*services))
	t.Log(msg)
}
