package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/douglasmakey/go-fcm"
	"github.com/pkg/errors"
	"github.com/pusher/pusher-http-go/v5"
)

// NotificationPayload struct
type NotificationPayload struct {
	Title            string `json:"title,omitempty"`
	Body             string `json:"body,omitempty"`
	BodyLocKey       string `json:"body_loc_key,omitempty"`
	BodyLocArgs      string `json:"body_loc_args,omitempty"`
	Icon             string `json:"icon,omitempty"`
	Tag              string `json:"tag,omitempty"`
	Sound            string `json:"sound,omitempty"`
	Badge            string `json:"badge,omitempty"`
	Color            string `json:"color,omitempty"`
	ClickAction      string `json:"click_action,omitempty"`
	TitleLocKey      string `json:"title_loc_key,omitempty"`
	TitleLocArgs     string `json:"title_loc_args,omitempty"`
	AndroidChannelID string `json:"android_channel_id,omitempty"`
}

// Notifiable struct
type Notifiable struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Message struct
type Message struct {
	Channel            string     `json:"channel"`
	Message            string     `json:"message"`
	ID                 int        `json:"id"`
	CreatedAt          string     `json:"created_at"`
	Action             string     `json:"action"`
	NotifiedID         int        `json:"notified_id"`
	NotifierID         int        `json:"notifier_id"`
	Notifiable         Notifiable `json:"notifiable"`
	CanSendFCM         bool       `json:"can_send_fcm"`
	IsDailyQuestAction bool       `json:"is_daily_quest_action"`
	FCMTokens          []string   `json:"fcm_tokens"`
}

//FCMPayload sd
type FCMPayload struct {
	Title string `json:"title"`
}

// Data sd
type Data struct {
	Action string `json:"action"`
}

// FCMNotification sd
type FCMNotification struct {
	Notification    FCMPayload `json:"notification"`
	RegistrationIDs []string   `json:"registration_ids"`
	Data            Data       `json:"data"`
}

func getMessage(message string) (msgStruct Message, err error) {
	err = json.Unmarshal([]byte(message), &msgStruct)
	return msgStruct, err
}

// SendFCMNotification sends push notification using FCM
func SendFCMNotification(msg Message) {

	url := "https://fcm.googleapis.com/fcm/send"
	method := "POST"

	notif := FCMNotification{
		Notification: FCMPayload{
			Title: msg.Message,
		},
		RegistrationIDs: msg.FCMTokens,
		Data: Data{
			Action: msg.Action,
		},
	}

	payload, err := json.Marshal(&notif)
	if err != nil {
		fmt.Println("ERROR marhsaling FCM Payload")
		return
	}

	strPayload := string(payload)
	fmt.Println("FCM PAYLOAD TO SEND: ", strPayload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, strings.NewReader(strPayload))

	if err != nil {
		fmt.Println(err)
		return
	}

	req.Header.Add("Authorization", fmt.Sprintf("key=%s", os.Getenv("FCM_KEY")))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

// SendPusherNotification sends push notification using pusher
func SendPusherNotification(msg Message) error {

	pusherClient := pusher.Client{
		AppID:   os.Getenv("PUSHER_APP_ID"),
		Key:     os.Getenv("PUSHER_KEY"),
		Secret:  os.Getenv("PUSHER_SECRET"),
		Cluster: os.Getenv("PUSHER_CLUSTER"),
	}
	pusherClient.Secure = true

	// using struct and marhsalling adds back slash excape character
	data := map[string]interface{}{
		"message":     msg.Message,
		"id":          strconv.Itoa(msg.ID),
		"created_at":  msg.CreatedAt,
		"action":      msg.Action,
		"notified_id": strconv.Itoa(msg.NotifiedID),
		"notifier_id": strconv.Itoa(msg.NotifierID),
		"notifiable":  fmt.Sprintf("{name: %s , type: %s}", msg.Notifiable.Name, msg.Notifiable.Type),
	}

	fmt.Println("Pusher Payload to send: ", data)
	return pusherClient.Trigger(msg.Channel, "message", data)
}

// Handler is here
func Handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		msg, err := getMessage(message.Body)
		if err != nil {
			return errors.Wrapf(err, "ERROR: Cannot Unmarshal Message Body.")
		}

		// SEND PUSHER Notification
		err = SendPusherNotification(msg)
		if err != nil {
			return err
		}
		fmt.Println("Pusher Notification Successfull")

		// Check and send FCM notification
		if msg.CanSendFCM == true && msg.IsDailyQuestAction == true {
			SendFCMNotification(msg)
		}
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}

// SendFCMNotification1 TODO: Delete  as ppart of cleaning
func SendFCMNotification1(msg Message) error {

	// android, err := json.Marshal(`{"priority":"High"}`)
	// if err != nil {
	// 	return errors.Wrapf(err, "ERROR: sendFCMNotification => Cannot Marshal androidPriority.")
	// }
	// data := map[string]interface{}{
	// 	"message": msg.Message,
	// 	"notification": map[string]string{
	// 		"title": msg.Notifiable.Name,
	// 	},
	// 	"android": string(android),
	// }

	client := fcm.NewClient(os.Getenv("FCM_KEY"))
	//("AAAA-LPBm28:APA91bHnicuKflqAfUZeHOJbo1htfX5SfCzkb4xsI8IHGsnw6QvyOhfW57h-llzmUL558Y575UtXLjR2lhBsC_jx1YIJ2sG1ZOxcYuayewePkSsYkgNufANIAETlm-qiZSQlWsc32oOB")
	// // client.SetHTTPClient(client)
	// notif := NotificationPayload{
	// 	Title: "TEST WQR title",
	// }

	data := map[string]interface{}{
		"message": "From Go-FCM",
		"details": map[string]string{
			"name":  "Name",
			"user":  "Admin",
			"thing": "none",
		},
	}

	if len(msg.FCMTokens) == 0 {
		fmt.Println("No FCM Token Present!")
		return nil
	}

	client.PushMultiple(msg.FCMTokens, data)
	client.Message.Priority = "high"

	badRegistrations := client.CleanRegistrationIds()
	log.Println(badRegistrations)

	status, err := client.Send()
	if err != nil {
		return errors.Wrapf(err, "ERROR: FCM Could not send FCM Notification")
	}
	fmt.Println("FCM Notification sent successfully:", status.Results)

	return nil
}
