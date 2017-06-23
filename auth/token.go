package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func GetClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()

	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}

	tok, err := tokenFromFile(cacheFile)

	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}

	return config.Client(ctx, tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Open the link in your browser. Then type the "+
		"authorization code: \n%v\nPaste the token Google returns:\n", authURL)

	// Get code from command line
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	tokenCacheDir := filepath.Join(usr.HomeDir, ".gmail")
	os.MkdirAll(tokenCacheDir, 0700)

	return filepath.Join(tokenCacheDir,
		url.QueryEscape("gmail.json")), err
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)

	return t, err
}

func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)

	f, err := os.Create(file)

	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}

	defer f.Close()

	json.NewEncoder(f).Encode(token)
}

func GetClientSecrets() ([]byte, error) {
	usr, err := user.Current()

	if err != nil {
		return []byte(""), err
	}

	credentialsDirectory := filepath.Join(usr.HomeDir, ".credentials")
	secretsFile := filepath.Join(credentialsDirectory,
		url.QueryEscape("client_secret.json"))

	b, err := ioutil.ReadFile(secretsFile)

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	return b, err
}
