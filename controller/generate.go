package controller

// package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type SendReplicateJson struct {
	Version string                  `json:"version"`
	Input   *SendInputReplicateJson `json:"input"`
}
type SendInputReplicateJson struct {
	Img     string `json:"img"`
	Version string `json:"version"`
	Scale   int    `json:"scale"`
}

func Generate(w http.ResponseWriter, r *http.Request) {
	// taking imageUrl element from req.body
	req_body := make(map[string]string)
	req_body["imageUrl"] = ""
	json.NewDecoder(r.Body).Decode(&req_body)

	// DEBUG
	fmt.Printf("=========================================================================\n")
	fmt.Printf("%#v\n", req_body)

	// resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(json_data))

	// new_req, err := http.NewRequest("POST", "https://api.replicate.com/v1/predictions", []byte(req_body))

	// new_body := make(map[string]interface{}{
	// 	"version": "9283608cc6b7be6b65a8e44983db012355fde4132009bf99d976b2f0896856a3",
	// 	"input":   make(map[string]interface{}{"img": req_body["imageUrl"], "version": "v1.4", "scale": 2}),
	// })

	sendReplicateBody := SendReplicateJson{Version: "9283608cc6b7be6b65a8e44983db012355fde4132009bf99d976b2f0896856a3", Input: &SendInputReplicateJson{Img: req_body["imageUrl"], Version: "v1.4", Scale: 2}}
	jsonObject, err := json.Marshal(sendReplicateBody)
	if err != nil {
		log.Fatal(err)
	}

	// DEBUG
	fmt.Printf("%#v\n", sendReplicateBody)

	// new_req, err := http.NewRequest("POST", "https://api.replicate.com/v1/predictions", bytes.NewBuffer(req_body.imageUrl))
	// sendReqReplicate, err := http.NewRequest("POST", "https://api.replicate.com/v1/predictions", []byte(jsonObject))

	authenticationString := "Token <Token>"
	sendReqReplicate, err := http.NewRequest("POST", "https://api.replicate.com/v1/predictions", bytes.NewBuffer(jsonObject))

	sendReqReplicate.Header.Set("Content-Type", "application/json")
	sendReqReplicate.Header.Add("Authorization", authenticationString)

	client := &http.Client{}
	sendRespReplicate, err := client.Do(sendReqReplicate)
	if err != nil {
		log.Fatal(err)
	}
	defer sendRespReplicate.Body.Close()

	respBody, err := ioutil.ReadAll(sendRespReplicate.Body)
	if err != nil {
		log.Fatal(err)
	}
	var sendRespReplicateJson map[string]interface{}
	json.Unmarshal(respBody, &sendRespReplicateJson)

	endImageUrl := "https://api.replicate.com/v1/predictions/" + sendRespReplicateJson["id"].(string)

	// DEBUG
	fmt.Printf("%v\n", endImageUrl)

	returnedImage := ""
	for returnedImage == "" {
		fmt.Println("polling for result...")

		sendGetReqReplicate, err := http.NewRequest("GET", endImageUrl, nil)
		if err != nil {
			log.Fatal(err)
		}

		sendGetReqReplicate.Header.Set("Content-Type", "application/json")
		sendGetReqReplicate.Header.Add("Authorization", authenticationString)

		getClient := &http.Client{}
		sendGetRespReplicate, err := getClient.Do(sendGetReqReplicate)
		if err != nil {
			log.Fatal(err)
		}

		getRespBody, err := ioutil.ReadAll(sendGetRespReplicate.Body)
		if err != nil {
			log.Fatal(err)
		}
		var sendGetRespReplicateJson map[string]interface{}
		json.Unmarshal(getRespBody, &sendGetRespReplicateJson)

		if sendGetRespReplicateJson["status"] == "succeeded" {
			// returnedOutput := sendGetRespReplicateJson["output"].([]string)
			// returnedImage = returnedOutput[0]
			returnedImage = sendGetRespReplicateJson["output"].(string)
		} else if sendGetRespReplicateJson["status"] == "failed" {
			break
		} else {
			time.Sleep(2 * time.Second)
		}
	}

	json.NewEncoder(w).Encode(returnedImage)
}

// func main() {
// 	Generate(&http.ResponseWriter, http.Request)
// }
