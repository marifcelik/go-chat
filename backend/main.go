package main

import (
	"net"
	"net/http"
	"strconv"

	"go-chat/config"
	"go-chat/db"
	"go-chat/internal/auth"
	"go-chat/internal/message"
	"go-chat/internal/ws"
	"go-chat/middlewares"
	st "go-chat/storage"

	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	app := chi.NewRouter()
	app.Use(
		middleware.Logger,
		middleware.RealIP,
		middleware.RedirectSlashes,
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.Heartbeat("/healthz"),
		cors.Handler(cors.Options{
			AllowedMethods:   []string{"DELETE", "GET", "OPTIONS", "PATCH", "POST", "PUT", "UPDATE"},
			ExposedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", config.C.HeaderKey.Session, config.C.HeaderKey.Expiry},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", config.C.HeaderKey.Session, config.C.HeaderKey.Expiry},
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		}),
	)

	// XXX may be i can create an interface for setup functions
	auth.Setup(app, db.DB)
	message.Setup(app, db.DB)
	ws.Setup(app, db.DB)

	app.With(middlewares.Auth).Get("/", func(w http.ResponseWriter, r *http.Request) {
		count := st.Session.GetInt(r.Context(), "count")
		count++
		st.Session.Put(r.Context(), "count", count)

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(strconv.Itoa(count) + "\n"))
	})

	addr := net.JoinHostPort(config.C.Host, config.C.Port)
	log.Info("Server listening on " + addr)
	log.Fatal(http.ListenAndServe(addr, st.Session.LoadAndServeHeader(app)))
}
