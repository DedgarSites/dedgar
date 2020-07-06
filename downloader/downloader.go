package downloader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

type certFile struct {
	FileName string `json:"FileName"`
}

// FileFromURL downloads file(s) from baseURL and writes it to the specified filePath.
func FileFromURL(downloadURL, filePath string, fileName ...string) error {
	for _, file := range fileName {
		var cFile certFile

		dest := path.Join(filePath, file)

		jsonStr, err := json.Marshal(cFile)
		if err != nil {
			fmt.Println("SendData: Error marshalling json: ", err)
		}

		req, err := http.NewRequest("POST", downloadURL, bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Printf("Error getting pod info: %v \n", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		//client := &http.Client{Transport: httpClientWithSelfSignedTLS}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("makeClient: Error making API request: %v", err)
		}

		defer resp.Body.Close()

		err = os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}

		out, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
	}
	return nil
}
