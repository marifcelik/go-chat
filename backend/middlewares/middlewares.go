package middlewares

import (
	"context"
	"net/http"

	"go-chat/storage"
	"go-chat/utils"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := storage.Session.GetString(r.Context(), "user")

		if user == "" {
			log.WithPrefix("AUTH MW").Warn("unauthorized request", "header", r.Header)
			utils.ErrResp(w, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// i can't use this middleware because if i do i have to parse the body twice
// TODO find a way to parse the body once or look for what other people do
func LoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := storage.Session.GetString(r.Context(), "user")
		log.Debug("user: %v\n", user)
		if user != "" {
			storage.Session.Put(r.Context(), "warn", storage.Session.GetInt(r.Context(), "warn")+1)
			utils.JsonResp(w, utils.M{"warn": "you already logged in"}, http.StatusConflict)
			return
		}
		// log.Warn("loggedin middleware unimplemented")
		next.ServeHTTP(w, r)
	})
}

func WsHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			utils.ErrResp(w, http.StatusUnauthorized)
			return
		}
		r.Header.Add("X-Session", token)
		ctx := context.WithValue(context.Background(), chi.RouteCtxKey, r.Context().Value(chi.RouteCtxKey))
		nc, err := storage.Session.Load(ctx, token)
		if err != nil {
			log.Error("session load", "err", err)
		}
		sr := r.WithContext(nc)
		next.ServeHTTP(w, sr)
	})
}

func UpgradeChecher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !(utils.ContainsI(r.Header.Get("connection"), "upgrade") && utils.ContainsI(r.Header.Get("upgrade"), "websocket")) {
			utils.ErrResp(w, http.StatusUpgradeRequired)
			return
		}
		next.ServeHTTP(w, r)
	})
}
