package facebook

import "log"

type Callback struct {
	Object string `json:"object,omitempty"`
	Entry  []struct {
		ID        string      `json:"id,omitempty"`
		Time      int         `json:"time,omitempty"`
		Messaging []Messaging `json:"messaging,omitempty"`
	} `json:"entry,omitempty"`
}

type Messaging struct {
	Sender    User    `json:"sender,omitempty"`
	Recipient User    `json:"recipient,omitempty"`
	Timestamp int     `json:"timestamp,omitempty"`
	Message   Message `json:"message,omitempty"`
}

type User struct {
	ID string `json:"id,omitempty"`
}

type QuickReply struct {
	ContentType string `json:"content_type,omitempty"`
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
	TemplateType    string       `json:"template_type,omitempty"`
	TopElementStyle string       `json:"top_element_style,omitempty"` // large, compact
	Coordinates     *Coordinates `json:"coordinates,omitempty"`
	Elements        *[]Element   `json:"elements,omitempty"`
	Buttons         *[]Button    `json:"buttons,omitempty"`
}

//ComposeLocation - response with location
func ComposeLocation(event Messaging) *Response {
	for _, attachment := range *event.Message.Attachments {
		coordinates := attachment.Payload.Coordinates
		log.Printf("User's location %f, %f", coordinates.Lat, coordinates.Long)
	}

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

//ComposeBrandList - response with brand list
func ComposeBrandList(event Messaging) *Response {
	var buttons []Button
	buttons = append(buttons, Button{
		Title:               "View",
		Type:                "web_url",
		URL:                 "https://peterssendreceiveapp.ngrok.io/collection",
		MessengerExtensions: true,
		WebViewHeightRatio:  "tall",
		FallbackURL:         "https://peterssendreceiveapp.ngrok.io/",
	})
	var elements []Element
	elements = append(elements, Element{
		Title:    "Classic T-Shirt Collection",
		SubTitle: "See all our colors",
		ImageURL: "https://peterssendreceiveapp.ngrok.io/img/collection.png",
		Buttons:  &buttons,
	})
	elements = append(elements, Element{
		Title:    "Classic Blue T-Shirt",
		SubTitle: "100% Cotton, 200% Comfortable",
		ImageURL: "https://peterssendreceiveapp.ngrok.io/img/collection.png",
		Buttons:  &buttons,
	})
	response := Response{
		Recipient: User{
			ID: event.Sender.ID,
		},
		Message: Message{
			Text: "Please tell me your location",
			Attachment: &Attachment{
				Type: "template",
				Payload: Payload{
					TemplateType:    "list",
					TopElementStyle: "large",
					Elements:        &elements,
				},
			},
		},
	}
	return &response
}
