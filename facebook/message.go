package facebook

import (
	"errors"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/samwang0723/genghis-khan/honestbee"
	"github.com/samwang0723/genghis-khan/utils"
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
	Sender         User            `json:"sender,omitempty"`
	Recipient      User            `json:"recipient,omitempty"`
	Timestamp      int             `json:"timestamp,omitempty"`
	Message        Message         `json:"message,omitempty"`
	PostBack       *PostBack       `json:"postback,omitempty"`
	AccountLinking *AccountLinking `json:"account_linking,omitempty"`
}

type PostBack struct {
	Title   string `json:"title,omitempty"`
	Payload string `json:"payload,omitempty"`
}

type User struct {
	ID string `json:"id,omitempty"`
}

type QuickReply struct {
	ContentType string `json:"content_type,omitempty"`
	Title       string `json:"title,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
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

type ActionResponse struct {
	Recipient    User   `json:"recipient,omitempty"`
	SenderAction string `json:"sender_action,omitempty"`
}

type AccountLinking struct {
	AuthorizationCode string `json:"authorization_code,omitempty"`
	Status            string `json:"status,omitempty"`
}

type Coordinates struct {
	Lat  float32 `json:"lat,omitempty"`
	Long float32 `json:"long,omitempty"`
}

type DefaultAction struct {
	Type                string `json:"type,omitempty"`
	URL                 string `json:"url,omitempty"`
	MessengerExtensions bool   `json:"messenger_extensions,omitempty"`
	WebViewHeightRatio  string `json:"webview_height_ratio,omitempty"`
}

type Button struct {
	Title               string `json:"title,omitempty"`
	Type                string `json:"type,omitempty"`
	URL                 string `json:"url,omitempty"`
	MessengerExtensions bool   `json:"messenger_extensions,omitempty"`
	WebViewHeightRatio  string `json:"webview_height_ratio,omitempty"`
	FallbackURL         string `json:"fallback_url,omitempty"`
	Payload             string `json:"payload,omitempty"`
}

type Element struct {
	Title         string         `json:"title,omitempty"`
	SubTitle      string         `json:"subtitle,omitempty"`
	ImageURL      string         `json:"image_url,omitempty"`
	Buttons       *[]Button      `json:"buttons,omitempty"`
	DefaultAction *DefaultAction `json:"default_action,omitempty"`
}

type Payload struct {
	URL              string       `json:"url,omitempty"`
	Text             string       `json:"text,omitempty"`
	TemplateType     string       `json:"template_type,omitempty"`
	ImageAspectRatio string       `json:"image_aspect_ratio,omitempty"`
	TopElementStyle  string       `json:"top_element_style,omitempty"` // large, compact
	Coordinates      *Coordinates `json:"coordinates,omitempty"`
	Elements         *[]Element   `json:"elements,omitempty"`
	Buttons          *[]Button    `json:"buttons,omitempty"`
}

func (event *Messaging) saveAccessToken() (string, error) {
	key := fmt.Sprintf("login_%s", event.Sender.ID)
	err := utils.RedisClient().Set(key, event.AccountLinking.AuthorizationCode, 0).Err()
	if err != nil {
		str := fmt.Sprintf("Session store error: %s", err.Error())
		return "", errors.New(str)
	}
	return event.AccountLinking.AuthorizationCode, nil
}

func (event *Messaging) viewingStoreID() string {
	key := fmt.Sprintf("viewing_store_%s", event.Sender.ID)
	val, err := utils.RedisClient().Get(key).Result()
	if err == nil {
		return val
	}
	return ""
}

func (event *Messaging) saveViewingStoreID(storeID string) error {
	key := fmt.Sprintf("viewing_store_%s", event.Sender.ID)
	err := utils.RedisClient().Set(key, storeID, 0).Err()
	return err
}

func (event *Messaging) location() *honestbee.Location {
	key := fmt.Sprintf("location_%s", event.Sender.ID)
	val, err := utils.RedisClient().Get(key).Result()
	if err == nil {
		location := new(honestbee.Location)
		jsoniter.UnmarshalFromString(val, &location)
		return location
	}
	return nil
}

func (event *Messaging) saveLocation(location *honestbee.Location) error {
	key := fmt.Sprintf("location_%s", event.Sender.ID)
	json := fmt.Sprintf(`{"latitude":%f,"longitude":%f}`, location.Latitude, location.Longitude)
	return utils.RedisClient().Set(key, json, 0).Err()
}

func (event *Messaging) senderTypingAction() *ActionResponse {
	response := ActionResponse{
		Recipient: User{
			ID: event.Sender.ID,
		},
		SenderAction: "typing_on",
	}
	return &response
}

func (event *Messaging) parseLocation() *Coordinates {
	if event.Message.Attachments != nil {
		for _, attachment := range *event.Message.Attachments {
			coordinates := attachment.Payload.Coordinates
			if coordinates != nil {
				return coordinates
			}
		}
	}
	return nil
}

func (event *Messaging) askLocation() *Response {
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
		},
	}
	return &response
}

func (event *Messaging) login() *Response {
	var buttons []Button
	buttons = append(buttons, Button{
		Type: "account_link",
		URL:  honestbee.LOGIN_URL,
	})

	response := Response{
		Recipient: User{
			ID: event.Sender.ID,
		},
		Message: Message{
			Attachment: &Attachment{
				Type: "template",
				Payload: Payload{
					TemplateType: "button",
					Text:         "Please login to have better shopping experience",
					Buttons:      &buttons,
				},
			},
		},
	}
	return &response
}

func (event *Messaging) listServices(services *[]honestbee.Service) *Response {
	var buttons []Button
	for _, service := range *services {
		if service.Avaliable {
			buttons = append(buttons, Button{
				Title:   service.ServiceType,
				Type:    "postback",
				Payload: fmt.Sprintf("%s:%s:1", honestbee.BRANDS, service.ServiceType),
			})
		}
	}

	response := Response{
		Recipient: User{
			ID: event.Sender.ID,
		},
		Message: Message{
			Attachment: &Attachment{
				Type: "template",
				Payload: Payload{
					TemplateType: "button",
					Text:         "These are the available services",
					Buttons:      &buttons,
				},
			},
		},
	}
	return &response
}

func (event *Messaging) showText(message string) *Response {
	response := Response{
		Recipient: User{
			ID: event.Sender.ID,
		},
		Message: Message{
			Text: message,
		},
	}
	return &response
}

func (event *Messaging) listProducts(products honestbee.Products) *Response {
	var elements []Element
	for _, product := range *products.Products {
		if product.Status == honestbee.STATUS_AVAILABLE {
			var buttons []Button
			buttons = append(buttons, Button{
				Title:   "Shop Now",
				Type:    "postback",
				Payload: fmt.Sprintf("%s:%d", honestbee.BUY_PRODUCT, product.ID),
			})
			elements = append(elements, Element{
				Title:    product.Title,
				SubTitle: fmt.Sprintf("%s (%s)\n$%s", product.ProductBrand, product.Size, product.Price),
				ImageURL: fmt.Sprintf("https://assets.honestbee.com/products/images/480/%s", product.ImageURLBasename),
				Buttons:  &buttons,
			})
		}
	}
	response := Response{
		Recipient: User{
			ID: event.Sender.ID,
		},
		Message: Message{
			Attachment: &Attachment{
				Type: "template",
				Payload: Payload{
					TemplateType:     "generic",
					ImageAspectRatio: "square",
					Elements:         &elements,
				},
			},
		},
	}
	return &response
}

func (event *Messaging) listDepartments(departments honestbee.Departments) *Response {
	index := 1
	var buttons []Button
	for _, department := range departments.Departments {
		buttons = append(buttons, Button{
			Title:   department.Name,
			Type:    "postback",
			Payload: fmt.Sprintf("%s:%d", honestbee.PRODUCTS, department.ID),
		})
		index = index + 1
		if index >= 3 {
			break
		}
	}

	response := Response{
		Recipient: User{
			ID: event.Sender.ID,
		},
		Message: Message{
			Attachment: &Attachment{
				Type: "template",
				Payload: Payload{
					TemplateType: "button",
					Text:         "Please choose one of the departments",
					Buttons:      &buttons,
				},
			},
		},
	}
	return &response
}

func (event *Messaging) listBrands(brands honestbee.Brands) *Response {
	var elements []Element
	for _, brand := range brands.Brands {
		var buttons []Button
		buttons = append(buttons, Button{
			Title:   "Browse",
			Type:    "postback",
			Payload: fmt.Sprintf("%s:%d", honestbee.SEARCH, brand.StoreID),
		})
		elements = append(elements, Element{
			Title:    brand.Name,
			SubTitle: brand.Description,
			ImageURL: brand.ImageURL,
			Buttons:  &buttons,
		})
	}

	var buttons []Button
	buttons = append(buttons, Button{
		Title:   "View More",
		Type:    "postback",
		Payload: fmt.Sprintf("%s:%s:%d", honestbee.BRANDS, brands.Brands[0].ServiceType, brands.Meta.CurrentPage+1),
	})

	response := Response{
		Recipient: User{
			ID: event.Sender.ID,
		},
		Message: Message{
			Attachment: &Attachment{
				Type: "template",
				Payload: Payload{
					TemplateType:    "list",
					TopElementStyle: "compact",
					Elements:        &elements,
					Buttons:         &buttons,
				},
			},
		},
	}
	return &response
}
