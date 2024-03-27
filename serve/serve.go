package serve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sotowang/sunotoapi/cfg"
	"github.com/sotowang/sunotoapi/models"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

var (
	SessionExp int64
	Session    string
)

func GetSession(clientSession string) string {
	_url := "https://clerk.suno.ai/v1/client?_clerk_js_version=4.70.5"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, _url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "__client="+clientSession)
	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return ""
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Print("Error")
		return ""
	}
	body, _ := io.ReadAll(res.Body)
	var data models.GetSessionData
	if err = json.Unmarshal(body, &data); err != nil {
		log.Print(err)
		return ""
	}
	SessionExp = data.Response.Sessions[0].ExpireAt
	return data.Response.Sessions[0].Id
}

func GetJwtToken(clientSession string) (string, error) {
	if time.Now().After(time.Unix(SessionExp, 0)) {
		Session = GetSession(clientSession)
	}
	_url := fmt.Sprintf("https://clerk.suno.ai/v1/client/sessions/%s/tokens?_clerk_js_version=4.70.5", Session)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, _url, nil)

	if err != nil {
		log.Print(err)
		return "", err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", "__client="+cfg.Config.App.Client)

	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != 200 {
		//log.Print(string(body))
		return "", fmt.Errorf(string(body))
	}
	var data models.GetTokenData
	if err = json.Unmarshal(body, &data); err != nil {
		log.Print(err)
		return "", err
	}
	//有效时间 1 分钟
	if len(data.Jwt) == 0 {
		log.Print("GetJwtToken: ", data.Jwt)
		return "", err
	}
	return data.Jwt, nil
}

func sendRequest(url, method string, data []byte) ([]byte, error) {
	jwt, err := IsJWTExpired(cfg.Config.App.Client)
	if err != nil {
		log.Println("Error getting JWT: ", err)
		return nil, err
	}

	client := &http.Client{}
	var req *http.Request
	if data != nil {
		req, err = http.NewRequest(method, url, bytes.NewReader(data))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		log.Print(err)
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	req.Header.Add("Authorization", "Bearer "+jwt)

	res, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	return body, nil
}

func V2Generate(d map[string]interface{}) ([]byte, error) {
	_url := "https://studio-api.suno.ai/api/generate/v2/"
	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Fatalf("Error marshalling request data: %v", err)
		return nil, err
	}
	body, err := sendRequest(_url, "POST", jsonData)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func V2GetFeedTask(ids string) ([]byte, error) {
	ids = url.QueryEscape(ids)
	_url := "https://studio-api.suno.ai/api/feed/?ids=" + ids
	body, err := sendRequest(_url, "GET", nil)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GenerateLyrics(d map[string]interface{}) ([]byte, error) {
	_url := "https://studio-api.suno.ai/api/generate/lyrics/"
	jsonData, err := json.Marshal(d)
	if err != nil {
		log.Fatalf("Error marshalling request data: %v", err)
		return nil, err
	}
	body, err := sendRequest(_url, "POST", jsonData)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetLyricsTask(ids string) ([]byte, error) {
	_url := "https://studio-api.suno.ai/api/generate/lyrics/" + ids
	body, err := sendRequest(_url, "GET", nil)
	if err != nil {
		return nil, err
	}
	return body, nil
}
