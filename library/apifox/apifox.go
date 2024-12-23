package apifox

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"weihu_server/library/config"
)

func Sync() {
	files, err := ioutil.ReadDir("./docs")
	if err != nil {
		fmt.Println(err.Error())
	}
	urlPath := config.GetString("apifox.url")
	token := config.GetString("apifox.token")
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") == false {
			continue
		}

		b, err := ioutil.ReadFile("./docs/" + f.Name())
		if err != nil {
			fmt.Println(err.Error())
		}

		payload := strings.NewReader(`{
			"input": ` + string(b) + `,
			"options": {
				"targetEndpointFolderId": 0,
				"targetSchemaFolderId": 0,
				"endpointOverwriteBehavior": "OVERWRITE_EXISTING",
				"schemaOverwriteBehavior": "OVERWRITE_EXISTING",
				"updateFolderOfChangedEndpoint": false,
				"prependBasePath": false
			}
		}`)

		client := &http.Client{}
		req, _ := http.NewRequest("POST", urlPath, payload)
		req.Header.Add("X-Apifox-Version", "2024-03-28")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer resp.Body.Close()
		bys, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		log.Println(string(bys))

	}

	//os.RemoveAll("./docs")
}
