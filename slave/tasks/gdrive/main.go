package gdrive

/*

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"google.golang.org/api/googleapi"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(oauthConfig *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := config.STORAGE_DIR + "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(oauthConfig)
		saveToken(tokFile, tok)
	}
	return oauthConfig.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	// FIXME make oauth2 auth more user friendly
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	// TODO make sure that token is valid and refreshed
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func ClientInit() *drive.Service {
	b, err := ioutil.ReadFile(config.STORAGE_DIR + "credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveAppdataScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	// FIXME don't use deprecated method
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	return srv
}

func UploadFile(service *drive.Service, content io.Reader, filename string) (*drive.File, error) {
	fileUploaded, err := service.Files.Create(&drive.File{
		MimeType: "application/gzip",
		Name:     filename,
		Parents:  []string{"appDataFolder"},
	}).Media(content).Do()

	if err != nil {
		logrus.Errorf("%v", err)
	}
	return fileUploaded, nil
}

func ListFiles(service *drive.Service) {
	fileList, err := service.Files.
		List().
		Spaces("appDataFolder").
		Fields(googleapi.Field("nextPageToken, files(id, name, md5Checksum, quotaBytesUsed)")).
		PageSize(10).
		Do()
	if err != nil {
		logrus.Errorf("Unable to retrieve files: %v", err)
	}

	if len(fileList.Files) == 0 {
		logrus.Info("No backups found")
	} else {
		for _, fileInfo := range fileList.Files {
			logrus.Infof("id: %s ", fileInfo.Id)
			logrus.Infof("name: %s ", fileInfo.Name)
			logrus.Infof("md5: %s ", fileInfo.Md5Checksum)
			logrus.Infof("size: %v MB", (fileInfo.QuotaBytesUsed/1024)/1024)
			logrus.Info("-")
		}
	}
}
*/
