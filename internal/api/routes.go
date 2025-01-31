package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thomasjinlo/gochatter/internal/models"
)

type UserRepository interface {
	GetAll() []*models.User
	Login(string) error
}

func setupRoutes(pc *PushClient, ur UserRepository) *chi.Mux {
	r := chi.NewRouter()
	r.Use(loggingMiddleWare)
	r.Get("/hello", handleHello())
	r.Get("/users", handleGetUsers(ur))
	r.Post("/login", handleLogin(ur))
	r.Post("/direct_message", handleDirectMessage(pc))
	return r
}

type GetUserResponse struct {
	Id   string
	Name string
}

type GetUsersResponse struct {
	Users []GetUserResponse
}

func handleGetUsers(ur UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			e := fmt.Sprintf("unsupported content type: %s", contentType)
			log.Print(e)
			http.Error(w, e, http.StatusUnsupportedMediaType)
			return
		}
		getUsersRes := GetUsersResponse{}
		for _, u := range ur.GetAll() {
			getUsersRes.Users = append(
				getUsersRes.Users,
				GetUserResponse{
					Id:   u.Username(),
					Name: u.Username(),
				},
			)
		}
		log.Printf("[gochatter-api] users %v", getUsersRes.Users)
		res, _ := json.Marshal(getUsersRes)
		w.Write(res)
		w.WriteHeader(http.StatusOK)
	}
}

type DirectMessageRequest struct {
	SourceAccountId string
	TargetAccountId string
	Content         string
}

func handleDirectMessage(pc *PushClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if contentType := r.Header.Get("Content-Type"); contentType != "application/json" {
			e := fmt.Sprintf("unsupported content type: %s", contentType)
			log.Print(e)
			http.Error(w, e, http.StatusUnsupportedMediaType)
			return
		}
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var msg DirectMessageRequest
		err = json.Unmarshal(body, &msg)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = pc.SendDirectMessage(msg.TargetAccountId, msg.SourceAccountId, msg.Content)
		if err != nil {
			log.Printf("[gochatter-api] error from redis: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type LoginRequest struct {
	Username string
}

func handleLogin(ur UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var lr LoginRequest
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &lr); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := ur.Login(lr.Username); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
		w.WriteHeader(http.StatusOK)
	}
}

func loggingMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[gochatter-api] received HTTP method %s on path %s", r.Method, r.URL.String())
		next.ServeHTTP(w, r)
	})
}
