package conn

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/indrenicloud/tricloud-agent/app/logg"

	"github.com/indrenicloud/tricloud-agent/app/cmd"
)

func updateSystemInfo(uuid string) {

	rawb := cmd.GetSystemInfo()
	if rawb == nil {
		panic("couldnot get systeminfo")
	}
	url := fmt.Sprintf("http://localhost:8081/updatesysinfo/%s", uuid)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rawb))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	logg.Log(string(body))
	defer resp.Body.Close()
}
