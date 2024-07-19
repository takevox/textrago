package textra

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/takevox/textrago/api"
	"golang.org/x/oauth2/clientcredentials"
)

type TexTra struct {
	Config            *Config
	AccessToken       string
	AccessTokenExpiry time.Time
}

func NewTexTra(conf *Config) (*TexTra, error) {
	tt := &TexTra{
		Config: conf,
	}
	return tt, nil
}

/*
アクセストークンをクリア
*/
func (tt *TexTra) _ClearAccessToken() {
	tt.AccessToken = ""
	tt.AccessTokenExpiry = time.Time{}
}

/*
アクセストークンを更新
*/
func (tt *TexTra) _RefreshAccessToken() error {
	conf := &clientcredentials.Config{
		ClientID:     tt.Config.API_KEY,
		ClientSecret: tt.Config.API_SECRET,
		TokenURL:     fmt.Sprintf("%s/oauth2/token.php", tt.Config.BaseURL),
	}

	ctx := context.Background()

	token, err := conf.Token(ctx)
	if err != nil {
		return err
	}

	tt.AccessToken = token.AccessToken
	tt.AccessTokenExpiry = token.Expiry

	return nil
}

/*
アクセストークンの取得
*/
func (tt *TexTra) GetAccessToken() (string, error) {
	if !tt.HasToken() {
		tt._ClearAccessToken()
	}
	err := tt._RefreshAccessToken()
	if err != nil {
		return "", err
	}
	return tt.AccessToken, nil
}

/*
有効なAccessTokenを保持しているか
*/
func (tt *TexTra) HasToken() bool {
	if tt.AccessToken == "" {
		return false
	}
	if tt.AccessTokenExpiry.Unix() < time.Now().Unix() {
		return false
	}
	return true
}

/*
言語を検出
*/
func (tt *TexTra) DetectLanguage(text string) ([]DetectLanguageResponse, error) {
	token, err := tt.GetAccessToken()
	if err != nil {
		return nil, err
	}

	values := &url.Values{}
	values.Add("access_token", token)
	values.Add("key", tt.Config.API_KEY)
	values.Add("name", tt.Config.UserName)
	values.Add("type", "json")
	values.Add("text", text)

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api/langdetect/", tt.Config.BaseURL),
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	res := &api.LangDetectResponse{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	if res.ResultSet.Code != 0 {
		return nil, fmt.Errorf("response error (%v) %v", res.ResultSet.Code, res.ResultSet.Message)
	}

	result := []DetectLanguageResponse{}

	for _, v := range res.ResultSet.Result.LangDetect {
		result = append(result, DetectLanguageResponse{
			Lang: v.Lang,
			Rate: v.Rate,
		})
	}

	return result, nil
}

/*
自動翻訳
*/
func (tt *TexTra) Translation(text string, lang_s string, lang_t string) (string, error) {
	token, err := tt.GetAccessToken()
	if err != nil {
		return "", err
	}

	values := &url.Values{}
	values.Add("access_token", token)
	values.Add("key", tt.Config.API_KEY)
	values.Add("name", tt.Config.UserName)
	values.Add("type", "json")
	values.Add("text", text)
	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/api/mt/generalNT_%s_%s/", tt.Config.BaseURL, lang_s, lang_t),
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	res := &api.TranslationResponse{}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return "", err
	}

	if res.ResultSet.Code != 0 {
		return "", fmt.Errorf("response error (%v) %v", res.ResultSet.Code, res.ResultSet.Message)
	}

	return res.ResultSet.Result.Text, nil
}
