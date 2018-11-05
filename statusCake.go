package main

import (
	"github.com/DreamItGetIT/statuscake"
	"github.com/mtulio/statuscake-exporter/collector"
	stk "github.com/mtulio/statuscake-exporter/statusCake"
	log "github.com/sirupsen/logrus"
)

var (
	Stk stk.StkOptions
)

func initClient() string {
	c, err := statuscake.New(statuscake.Auth{
		Username: config.StkUsername,
		Apikey:   config.StkApikey,
	},
	)
	if err != nil {
		log.Fatal(err)
	}
	Stk.Client = *c
	collector.Stk = Stk
	return "Success"
}
