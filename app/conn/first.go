package conn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/indrenicloud/tricloud-agent/app/cmd"
	"github.com/indrenicloud/tricloud-agent/app/logg"
)

func RegisterAgent() {

	cf := GetConfig()

	if cf.UUID == "" {
		if cf.ApiKey == "" {
			logg.Log("Need api key")
			os.Exit(1)
		}

		client := &http.Client{}
		url := fmt.Sprintf("http://%s/registeragent/", cf.Url)
		//var url = cf.Url + "registeragent"
		req, _ := http.NewRequest("POST", url, nil)
		req.Header.Add("Api-key", cf.ApiKey)

		resp, err := client.Do(req)

		if err != nil {
			panic("server error")
		}
		body, err := ioutil.ReadAll(resp.Body)

		resbody := make(map[string]string)
		json.Unmarshal(body, &resbody)

		uid := resbody["data"]
		logg.Log("My ID:", uid)

		if uid == "" {
			logg.Log("Server didnot register us, every man for himself")
			os.Exit(1)
		}
		cf.UUID = uid
		SaveConfig()
	}

	updateSystemInfo()

}

func updateSystemInfo() {
	cf := GetConfig()
	rawb := cmd.GetSystemInfo()
	if rawb == nil {
		logg.Log("couldnot get systeminfo")
		return
	}
	//url := cf.Url + fmt.Sprintf("updatesysinfo/%s", cf.UUID)
	url := fmt.Sprintf("http://%s/updatesysinfo/%s", cf.Url, cf.UUID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(rawb))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logg.Log(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		logg.Log(err)
		return
	}
	logg.Log(string(body))
}
