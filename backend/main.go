package main

import (
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
)

func main() {
	app := chi.NewRouter()
	app.Use(
		middleware.Logger,
		middleware.RealIP,
		middleware.RedirectSlashes,
		middleware.StripSlashes,
		middleware.Recoverer,
	)

	// XXX may be i can create an interface for setup functions
	auth.Setup(app, db.DB)
	message.Setup(app)
	ws.Setup(app, db.DB)

	app.With(middlewares.Auth).Get("/", func(w http.ResponseWriter, r *http.Request) {
		count := st.Session.GetInt(r.Context(), "count")
		count++
		st.Session.Put(r.Context(), "count", count)

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(strconv.Itoa(count) + "\n"))
	})

	addr := config.GetListenAddr()
	log.Info("Server listening on" + addr)
	log.Fatal(http.ListenAndServe(addr, st.Session.LoadAndServeHeader(app)))
}
