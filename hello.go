package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"What is your name?"`
}

type DashEvent struct {
	DeviceEvent struct {
		ButtonClicked struct {
			ClickType    string
			ReportedTime string
		}
	}
}

type MyResponse struct {
	Message string `json:"Answer:"`
}

type SlackMessage struct {
	Text string `json:"text"`
}

func main() {
	lambda.Start(hello)
}

func hello(event DashEvent) (MyResponse, error) {
	url := "INCOMING_WEBHOOK_URL"
	SendSlackMessage(url, event.DeviceEvent.ButtonClicked.ClickType)
	return MyResponse{Message: fmt.Sprintf("Hello %s!!", "event.Name")}, nil
}

func SendSlackMessage(url string, clickType string) error {

	jstZone := time.FixedZone("Asia/Tokyo", 9*60*60)
	jst := time.Now().UTC().In(jstZone)
	timeFormat := "15時04分"
	nowDateTimeString := fmt.Sprint(jst.Format(timeFormat))

	msg := ""

	if clickType == "SINGLE" {
		msg = "前駆陣痛開始"
	} else {
		msg = "前駆陣痛終了"
	}

	message := SlackMessage{Text: fmt.Sprintf("*【%v】* `%v`", nowDateTimeString, msg)}

	SendLine(fmt.Sprintf("*【%v】* `%v`", nowDateTimeString, msg))

	requestParam, _ := json.Marshal(message)

	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(requestParam),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}

func SendLine(message string) {

	url := "https://notify-api.line.me/api/notify"
	method := "POST"

	payload := strings.NewReader("message=" + message)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer {LINE_ACCESS_TOKEN}")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

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

// build command
// GOOS=linux GOARCH=amd64 go build -o hello
// zip handler.zip ./hello
