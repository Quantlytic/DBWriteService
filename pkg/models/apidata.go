package models

type ApiData struct {
	Timestamp string `json:"timestamp"`
	Source    string `json:"source"`
	Exchange  string `json:"exchange"`
	Symbol    string `json:"symbol"`
	Data      string `json:"data"`
}
