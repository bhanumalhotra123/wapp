package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	pollInterval = time.Second * 2
)

const (
	endpoint = "https://api.open-meteo.com/v1/forecast" //?latitude=52.52&longitude=13.41&hourly=temperature_2m
)

type any = interface{}

type WeatherData struct {
	Elevation float64        `json:"elevation"`
	Hourly    map[string]any `json:"hourly"`
}

type Sender interface {
	Send(*WeatherData) error
}

type SMSSender struct {
	number string
}

func NewSMSSender(number string) *SMSSender {
	return &SMSSender{
		number: number,
	}
}

func (s *SMSSender) Send(data *WeatherData) error {
	fmt.Println("sending weather to number: ", s.number)
	return nil
}


type EmailSender struct {
	number string
}

func NewEmailSender(number string) *EmailSender {
	return &EmailSender{
		number: number,
	}
}

func (e *EmailSender) Send(data *WeatherData) error {
	fmt.Println("sending weather to email: ", e.email)
	return nil
}


type WPoller struct {
	closech chan struct{}
	senders  Sender
}

func NewPoller(senders ...Sender) *WPoller {
	return &WPoller{
		closech: make(chan struct{}),
		senders:  senders,
	}

}

func (wp *WPoller) close() {
	close(wp.closech)
}

func (wp *WPoller) start() {

	fmt.Println("starting the wpoller")

	ticker := time.NewTicker(pollInterval)

outer:
	for {
		select {
		case <-ticker.C:

			data, err := getWeatherResult(52.52, 13.41)
			if err != nil {
				log.Fatal(err)
			}
			if err := wp.handleData(data); err != nil {
				log.Fatal(err)

			}

		case <-wp.closech:
			//handle the graceful shutdown
			break outer
		}

	}
	fmt.Println("wpoller stopped gracefully")
}

func (wp *WPoller) handleData(data *WeatherData) error {
	// handle the data (store it in db)
	//send
	for _, s := range wp.senders {
		if err := s.Send(data); err != nil {
			return err := s.Send(data); err != nil {
				fmt.Println(err)
			}
		}
	}
	return wp.sender.Send(data)

}

func getWeatherResult(lat, long float64) (*WeatherData, error) {
	uri := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m", endpoint, lat, long)

	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data WeatherData

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func main() {
	smsSender := NewSMSSender("7404817684")
	WPoller := NewPoller(smsSender)
	WPoller.start()

}
