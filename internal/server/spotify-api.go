package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-spordlfy/internal/models"
	"io"
	"net/http"
	"strings"
)

var client = &http.Client{}

func SearchTracks(accessToken string, searchTerm string) (*models.Tracks, error) {
	if searchTerm == "" {
		return nil, nil
	}
	urlEncodedSearchTerm := strings.ReplaceAll(searchTerm, " ", "%20")
	responseBody, err := makeHttpRequest(http.MethodGet, "https://api.spotify.com/v1/search?q="+urlEncodedSearchTerm+"&type=track", accessToken)
	if err != nil {
		return nil, err
	}
	var searchResponse models.SearchResponse
	json.Unmarshal(*responseBody, &searchResponse)
	return &searchResponse.Tracks, nil
}

func SearchAlbums(accessToken string, searchTerm string) (*models.Albums, error) {
	if searchTerm == "" {
		return nil, nil
	}
	urlEncodedSearchTerm := strings.ReplaceAll(searchTerm, " ", "%20")
	responseBody, err := makeHttpRequest(http.MethodGet, "https://api.spotify.com/v1/search?q="+urlEncodedSearchTerm+"&type=album", accessToken)
	if err != nil {
		return nil, err
	}
	var searchResponse models.SearchResponse
	json.Unmarshal(*responseBody, &searchResponse)
	return &searchResponse.Albums, nil
}

func PlayLists(accessToken string) (*models.PlayLists, error) {
	responseBody, err := makeHttpRequest(http.MethodGet, "https://api.spotify.com/v1/me/playlists", accessToken)
	if err != nil {
		return nil, err
	}
	var searchResponse models.PlayLists
	json.Unmarshal(*responseBody, &searchResponse)
	return &searchResponse, nil
}

func Queue(accessToken string) (*models.Queue, error) {
	responseBody, err := makeHttpRequest(http.MethodGet, "https://api.spotify.com/v1/me/player/queue", accessToken)
	if err != nil {
		return nil, err
	}

	var queueResponse models.Queue
	json.Unmarshal(*responseBody, &queueResponse)
	return &queueResponse, nil
}

func AlbumInfo(accessToken string, albumId string) (*models.Album, error) {
	responseBody, err := makeHttpRequest(http.MethodGet, "https://api.spotify.com/v1/albums/"+albumId, accessToken)
	if err != nil {
		return nil, err
	}
	var album models.Album
	json.Unmarshal(*responseBody, &album)
	return &album, nil
}

func ArtistInfo(accessToken string, artistId string) (*models.Artist, error) {
	responseBody, err := makeHttpRequest(http.MethodGet, "https://api.spotify.com/v1/artists/"+artistId, accessToken)
	if err != nil {
		return nil, err
	}
	var album models.Artist
	json.Unmarshal(*responseBody, &album)
	return &album, nil
}

func AlbumsOfArtist(accessToken string, artistId string) (*models.Albums, error) {
	responseBody, err := makeHttpRequest(http.MethodGet, "https://api.spotify.com/v1/artists/"+artistId+"/albums?limit=50", accessToken)
	if err != nil {
		return nil, err
	}
	var albums models.Albums
	json.Unmarshal(*responseBody, &albums)
	return &albums, nil
}

func TopTracksOfArtist(accessToken string, artistId string) (*models.TopTracksOfArtist, error) {
	responseBody, err := makeHttpRequest(http.MethodGet, "https://api.spotify.com/v1/artists/"+artistId+"/top-tracks?market=US", accessToken)
	if err != nil {
		return nil, err
	}
	var tracks models.TopTracksOfArtist
	json.Unmarshal(*responseBody, &tracks)
	return &tracks, nil
}

func makeHttpRequest(method string, url string, accessToken string) (*[]byte, error) {
	return makeHttpRequestWithBody(method, url, accessToken, nil)
}

func makeHttpRequestWithBody(method string, url string, accessToken string, body *bytes.Buffer) (*[]byte, error) {
	if body == nil {
		body = bytes.NewBuffer([]byte{})
	}
	req, err := http.NewRequest(method, url, body)
	fmt.Println("req", req)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("spotify return invalid status code " + resp.Status + " ( " + string(responseBody) + ")")
	}
	return &responseBody, nil
}
