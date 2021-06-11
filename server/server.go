package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"../connect"
	"../structures"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func StartServer() {
	connect.InitializeDatabase()
	defer connect.CloseConnection()
	// crear enrutador
	r := chi.NewRouter()
	r.Use(middleware.Timeout(10000 * time.Millisecond)) // 10 Seg
	r.Use(middleware.NoCache)                           // deshabilitar cache
	r.Use(middleware.Logger)                            // registros de solicitudes
	r.Use(render.SetContentType(render.ContentTypeJSON))
	// habilitar las conexiones desde cualquier sitio
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})
	r.Use(cors.Handler)
	r.Route("/users", func(r chi.Router) {
		r.Post("/", GetUser)
		r.Get("/best-scores", GetBestScores)
		r.Post("/{id}/score", NewScore)
	})

	log.Fatal(http.ListenAndServe(":8000", r))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	var response structures.ResponseUser
	userRequest, err := GetUserRequest(r)

	if err != nil {
		response = ServerResponseError("No se han podido leer los datos")
	} else {
		user := connect.GetUserByUsername(userRequest.Username)
		if !user.IsValid() {
			user = connect.CreateUser(userRequest)
			response = ServerResponseUserOk(user, "create and login success!.")
		} else {
			response = ServerResponseUserOk(user, "Login success!.")
		}
	}
	json.NewEncoder(w).Encode(response)
}

func NewScore(w http.ResponseWriter, r *http.Request) {
	scoreRequest, err := GetScoreRequest(r)
	var response structures.ResponseUser
	if err != nil {
		response = ServerResponseError("No se han podido leer los datos")
	} else {
		user := connect.NewScore(scoreRequest)
		response = ServerResponseUserOk(user, "puntaje guardado correctamente")
	}
	json.NewEncoder(w).Encode(response)
}

func GetUserRequest(r *http.Request) (structures.User, error) {
	var user structures.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return structures.User{}, errors.New("NO se pudo crear el usuario")
	}
	return user, nil
}

func GetScoreRequest(r *http.Request) (structures.Score, error) {
	var score structures.Score

	err := json.NewDecoder(r.Body).Decode(&score)
	if err != nil {
		return structures.Score{}, errors.New("NO se pudo crear el score")
	}
	return score, nil
}

func GetBestScores(w http.ResponseWriter, r *http.Request) {
	scores := connect.GetBestScores()
	response := ServerResponseScoresOk(scores, "Información optenido correctamente")
	json.NewEncoder(w).Encode(response)
}

func ServerResponseScoresOk(data []structures.BestScores, message string) structures.ResponseScores {
	return ServerResponseScores(
		http.StatusBadRequest,
		data,
		message,
	)
}

func ServerResponseUserOk(data structures.User, message string) structures.ResponseUser {
	return ServerResponseUser(
		http.StatusBadRequest,
		data,
		message,
	)
}

func ServerResponseError(message string) structures.ResponseUser {
	return ServerResponseUser(
		http.StatusOK,
		structures.User{},
		message,
	)
}

func ServerResponseUser(status int, data structures.User, message string) structures.ResponseUser {
	return structures.ResponseUser{Status: status, Data: data, Message: message}
}

func ServerResponseScores(status int, data []structures.BestScores, message string) structures.ResponseScores {
	return structures.ResponseScores{Status: status, Data: data, Message: message}
}
