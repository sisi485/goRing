package hue

import (
	"net/http"
	"fmt"
	"encoding/json"
	"bytes"
)

const rootUrl = "http://192.168.178.22/api/p2Hk18EfsISvatK3lrxF13j3rYnCoOF2XfRMwPFG"

func GetLights() (map[string]interface{}, error) {

	resp, err := http.Get(fmt.Sprintf("%s/lights", rootUrl))
	if err != nil {
		return nil, err
	}

	var buff map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&buff)

	return buff, err
}

func SetGroup(grpId int, grp map[string]interface{}) (map[string]interface{}, error) {

	data, err := json.Marshal(grp)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(data))
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/groups/%d/action", rootUrl, grpId), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var buff map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&buff)

	return buff, err
}


func CreateScene(name string, lights []string) ([]interface{}, error) {

	bLights := map[string]interface{}{
		"name": name,
		"lights": lights,
		"recycle": true,
	}

	data, err := json.Marshal(bLights)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/scenes/", rootUrl), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var buff []interface{}
	json.NewDecoder(resp.Body).Decode(&buff)

	return buff, err
}

func DeleteScene(id string) (map[string]interface{}, error) {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/scenes/%s", rootUrl, id), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}


	var buff map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&buff)

	return buff, err
}