package server

import (
	"context"
	"fmt"
	"go-spordlfy/internal/models"
	"go-spordlfy/internal/view"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//TODO: return 403 if no session found
		if r.URL.Path == "/callback" {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/static") {
			next.ServeHTTP(w, r)
			return
		}
		session, err := s.loadUserSession(r)
		if err != nil {
			baseUrl := os.Getenv("BASEURL")
			w.Header().Set("HX-Retarget", "html")
			view.Main(baseUrl, "", buildSpotifyURL()).Render(r.Context(), w)
			return
		}
		ctx := context.WithValue(r.Context(), "session", session)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (s *Server) loadUserSession(r *http.Request) (*models.UserSession, error) {
	sessionCookie, err := r.Cookie("session_id")
	fmt.Println("sessionCookie: ", sessionCookie)
	if err != nil {
		return nil, err
	}
	session, err := s.db.LoadSessionBySessionId(sessionCookie.Value)
	if err != nil {
		return nil, err
	}
	fmt.Println("session: ", session)
	if session.ExpiryTime.Add(1 * time.Minute).After(time.Now()) {
		err := s.RefreshAccessToken(session)
		if err != nil {
			return nil, err
		}
		return session, nil
	}
	return session, nil
}

func noCacheMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache")
		next.ServeHTTP(w, r)
	})
}
