package auth

import (
	"go-chat/middlewares"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(c *chi.Mux, db *mongo.Database) {
	repo := NewAuthRepo(db)
	handler := NewAuthHandler(repo)

	c.Route("/auth", func(r chi.Router) {
		r.With(middlewares.LoggedIn).Post("/login", handler.Login)
		r.Post("/logout", handler.Logout)
		r.Post("/register", handler.Register)
	})
}
