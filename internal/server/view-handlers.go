package server

import (
	"go-spordlfy/internal/models"
	"go-spordlfy/internal/view"
	"net/http"
)

func SearchTracksHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(sessionContextKey).(*models.UserSession)
	if !ok {
		http.Error(w, "failed to get session info", http.StatusInternalServerError)
		return
	}

	searchTerm := r.FormValue("search")

	searchResponse, err := SearchTracks(session.AccessToken, searchTerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if searchResponse != nil {
		view.SearchTracks(*searchResponse).Render(r.Context(), w)
	} else {
		view.NoResults().Render(r.Context(), w)
	}
}

func SearchAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(sessionContextKey).(*models.UserSession)
	if !ok {
		http.Error(w, "failed to get session info", http.StatusInternalServerError)
	}

	searchTerm := r.FormValue("search")
	searchResponse, err := SearchAlbums(session.AccessToken, searchTerm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if searchResponse != nil {
		view.SearchAlbums(*searchResponse).Render(r.Context(), w)
	} else {
		view.NoResults().Render(r.Context(), w)
	}
}

func PlayListsHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(sessionContextKey).(*models.UserSession)
	if !ok {
		http.Error(w, "failed to get session info", http.StatusInternalServerError)
	}

	playLists, err := PlayLists(session.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.PlayLists(*playLists).Render(r.Context(), w)
}

func QueueHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(sessionContextKey).(*models.UserSession)
	if !ok {
		http.Error(w, "failed to get session info", http.StatusInternalServerError)
	}

	queue, err := Queue(session.AccessToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.Queue(queue).Render(r.Context(), w)
}

func AlbumInfoHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(sessionContextKey).(*models.UserSession)
	if !ok {
		http.Error(w, "failed to get session info", http.StatusInternalServerError)
	}

	albumID := r.PathValue("albumid")
	if albumID == "" {
		http.Error(w, "missing album id", http.StatusBadRequest)
		return
	}

	albumInfo, err := AlbumInfo(session.AccessToken, albumID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.AlbumInfo(albumInfo).Render(r.Context(), w)
}

func ArtistInfoHandler(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(sessionContextKey).(*models.UserSession)
	if !ok {
		http.Error(w, "failed to get session info", http.StatusInternalServerError)
	}

	artistId := r.PathValue("artistid")
	if artistId == "" {
		http.Error(w, "missing artist id", http.StatusBadRequest)
		return
	}

	artistInfo, err := ArtistInfo(session.AccessToken, artistId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	albumsOfArtist, err := AlbumsOfArtist(session.AccessToken, artistId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	topTracksOfArtist, err := TopTracksOfArtist(session.AccessToken, artistId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	view.ArtistInfo(artistInfo, albumsOfArtist, topTracksOfArtist).Render(r.Context(), w)
}
