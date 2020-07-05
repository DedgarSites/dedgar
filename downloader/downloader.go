package downloader

import (
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

		resp, err := http.Get(fullURL)
		if err != nil {
			return err
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
