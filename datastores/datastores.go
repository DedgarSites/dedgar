package datastores

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dedgarsites/dedgar/models"
	"github.com/jinzhu/gorm"

	// Convention for gorm usage
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	// PostMap containes the names of eligible posts and their paths
	PostMap      = make(map[string]string)
	CookieSecret string
	OAuthID      string
	OAuthKey     string
	dbHost       string
	dbPort       string
	dbUser       string
	dbPass       string
	dbName       string
	Subject      string
	CharSet      string
	Sender       string
	Recipient    string
	AuthMap      map[string]bool
	psqlInfo     = fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	DB, _        = gorm.Open("postgres", psqlInfo)
)

func FindSummary(fpath string) string {
	file, err := os.Open(fpath + "_summary")
	if err != nil {
		return "No summary"
	}
	defer file.Close()

	var buffer bytes.Buffer
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		buffer.WriteString(line)
		//    if line == "<!--more-->" {
		//      break
		//    }
		//fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return buffer.String()
}

// Populates a map of postnames that gets checked every call to GET /post/:postname.
// We're running in a container, so populating this on startup works fine as we won't be adding
// any new posts while the container is running.
func FindPosts(dirpath string, extension string) map[string]string {
	if err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
		}
		if strings.HasSuffix(path, extension) {
			postname := strings.Split(path, extension)[0]
			summary := FindSummary(postname)
			//fmt.Println(summary)
			//fmt.Println(fmt.Sprintf("%T", summary))
			PostMap[filepath.Base(postname)] = summary
		}
		return err
	}); err != nil {
		panic(err)
	}
	return PostMap
}

func CheckDB() {
	if !DB.HasTable(&models.User{}) {
		fmt.Println("Creating users table")
		DB.CreateTable(&models.User{})
	}
}

func init() {
	var appSecrets models.AppSecrets

	filePath := "/secrets/dedgar_secrets.json"
	fileBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Println("Error loading secrets json: ", err)
	}

	err = json.Unmarshal(fileBytes, &appSecrets)
	if err != nil {
		fmt.Println("Error Unmarshaling secrets json: ", err)
	}

	CookieSecret = appSecrets.CookieSecret
	OAuthID = appSecrets.GoogleAuthID
	OAuthKey = appSecrets.GoogleAuthKey
	dbPass = appSecrets.PsqlPassword
	dbUser = appSecrets.PsqlUser
	dbPort = appSecrets.PsqlServicePort
	dbName = appSecrets.PsqlDatabase
	dbHost = appSecrets.PsqlServiceHost
	Subject = appSecrets.Subject
	CharSet = appSecrets.CharSet
	Sender = appSecrets.Sender
	AuthMap = appSecrets.AuthMap
	Recipient = appSecrets.Recipient
}
