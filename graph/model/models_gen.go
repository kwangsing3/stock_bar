// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type DailyRecord struct {
	Date          string `json:"date"`
	TradingVolume string `json:"tradingVolume"`
	TradingPrice  string `json:"tradingPrice"`
	OpenPrice     string `json:"openPrice"`
	HighestPrice  string `json:"highestPrice"`
	LowestPrice   string `json:"lowestPrice"`
	ClosePrice    string `json:"closePrice"`
	PriceDiff     string `json:"priceDiff"`
	TransAmount   string `json:"transAmount"`
}

type DeleteRecord struct {
	Code string `json:"code"`
	Date string `json:"date"`
}

type NewRecord struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Date          string `json:"date"`
	TradingVolume string `json:"tradingVolume"`
	TradingPrice  string `json:"tradingPrice"`
	OpenPrice     string `json:"openPrice"`
	HighestPrice  string `json:"highestPrice"`
	LowestPrice   string `json:"lowestPrice"`
	ClosePrice    string `json:"closePrice"`
	PriceDiff     string `json:"priceDiff"`
	TransAmount   string `json:"transAmount"`
}

type NewStock struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Stock struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
