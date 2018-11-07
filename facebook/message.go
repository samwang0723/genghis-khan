package facebook

import (
	"fmt"

	"github.com/samwang0723/genghis-khan/honestbee"
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
	Sender    User     `json:"sender,omitempty"`
	Recipient User     `json:"recipient,omitempty"`
	Timestamp int      `json:"timestamp,omitempty"`
	Message   Message  `json:"message,omitempty"`
	PostBack  PostBack `json:"postback,omitempty"`
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
	URL             string       `json:"url,omitempty"`
	Text            string       `json:"text,omitempty"`
	TemplateType    string       `json:"template_type,omitempty"`
	TopElementStyle string       `json:"top_element_style,omitempty"` // large, compact
	Coordinates     *Coordinates `json:"coordinates,omitempty"`
	Elements        *[]Element   `json:"elements,omitempty"`
	Buttons         *[]Button    `json:"buttons,omitempty"`
}

// SenderTypingAction - response with typing actions
func SenderTypingAction(event Messaging) *ActionResponse {
	response := ActionResponse{
		Recipient: User{
			ID: event.Sender.ID,
		},
		SenderAction: "typing_on",
	}
	return &response
}

// ParseLocation - parse latitude and longitude
func ParseLocation(event Messaging) *Coordinates {
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

func ComposeServicesButton(SenderID string, services *[]honestbee.Service) *Response {
	var buttons []Button
	for _, service := range *services {
		if service.Avaliable {
			buttons = append(buttons, Button{
				Title:   service.ServiceType,
				Type:    "postback",
				Payload: "selected_service",
			})
		}
	}

	response := Response{
		Recipient: User{
			ID: SenderID,
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

// ComposeLocation - response with location
func ComposeLocation(event Messaging) *Response {
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

func ComposeText(senderID string, message string) *Response {
	response := Response{
		Recipient: User{
			ID: senderID,
		},
		Message: Message{
			Text: message,
		},
	}
	return &response
}

//ComposeBrandList - response with brand list
func ComposeBrandList(event Messaging, brands honestbee.Brands) *Response {
	var elements []Element
	for _, brand := range brands.Brands {
		brandURL := fmt.Sprintf("https://www.honestbee.tw/zh-TW/%s/stores/%s", brand.ServiceType, brand.Slug)
		var buttons []Button
		buttons = append(buttons, Button{
			Title:               "View",
			Type:                "web_url",
			URL:                 brandURL,
			MessengerExtensions: true,
			WebViewHeightRatio:  "tall",
		})
		elements = append(elements, Element{
			Title:    brand.Name,
			SubTitle: brand.Description,
			ImageURL: brand.ImageURL,
			Buttons:  &buttons,
		})
	}

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
				},
			},
		},
	}
	return &response
}
