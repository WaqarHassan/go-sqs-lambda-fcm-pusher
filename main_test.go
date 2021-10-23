package main

import (
	"fmt"
	"log"
	"testing"
)

func Test_GetMessage(t *testing.T) {
	msg := `{ "channel": "notifications_NDZjqiuiKh54ezuX71iuf4Um", "message": "Your daily quest is ready", "id": 101284, "created_at": "2021-10-12T18:48:10.513+08:00", "action": "daily_quest_ready", "notified_id": 13, "notifiable": { "name": "Demo1stg", "type": "company" }, "can_send_fcm": false, "is_daily_quest_action": true, "fcm_tokens": [ "cb-M9-wYTj2ye1l-SG1_Em:APA91bHDML-Dj-ik6EO9y-_oY_Y-rsWA-1VKmBG0v5pWAlUvCRHFTk6ML9engug_E-BykQdgvCMie5HIPM93-EKiAuuuhQ9YDAFyhCOLrhMq6ObiC-p98kTMY5MQZZv3j6L7CXzpjeZX" ] }`
	// msg := `{ "channel": "notifications_NDZjqiuiKh54ezuX71iuf4Um"}`
	fmt.Println("=============================1111========================")

	msg2, err := getMessage(msg)
	log.Default().Println("Message", msg2)
	fmt.Println("\n\n MESSAGE \t\t", msg2.Notifiable.Type)
	if err != nil {
		log.Fatal(err)
	}

}
