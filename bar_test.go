package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIGetTimesFromPingPostSuccessfully(t *testing.T){
	type Request struct{
		Times int `json:"times"`
	}
	testRequest := Request{Times: 3}
	app:= fiber.New()

	app.Post("/ping", func(ctx *fiber.Ctx) error {
		body := Request{}
		if err := ctx.BodyParser(&body); err != nil{
			return err
		}
		return ctx.JSON(body)
	})

	testRequestAsByte,err := json.Marshal(testRequest)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost,"/ping", strings.NewReader(string(testRequestAsByte)))
	req.Header.Set(fiber.HeaderContentType,fiber.MIMEApplicationJSON)
	var resp *http.Response
	resp ,err = app.Test(req,1)
	if err != nil {
		t.Fatal(err)
	}
	body,_ := ioutil.ReadAll(resp.Body)
	resentRequest := Request{}
	err = json.Unmarshal(body,&resentRequest)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equalf(t, resentRequest.Times,testRequest.Times,"Sended times came correctly to the api")
}
func TestIProducePongsAsManyAsGivenTimes(t *testing.T){
	type Request struct{
		Times int `json:"times"`
	}
	testRequest := Request{Times: 3}
	app:= fiber.New()

	type PongsResponse struct{
		Pongs []string `json:"pongs"`
	}
	const PONG = "pong"
	constructPongResponseAmountOfGivenTimes := func(times int) PongsResponse{
		var pongsResponse = PongsResponse{}
		for i:=0 ; i<times; i++{
			pongsResponse.Pongs = append(pongsResponse.Pongs,PONG)
		}
		return pongsResponse
	}
	app.Post("/ping", func(ctx *fiber.Ctx) error {
		body := Request{}
		if err := ctx.BodyParser(&body); err != nil{
			return err
		}
		pongsResponse := constructPongResponseAmountOfGivenTimes(body.Times)
		return ctx.JSON(pongsResponse)
	})
	testRequestAsByte,err := json.Marshal(testRequest)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost,"/ping",strings.NewReader(string(testRequestAsByte)))
	req.Header.Set(fiber.HeaderContentType,fiber.MIMEApplicationJSON)
	resp,_ := app.Test(req,1)
	body,_ := ioutil.ReadAll(resp.Body)

	sentPongsResponse := PongsResponse{}
	err = json.Unmarshal(body,&sentPongsResponse)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equalf(t, len(sentPongsResponse.Pongs),3,"Pongs response that contains 'pong' word amount of given count")
}