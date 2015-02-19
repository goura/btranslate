package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	MSTRANS_SCOPE          = "http://api.microsofttranslator.com"
	TOKEN_SERVICE_URL      = "https://datamarket.accesscontrol.windows.net/v2/OAuth2-13"
	TRANSLATE_ENDPOINT_URL = "http://api.microsofttranslator.com/V2/Http.svc/Translate"
)

type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn string `json:"expires_in"`
}

type ClientCreds struct {
	ClientID     string
	ClientSecret string
}

func ObtainClientCreds() *ClientCreds {
	creds := ClientCreds{}
	creds.ClientID = os.Getenv("BTRANSLATE_CLIENT_ID")
	creds.ClientSecret = os.Getenv("BTRANSLATE_CLIENT_SECRET")
	return &creds
}

func ObtainAccessToken(creds *ClientCreds) (*AccessToken, error) {
	at := &AccessToken{}
	c := http.Client{}

	data := url.Values{}
	data.Set("client_id", creds.ClientID)
	data.Set("client_secret", creds.ClientSecret)
	data.Set("scope", MSTRANS_SCOPE)
	data.Set("grant_type", "client_credentials")

	u := TOKEN_SERVICE_URL

	resp, err := c.PostForm(u, data)
	if err != nil {
		fmt.Println(err)
		return at, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return at, err
	}

	err = json.Unmarshal(buf, at)
	if err != nil {
		fmt.Println(err)
		return at, err
	}

	return at, nil
}

func TranslateQuery(from, to, text string, accessToken *AccessToken) (string, error) {
	// Build request
	v := url.Values{}
	v.Set("from", from)
	v.Set("to", to)
	v.Set("text", text)

	u := TRANSLATE_ENDPOINT_URL + "?" + v.Encode()

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+accessToken.Token)

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	result := ""
	err = xml.Unmarshal(buf, &result)

	return result, err
}

func main() {
	// Obtain Creds
	creds := ObtainClientCreds()

	// Obtain the Access Token
	accessToken, err := ObtainAccessToken(creds)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ObtainAccessToken failed: %s\n", err)
		os.Exit(1)
	}

	// Obtain from, to, text, round_trip flag
	var from, to, text string
	var roundTrip, jsonOut bool

	// Set flags from command ling args
	flag.StringVar(&from, "from", "ja", "Language code of the text to translate")
	flag.StringVar(&to, "to", "en", "Language code to translate the text to")
	flag.StringVar(&text, "text", "", "Text to translate. If this is omitted, the program reads text from the standard input.")
	flag.BoolVar(&jsonOut, "json", false, "Output in json")
	flag.BoolVar(&roundTrip, "round_trip", false, "To round trip or not")
	flag.Parse()

	// Get text from stdin
	if text == "" {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Read from stdin failed: %s\n", err)
			os.Exit(1)
		}
		text = string(bytes)
	}

	// Fire query
	s, err := TranslateQuery(from, to, text, accessToken)

	if err != nil {
		fmt.Fprintf(os.Stderr, "TranslateQuery failed: %s\n", err)
		os.Exit(1)
	}

	// If round_trip flag is false, end here
	if !(roundTrip || jsonOut) {
		fmt.Println(s)
		os.Exit(0)
	}

	// Reverse translate
	rs, err := TranslateQuery(to, from, s, accessToken)
	if err != nil {
		fmt.Fprintf(os.Stderr, "TranslateQuery failed (reverse): %s\n", err)
		os.Exit(1)
	}

	if !jsonOut {
		fmt.Println(rs)
		os.Exit(0)
	}

	m := map[string]string{
		"from":          from,
		"to":            to,
		"original":      text,
		"translated":    s,
		"round_tripped": rs,
	}

	jsbs, err := json.Marshal(m)
	if err != nil {
		fmt.Fprintf(os.Stderr, "json.Unmarshal failed: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsbs))
}
