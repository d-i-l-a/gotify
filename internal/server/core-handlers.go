package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gotify/internal/models"
	"gotify/internal/view"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

var clientId = os.Getenv("SPOTIFY_CLIENT_ID")
var clientSecret = os.Getenv("SPOTIFY_CLIENT_SECRET")

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type Profile struct {
	DisplayName  string `json:"display_name"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href      string `json:"href"`
	ID        string `json:"id"`
	Images    []any  `json:"images"`
	Type      string `json:"type"`
	URI       string `json:"uri"`
	Followers struct {
		Href  any `json:"href"`
		Total int `json:"total"`
	} `json:"followers"`
	Country         string `json:"country"`
	Product         string `json:"product"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`
	Email string `json:"email"`
}

func buildSpotifyURL() string {
	baseUrl := os.Getenv("BASEURL")
	return fmt.Sprintf(
		"https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=%s",
		clientId,
		baseUrl+"/callback",
		"streaming user-read-private user-read-email user-read-playback-state",
	)
}

func (s *Server) RefreshAccessToken(session *models.UserSession) error {
	data := setRefreshTokenQueryParams(session.RefreshToken)
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientId+":"+clientSecret)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("spotify: got %d status code: %s", resp.StatusCode, body)
	}
	var response RefreshTokenResponse
	err = json.Unmarshal([]byte(string(body)), &response)

	if err != nil {
		return err
	}
	session.AccessToken = response.AccessToken
	session.ExpiryTime = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	session.RefreshToken = response.RefreshToken
	s.db.UpdateUserSession(session)

	return nil
}

func (s *Server) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	baseUrl := os.Getenv("BASEURL")

	data := setAuthTokenQueryParams(code, baseUrl+"/callback")
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if len(clientId) < 1 || len(clientSecret) < 1 {
		http.Error(w, "ClientId or ClientSecret not provide in env", http.StatusInternalServerError)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(clientId+":"+clientSecret)))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Errorf(err.Error())
		}
	}()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if resp.StatusCode != 200 {
		http.Error(w, fmt.Sprintf("spotify: got %d status code: %s", resp.StatusCode, body), http.StatusInternalServerError)
	}
	var response TokenResponse
	err = json.Unmarshal([]byte(string(body)), &response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	userSession := models.UserSession{
		ID:           uuid.New().String(),
		Name:         "Not implemented",
		SessionID:    uuid.New().String(),
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
		ExpiryTime:   time.Now().Add(time.Duration(response.ExpiresIn) * time.Second),
	}

	s.db.CreateUserSession(&userSession)

	cookie := new(http.Cookie)
	cookie.Name = "session_id"
	cookie.Value = userSession.SessionID
	cookie.Expires = userSession.ExpiryTime
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func setAuthTokenQueryParams(authCode string, redirectURI string) url.Values {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("redirect_uri", redirectURI)
	return data
}

func setRefreshTokenQueryParams(refreshToken string) url.Values {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", clientId)
	return data
}

func MainHandler(w http.ResponseWriter, r *http.Request) {
	sessionCtxEntry := r.Context().Value(sessionContextKey)
	baseUrl := os.Getenv("BASEURL")
	if sessionCtxEntry == nil {
		view.Main(baseUrl, "", buildSpotifyURL()).Render(r.Context(), w)
		return
	}
	session, ok := sessionCtxEntry.(*models.UserSession)
	if !ok {
		panic("session in context is not a models.UserSession!")
	}
	view.Main(baseUrl, session.AccessToken, "").Render(r.Context(), w)
}
