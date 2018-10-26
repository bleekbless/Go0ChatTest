package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type RecipientResponse struct {
	UserId                string `json:"userId"`
	ConversationId        string `json:"conversationId"`
	ConversationReference string `json:"conversationReference"`
	UserName              string `json:"userName"`
}

type RecipientFiled struct {
	ID string `json:"id"`
}

func getRule(url, body string) (rule ChatRule, ok bool) {
	for _, r := range config.Rules {
		if r.Request.URL == url && strings.Contains(body, r.Request.BodySegment) {
			rule, ok = r, true
			return
		}
	}
	return
}

func sendResponse(rule ChatRule, reqBody string) {
	statChan <- Stat{Name: rule.Request.Name, ReqCount: 1, LastRequest: time.Now()}

	if config.ResponsePauseMax > 0 {
		ms := rand.Intn(config.ResponsePauseMax - config.ResponsePauseMin)
		ms += config.ResponsePauseMin
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
	fmt.Printf("Estou entrando aqui agora\n\n")
	body := rule.Response.Body

	// get RecipientId from request and use in Response
	var recResp RecipientResponse
	err := json.Unmarshal([]byte(reqBody), &recResp)
	if err != nil {
		fmt.Errorf("Error on repicient unmarshaling: %v", err)
	} else if recResp.ConversationId != "" {
		body = strings.Replace(body, "[cId]", recResp.ConversationId, -1)
		body = strings.Replace(body, "[cReference]", recResp.ConversationReference, -1)
		body = strings.Replace(body, "[uName]", recResp.UserName, -1)
		body = strings.Replace(body, "[uId]", recResp.UserId, -1)

	}
	fmt.Printf("[%v]\n", body)

	// set valid timestamp
	// t := time.Now()
	// ts := strconv.Itoa(int(t.UnixNano() / 1000000))
	// body = strings.Replace(body, "[timestamp]", ts, -1)

	// right quotes
	body = strings.Replace(body, "'", "\"", -1)

	client := &http.Client{}

	req, err := http.NewRequest("POST", config.AppURL, strings.NewReader(body))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Auth-Token", "aec6a4ad-22ba-4d85-bb88-ca5a0c259066")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		statChan <- Stat{Name: rule.Request.Name, RespCount: 1, RespNetErr: 1, LastResponse: time.Now()}
	} else {
		fmt.Printf("Olha só %v", res)
		if res.StatusCode == 200 {
			statChan <- Stat{Name: rule.Request.Name, RespCount: 1, LastResponse: time.Now()}
		} else {
			statChan <- Stat{Name: rule.Request.Name, RespCount: 1, RespStatusErr: 1, LastResponse: time.Now()}
		}
	}
}

func readRequest(r *http.Request) error {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	body := string(bodyBytes)

	if debug {
		fmt.Printf("[%v] %v == %v\n", r.Method, r.URL, body)
	}

	if config.RequestTimeMax > 0 {
		ms := rand.Intn(config.RequestTimeMax - config.RequestTimeMin)
		ms += config.RequestTimeMin
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}

	rule, ok := getRule(r.URL.Path, body)
	if !ok {
		statChan <- Stat{Name: otherRequests, ReqCount: 1, LastRequest: time.Now()}
		return nil
	}
	if rule.Response.URL == "" {
		statChan <- Stat{Name: rule.Request.Name, ReqCount: 1, LastRequest: time.Now()}
		return nil
	}

	go sendResponse(rule, body)

	return nil
}
