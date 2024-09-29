package sms

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type TwilioClient struct {
	accountSID          string
	authToken           string
	messagingServiceSid string
}

func NewTwilioClient(accountSID, authToken, messagingServiceSid string) *TwilioClient {
	return &TwilioClient{
		accountSID:          accountSID,
		authToken:           authToken,
		messagingServiceSid: messagingServiceSid,
	}
}

func (c *TwilioClient) SendSMSAsync(to, message string) {
	go func() {
		err := c.SendSMS(to, message)
		if err != nil {
			// Handle error, perhaps log it
			log.Println(err)
		}
	}()
}

func (c *TwilioClient) SendSMS(to, message string) error {
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", c.accountSID)

	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("messagingServiceSid", c.messagingServiceSid)
	msgData.Set("Body", message)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(msgData.Encode()))
	req.SetBasicAuth(c.accountSID, c.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("Failed to send SMS: %v")
}
