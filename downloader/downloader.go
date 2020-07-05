package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

// FileFromURL downloads file(s) from baseURL and writes it to the specified filePath.
func FileFromURL(baseURL, filePath string, fileName ...string) error {
	for _, file := range fileName {
		dest := path.Join(filePath, file)
		fullURL := path.Join(baseURL, file)

		req, err := http.NewRequest("POST", fullURL, nil)
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
