package downloader

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

var (
	clusterCABundle = "/var/run/secrets/kubernetes.io/serviceaccount/service-ca.crt"
)

type certFile struct {
	FileName string `json:"FileName"`
}

// FileFromURL downloads file(s) from baseURL and writes it to the specified filePath.
func FileFromURL(downloadURL, filePath string, fileName ...string) error {
	insecure := flag.Bool("insecure-ssl", false, "Accept/Ignore all server SSL certificates")
	flag.Parse()

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	certs, err := ioutil.ReadFile(clusterCABundle)
	if err != nil {
		fmt.Printf("Failed to append %q to RootCAs: %v \n", clusterCABundle, err)
	}

	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		fmt.Println("No certs appended, using system certs only")
	}

	config := &tls.Config{
		InsecureSkipVerify: *insecure,
		RootCAs:            rootCAs,
	}
	tr := &http.Transport{TLSClientConfig: config}

	for _, file := range fileName {
		cFile := certFile{FileName: file}
		fmt.Println("searching for", downloadURL, filePath, file)

		dest := path.Join(filePath, file)

		jsonStr, err := json.Marshal(cFile)
		if err != nil {
			fmt.Println("Error marshalling json: ", err)
		}

		req, err := http.NewRequest("POST", downloadURL, bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Printf("Error making request: %v \n", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{Transport: tr}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making API request: %v \n", err)
		}
		defer resp.Body.Close()

		fmt.Println(resp.Status, resp.StatusCode)

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
