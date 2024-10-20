package message

import (
	"go-chat/middlewares"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(c *chi.Mux, db *mongo.Database) {
	repo := NewMessageRepo(db)
	handler := NewMessageHandler(repo)

	c.Route("/messages", func(r chi.Router) {
		r.Use(middlewares.Auth)

		// TODO implement get message queries like sender=x, receiver=x
		r.Get("/{userID}", handler.GetUserMessages)
		r.Get("/groups/{id}", handler.GetGroupMessages)
	})
}
