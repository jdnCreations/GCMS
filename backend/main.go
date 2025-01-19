package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/jdnCreations/gcms/internal/auth"
	"github.com/jdnCreations/gcms/internal/database"
	"github.com/jdnCreations/gcms/internal/models"
	"github.com/jdnCreations/gcms/internal/utils"
)

type apiConfig struct {
  db *database.Queries
  secret string
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}


func respondWithError(w http.ResponseWriter, code int, msg string) {
  type returnErr struct {
    Error string `json:"error"`
  }

  respBody := returnErr{
    Error: msg,
  }
  dat, err := json.Marshal(respBody)
  if err != nil {
    log.Printf("Error marshalling JSON %s", err)
    w.WriteHeader(500)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
  dat, err := json.Marshal(payload)
  if err != nil {
    respondWithError(w, 500, "Error marshalling JSON")
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(dat)
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type User struct {
		FirstName string;
		LastName string;
		Email string;
    Password string;
	}
	log.Println("Attempting to create a user")
	decoder := json.NewDecoder(r.Body)
	params := User{} 
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

  if params.Password == "" {
    respondWithError(w, 422, "Password is required")
    return
  }

  pw, err := auth.HashPassword(params.Password)
  if err != nil {
    log.Printf("Could not hash password: %s", err.Error())
    respondWithError(w, 422, "Could not hash password")
    return
  }


	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		FirstName: params.FirstName,
		LastName: params.LastName,
		Email: params.Email,
    HashedPassword: pw,
	}) 
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		respondWithError(w, 422, "Could not create user")
		return
	}

  type UserResponse struct {
    ID string;
    FirstName string;
    LastName string;
    Email string;
  }

	respondWithJSON(w, 201, UserResponse{
    ID: user.ID.String(),
    FirstName: user.FirstName,
    LastName: user.LastName,
    Email: user.Email,
  })
}

func (cfg *apiConfig) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	refresh, err := r.Cookie("refresh_token")
	if err != nil {
		log.Printf("No cookie found: %v", err)
	} else {
		log.Printf("Found cookie with token: %v", refresh.Value)
		token, err := cfg.db.RevokeToken(r.Context(), refresh.Value)
		if err != nil {
			log.Printf("Could not revoke token: %v", token)
		} else {
			log.Printf("Successfully revoked token: %v", token)
		}	
	}

	cookie := http.Cookie{
		Name: "refresh_token", 
		Value: "", 
		Path: "/",
		MaxAge: -1,
		HttpOnly: true, 
		SameSite: http.SameSiteStrictMode,
	}
  http.SetCookie(w, &cookie)

	respondWithJSON(w, 201, "");
}

func (cfg *apiConfig) handleLoginUser(w http.ResponseWriter, r *http.Request) {
  type ExpectedBody struct {
    Password string;
    Email string;
  }

  decoder := json.NewDecoder(r.Body)
	params := ExpectedBody{} 
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, err.Error())
		return
	}

  if params.Password == "" {
    respondWithError(w, 422, "Password is required")
    return
  }

  if params.Email == "" {
    respondWithError(w, 422, "Email is required")
    return
  }

  user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
  if err != nil {
    respondWithError(w, 404, "user not found")
    return
  }

  err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
  if err != nil {
    respondWithError(w, 401, "invalid login details")
    return
  }

  jwt, err := auth.MakeJWT(user.ID, cfg.secret)
  if err != nil {
    respondWithError(w, 500, "could not create jwt token")
    return
  }

	refresh, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, "could not create refresh token")
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refresh,
		UserID: user.ID,
		ExpiresAt: pgtype.Timestamp{Time: time.Now().Add(60 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		log.Println("error saving token to db", err.Error())
		respondWithError(w, 500, "could not save token to db")
		return
	}

  cookie := http.Cookie{
		Name: "refresh_token", 
		Value: refresh, 
		Path: "/", 
		HttpOnly: true, 
		MaxAge: 60 * 24 * 60 * 60, 
		SameSite: http.SameSiteStrictMode,
	}
  http.SetCookie(w, &cookie)

  type UserResponse struct {
    ID string
    Email string
		FirstName string
    Token string
		RefreshToken string
    IsAdmin bool
  }

  respondWithJSON(w, 200, UserResponse{ID: user.ID.String(), Email: user.Email, FirstName: user.FirstName, Token: jwt, RefreshToken: refresh, IsAdmin: user.IsAdmin})
}

func (cfg *apiConfig) handleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to retrieve all users")
	users, err := cfg.db.GetAllUsers(context.Background())
	if err != nil {
		log.Printf("Error: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, "Could not retrieve users")
		return
	}

	respondWithJSON(w, http.StatusOK, users)	
}

func (cfg *apiConfig) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, 401, "invalid bearer token")
    return
  }

  userId, err := auth.ValidateJWT(token, cfg.secret)
  if err != nil {
    respondWithError(w, 401, "invalid jwt")
    return
  }

	log.Println("Attempting to delete user")
	vars := mux.Vars(r)
	id := vars["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid user ID format")
		return
	}

  if userId != uuid {
    respondWithError(w, 401, "you cannot delete another user")
    return
  }

	err = cfg.db.DeleteUserById(r.Context(), uuid) 
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		respondWithError(w, http.StatusBadRequest, "Unable to delete user")
		return
	}

	respondWithJSON(w, http.StatusOK, "Deleted user")
}

func (cfg *apiConfig) handleGetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid user ID format")
		return
	}

	user, err := cfg.db.GetUserById(context.Background(), uuid)
	if err != nil {
		respondWithError(w, 404, "user not found")
		return
	}

	respondWithJSON(w, 200, user)
}

func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, 401, "invalid bearer token")
    return
  }

  userId, err := auth.ValidateJWT(token, cfg.secret)
  if err != nil {
    respondWithError(w, 401, "invalid jwt")
    return
  }


	vars := mux.Vars(r)
	id := vars["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid user ID format")
		return
	}

  if userId != uuid {
    respondWithError(w, 401, "cannot update another user")
    return
  }

	type User struct {
		FirstName string
		LastName string
		Email string
	}

	decoder := json.NewDecoder(r.Body)
	params := User{} 
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid input data")
		return
	}

	cust, err := cfg.db.UpdateUser(context.Background(), database.UpdateUserParams{
		Column1: params.FirstName,
		Column2: params.LastName,
		Column3: params.Email,
		ID: uuid,
	})
	if err != nil {
		respondWithError(w, 500, "Could not update user")
		return
	}

	respondWithJSON(w, http.StatusOK, cust)
}

func (cfg *apiConfig) handleGetAllGames(w http.ResponseWriter, r *http.Request) {
	games, err := cfg.db.GetAllGames(r.Context())
	if err != nil {
		log.Fatal(err)
		respondWithError(w, 500, "Could not retrieve all games")
	}

	respondWithJSON(w, 200, games)
}

func (cfg *apiConfig) handleCreateGame(w http.ResponseWriter, r *http.Request) {
  log.Println("Attempting to create game")
  decoder := json.NewDecoder(r.Body)
  gameInfo := models.GameInfo{}
	err := decoder.Decode(&gameInfo)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// validate game 
	validate := validator.New()
	err = validate.Struct(gameInfo)
	if err != nil {
		respondWithError(w, 422, err.Error())
		return
	}

	game, err := cfg.db.CreateGame(r.Context(), 
	database.CreateGameParams{
		Title: gameInfo.Title,
		Copies: gameInfo.Copies,
	})
	if err != nil {
		respondWithError(w, 422, "Could not create game, name must be unique")
		return
	}

  respondWithJSON(w, 200, game)
}

func (cfg *apiConfig) handleDeleteGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid game ID format")
		return
	}

	err = cfg.db.DeleteGameById(r.Context(), uuid)
	if err != nil {
		respondWithError(w, 404, "Game not found")
		return
	}

	respondWithJSON(w, 200, nil)
}

func (cfg *apiConfig) handleGetGameById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid game ID format")
		return
	}

	game, err := cfg.db.GetGameById(context.Background(), uuid)
	if err != nil {
		respondWithError(w, 404, "Game not found")
		return
	}

	respondWithJSON(w, 200, game)
}

func (cfg *apiConfig) handleUpdateGame(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	uuid, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, 400, "Invalid game ID format")
		return
	}

	gameToUpdate, err := cfg.db.GetGameById(r.Context(), uuid)
	if err != nil {
		respondWithError(w, 404, "Game does not exist with that id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := models.UpdateGameInfo{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

	if params.Title != nil {
		gameToUpdate.Title = *params.Title
	}

	if params.Copies != nil {
		gameToUpdate.Copies = *params.Copies
	}
	
	game, err := cfg.db.UpdateGame(r.Context(), database.UpdateGameParams{
		Column1: gameToUpdate.Title,
		Column2: gameToUpdate.Copies,
		ID: gameToUpdate.ID,
	})
	if err != nil {
		respondWithError(w, 422, "Could not update game")
		return
	}

	respondWithJSON(w, 200, game)
}

func (cfg *apiConfig) handleCreateGenre(w http.ResponseWriter, r *http.Request) {
  decoder := json.NewDecoder(r.Body)
  genreInfo := models.GenreInfo{}
  err := decoder.Decode(&genreInfo)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, err.Error())
		return
  }

  genre, err := cfg.db.CreateGenre(r.Context(), genreInfo.Name)
  if err != nil {
    respondWithError(w, 422, "Genre may already exist")
    return
  }

  respondWithJSON(w, 200, genre)
}

func (cfg *apiConfig) handleGetAllGenres(w http.ResponseWriter, r *http.Request) {
  genres, err := cfg.db.GetAllGenres(r.Context())
  if err != nil {
    respondWithError(w, 422, "Could not retrieve genres")
    return
  }

  respondWithJSON(w, 200, genres)
}

func (cfg *apiConfig) handleGetGenreById(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id := vars["id"]
  uuid, err := uuid.Parse(id)
  if err != nil {
    respondWithError(w, 400, "Invalid ID format")
    return
  }

  genre, err := cfg.db.GetGenreById(r.Context(), uuid)
  if err != nil {
    respondWithError(w, 404, "No genre with that id")
    return
  }

  respondWithJSON(w, 200, genre)
}

func (cfg *apiConfig) handleDeleteGenre(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id := vars["id"]
  uuid, err := uuid.Parse(id)
  if err != nil {
    respondWithError(w, 400, "invalid ID format")
    return
  }

  err = cfg.db.DeleteGenreById(r.Context(), uuid)
  if err != nil {
    respondWithError(w, 404, "No genre with that id")
    return
  }

  respondWithJSON(w, 200, "deleted")
}

func (cfg *apiConfig) handleUpdateGenre(w http.ResponseWriter, r *http.Request) {
  uuid, err := utils.GetIdFromRequest("id", r)
  if err != nil {
    respondWithError(w, 400, "Invalid ID format")
  }

  genreToUpdate, err := cfg.db.GetGenreById(r.Context(), uuid)
  if err != nil {
    respondWithError(w, 404, "No genre with that id")
    return
  }

  type p struct {
    Name *string
  }
  params := p{}
  decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}

  if params.Name != nil {
    genreToUpdate.Name = *params.Name
  }
  
  genre, err := cfg.db.UpdateGenre(r.Context(), database.UpdateGenreParams{
    Column1: genreToUpdate.Name,
    ID: genreToUpdate.ID,
  })
  if err != nil {
    respondWithError(w, 422, "Unable to update genre")
    return
  }

  respondWithJSON(w, 200, genre)
}

func (cfg *apiConfig) handleAddGenreToGame(w http.ResponseWriter, r *http.Request) {
  gameID, err := utils.GetIdFromRequest("gameID", r)
  if err != nil {
    respondWithError(w, 400, "Invalid ID format - gameid")
    return
  }

  type p struct {
    GenreID string `json:"genre_id"`
  }
  params := p{}
  decoder := json.NewDecoder(r.Body)
  err = decoder.Decode(&params)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "invalid request payload")
    return
  }

  genreID, err := uuid.Parse(params.GenreID)
  if err != nil {
    respondWithError(w, 400, "Invalid ID format - genreID")
    return
  }

  gg, err := cfg.db.AddGenreToGame(r.Context(), database.AddGenreToGameParams{
    GameID: gameID,
    GenreID: genreID,
  })
  if err != nil {
    respondWithError(w, 422, "Could not add genre to game")
    return
  }

  respondWithJSON(w, 200, gg)
}

func (cfg *apiConfig) handleActiveReservations(w http.ResponseWriter, r *http.Request) {
	res, err := cfg.db.GetAllActiveReservations(r.Context())
  if err != nil {
    respondWithError(w, 422, "Could not retrieve active reservations")
    return
  }

  respondWithJSON(w, 200, res)
}

func (cfg *apiConfig) handleGetAllReservations(w http.ResponseWriter, r *http.Request) {
	res, err := cfg.db.GetAllReservations(r.Context())
  if err != nil {
    respondWithError(w, 422, "Could not retrieve reservations")
    return
  }

  respondWithJSON(w, 200, res)
}

func (cfg *apiConfig) handleCreateReservation(w http.ResponseWriter, r *http.Request) {
	// check user id matches logged in user details
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, 404, "invalid bearer token")
    return
  }

  uuid, err := auth.ValidateJWT(token, cfg.secret)
  if err != nil {
    respondWithError(w, 404, "invalid jwt token")
    return
  }


  type p struct {
    Date string
    StartTime string
    EndTime string
    UserID string
    GameID string
  }
  params := p{}
  decoder := json.NewDecoder(r.Body)
  err = decoder.Decode(&params)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "invalid request payload")
    return
  }

  // make sure user ids match
  if params.UserID != uuid.String() {
    respondWithError(w, 403, "You cannot create a reservation for another user")
    return
  }
  
	// convert to appropriate types for db ?
  resDate, err := utils.ConvertToPGDate(params.Date)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "invalid reservation date")
    return
  }
  
  
	startTime, err := utils.ConvertToPGTime(params.StartTime)
	if err != nil {
    respondWithError(w, http.StatusBadRequest, "invalid start time")
		return
	}
  
  endTime, err := utils.ConvertToPGTime(params.EndTime)
	if err != nil {
    respondWithError(w, http.StatusBadRequest, "invalid end time")
		return
	}
	
	userId, err := utils.ConvertToPGUUID(params.UserID)
	if err != nil {
    respondWithError(w, http.StatusBadRequest, "invalid userid")
		return 
	}
	
	gameId, err := utils.ConvertToPGUUID(params.GameID)
	if err != nil {
    respondWithError(w, http.StatusBadRequest, "invalid gameid")
		return 
	}

  // check if game is available
  copies, err := cfg.db.GetCurrentCopies(r.Context(), gameId.Bytes)
  if err != nil {
    respondWithError(w, 422, "invalid game id")
    return
  }

  if copies < 1 {
    respondWithError(w, 422, "No copies avaiable")
    return
  }
	
  res, err := cfg.db.CreateReservation(r.Context(), database.CreateReservationParams{
    ResDate: resDate,
    StartTime: startTime,
		EndTime: endTime,
		UserID: userId,
		GameID: gameId,
  })
	if err != nil {
		log.Printf("Error: %s", err.Error())
		respondWithError(w, 422, "Could not create reservation")
		return
	}

  err = cfg.db.DecCurrentCopies(r.Context(), gameId.Bytes)
  if err != nil {
    log.Printf("Could not decrement current copies: %s", err.Error())
  }

	respondWithJSON(w, 201, res)
}

func (cfg *apiConfig) handleGetReservationsForUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	uuid, err := utils.ConvertToPGUUID(userID) 
  if err != nil {
    respondWithError(w, 400, "Invalid ID format")
    return
  }

	res, err := cfg.db.GetReservationsForUser(r.Context(), uuid)
	if err != nil {
    respondWithError(w, 422, "Could not retrieve reservations")
    return
	}

  respondWithJSON(w, 200, res)
}

func (cfg *apiConfig) handleCheckGameAvailable(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	gameId := vars["gameID"]

	uuid, err := utils.ConvertToPGUUID(gameId)
	if err != nil {
		respondWithError(w, 422, "Invalid ID format")
		return
	}

	type p struct {
    StartTime string
    EndTime string
    UserID string
  }
  params := p{}
  decoder := json.NewDecoder(r.Body)
  err = decoder.Decode(&params)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "invalid request payload")
    return
  }

	pgStartTime, err := utils.ConvertToPGTime(params.StartTime);
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid start time")
		return
	}

	pgEndTime, err := utils.ConvertToPGTime(params.EndTime);
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid end time")
		return
	}

	num, err := cfg.db.CheckGameReservation(r.Context(), database.CheckGameReservationParams{
		GameID: uuid,
		StartTime: pgStartTime,
		EndTime: pgEndTime,
	})
	if err != nil {
		respondWithError(w, 422, "not enough copies")
		return
	}

	if num > 0 {
		respondWithJSON(w, 200, "copy is available")
	}
	respondWithError(w, 422, "not enough copies")
}

func (cfg *apiConfig) handleDeleteReservation(w http.ResponseWriter, r *http.Request) {
  // authorize user
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, 401, "invalid bearer token")
    return
  }

  id, err := auth.ValidateJWT(token, cfg.secret)
  if err != nil {
    respondWithError(w, 401, "invalid jwt token")
    return
  }

	vars := mux.Vars(r)
	resID:= vars["id"]

	resId, err := uuid.Parse(resID)
  if err != nil {
    respondWithError(w, 400, "Invalid ID format")
    return
  }

  res, err := cfg.db.GetReservationById(r.Context(), resId)
  if err != nil {
    respondWithError(w, 404, "invalid reservation id")
    return
  }

  // check if reservation user id matches user id in bearer
  if id != res.UserID.Bytes {
    respondWithError(w, 403, "You cannot delete another user's reservation")
    return
  }

	err = cfg.db.DeleteReservation(r.Context(), resId)
	if err != nil {
		respondWithError(w, 422, "Could not delete reservation")
		return
	}

  _, err = cfg.db.IncCurrentCopies(r.Context(), res.GameID.Bytes)
  if err != nil {
    log.Printf("Could not increment copies for game: %v", resId)
  }

	respondWithJSON(w, 200, "")
}

func (cfg *apiConfig) handleGetCurrentCopies(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
	resID:= vars["id"]

	uuid, err := uuid.Parse(resID)
  if err != nil {
    respondWithError(w, 400, "Invalid ID format")
    return
  }
  copies, err := cfg.db.GetCurrentCopies(r.Context(), uuid)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, "could not get current copies")
    return
  }

  respondWithJSON(w, 200, copies)
}

func (cfg *apiConfig) runPeriodicReservationChecker(interval time.Duration) {
  log.Println("running periodic reservation checker...")
  ticker := time.NewTicker(interval)
  defer ticker.Stop()

  for range ticker.C {
    cfg.checkReservationsAndUpdate()
  }
}

func (cfg *apiConfig) checkReservationsAndUpdate() {
  // check for any expired reservations, if expired change active to false, increment that games current_copies
  expired, err := cfg.db.GetExpiredReservations(context.Background())
  if err != nil {
    log.Printf("cannot get expired reservations: %s", err.Error())
  }

  for _, res := range expired {
    // set active to false
    _, err := cfg.db.SetReservationInactive(context.Background(), res.ID)
    if err != nil {
      log.Printf("invalid id, could not set reservation to inactive: %s", err.Error())
    }

    // increment game count for the game
    _, err = cfg.db.IncCurrentCopies(context.Background(), res.GameID.Bytes)
    if err != nil {
      log.Printf("could not increment copies for reservation: %v, err: %s", res.ID, err.Error())
    }
    log.Printf("set reservation: %s to inactive", res.ID)
  }
}

// func (cfg *apiConfig) handleSetAdmin(w http.ResponseWriter, r *http.Request) {
// 	log.Println("HIT DA ADMIN ENDPOINT")
// }

func (cfg *apiConfig) handleSetAdmin(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Could not extract bearer")
		return
	}

	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "invalid token")
		return
	}

	user, err := cfg.db.GetUserById(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, "user does not exist")
		return
	}

	if !user.IsAdmin {
		respondWithError(w, 404, "invalid permissions")
		return
	}

	type p struct {
		SetAdmin bool 
    UserID string
  }
  params := p{}
  decoder := json.NewDecoder(r.Body)
  err = decoder.Decode(&params)	
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
    return
	}

	log.Printf("%s is trying to set admin on: %s", user.ID, params.UserID)

	userId, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, 400, "Invalid ID format")
    return
	}

	err = cfg.db.SetAdmin(r.Context(), database.SetAdminParams{
		IsAdmin: params.SetAdmin,
		ID: userId,
	})
	if err != nil {
		respondWithError(w, 422, "could not set admin")
		return
	}

	respondWithJSON(w, 200, "admin set")
}

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

  token, err := cfg.db.GetRefreshToken(r.Context(), cookie.Value)
	if err != nil {
		respondWithError(w, 401, "invalid refresh token")
		return
	}

	if time.Now().After(token.ExpiresAt.Time) {
		respondWithError(w, 401, "invalid refresh token")
		return
	}

	if token.RevokedAt.Valid {
		respondWithError(w, 401, "invalid refresh token")
		return
	}

  user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token.Token)
  if err != nil {
    respondWithError(w, 401, "no user exists with refresh token")
    return
  }

  newToken, err := auth.MakeJWT(user.ID, cfg.secret)
  if err != nil {
    respondWithError(w, 500, "could not create jwt token")
    return
  }

	type res struct {
		Token string
		Name string
		IsAdmin bool
		Email string
		ID string
	}

	respondWithJSON(w, 200, res{
		Token: newToken,
		Name: user.FirstName,
		Email: user.Email,
		IsAdmin: user.IsAdmin,
		ID: user.ID.String(),
	});
}


func (cfg *apiConfig) handleRevokeToken(w http.ResponseWriter, r *http.Request) {
  bearer, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, 401, "invalid bearer token")
    return
  }

  _, err = cfg.db.RevokeToken(r.Context(), bearer)
  if err != nil {
    respondWithError(w, 500, "could not revoke token")
    return
  }

  respondWithJSON(w, 204, "")
}

func (cfg *apiConfig) handleVerifyToken(w http.ResponseWriter, r *http.Request) {
  token, err := auth.GetBearerToken(r.Header)
  if err != nil {
    respondWithError(w, 401, "invalid bearer token")
    return
  }

  userId, err := auth.ValidateJWT(token, cfg.secret)
  if err != nil {
    respondWithError(w, 401, "invalid access token")
    return
  }

  user, err := cfg.db.GetUserById(r.Context(), userId)
  if err != nil {
    respondWithError(w, 404, "no user with that id found")
    return
  }

  respondWithJSON(w, 200, user)
}


func main() {
    err := godotenv.Load(".env")
    if err != nil {
      log.Fatalf("Could not load .env: %v", err)
    }
		port := os.Getenv("PORT")
		dbURL := os.Getenv("DATABASE_URL")
    secret := os.Getenv("SECRET")
		if port == "" {
			log.Fatal("PORT environment variable is not set")
		}

		dbpool, err := pgxpool.New(context.Background(), dbURL)
		if err != nil {
			log.Fatal(err)
		}
		defer dbpool.Close()

		dbQueries := database.New(dbpool)
		apiCfg := apiConfig{}
		apiCfg.db = dbQueries
    apiCfg.secret = secret

		r := mux.NewRouter()	
		// wrap with cors and json content type
		corsRouter := enableCORS(jsonContentTypeMiddleware(r))
		server := &http.Server{
			Addr: ":8080",
			Handler: corsRouter,
		}
		r.HandleFunc("/api/healthz", handleReadiness).Methods("GET")

		// refresh token
		r.HandleFunc("/api/refresh", apiCfg.handleRefreshToken).Methods("POST");
    r.HandleFunc("/api/revoke", apiCfg.handleRevokeToken).Methods("POST");
    r.HandleFunc("/api/verify", apiCfg.handleVerifyToken).Methods("GET");

		// users
		r.HandleFunc("/api/users/admin", apiCfg.handleSetAdmin).Methods("PUT")
		r.HandleFunc("/api/users", apiCfg.handleCreateUser).Methods("POST")
    r.HandleFunc("/api/users/login", apiCfg.handleLoginUser).Methods("POST")
		r.HandleFunc("/api/users/logout", apiCfg.handleLogoutUser).Methods("POST")
		r.HandleFunc("/api/users/{id}", apiCfg.handleGetUserById).Methods("GET")
		r.HandleFunc("/api/users", apiCfg.handleGetAllUsers).Methods("GET")
    r.HandleFunc("/api/users/{id}", apiCfg.handleUpdateUser).Methods("PUT")
		r.HandleFunc("/api/users/{id}", apiCfg.handleDeleteUser).Methods("DELETE")

		// games
		r.HandleFunc("/api/games", apiCfg.handleGetAllGames).Methods("GET")
		r.HandleFunc("/api/games", apiCfg.handleCreateGame).Methods("POST")
		r.HandleFunc("/api/games/{id}", apiCfg.handleGetGameById).Methods("GET")
		r.HandleFunc("/api/games/{id}", apiCfg.handleDeleteGame).Methods("DELETE")
		r.HandleFunc("/api/games/{id}", apiCfg.handleUpdateGame).Methods("PUT")
    r.HandleFunc("/api/games/{id}/copies", apiCfg.handleGetCurrentCopies).Methods("GET")

    // genres
    r.HandleFunc("/api/genres", apiCfg.handleCreateGenre).Methods("POST")
    r.HandleFunc("/api/genres", apiCfg.handleGetAllGenres).Methods("GET")
    r.HandleFunc("/api/genres/{id}", apiCfg.handleGetGenreById).Methods("GET")
    r.HandleFunc("/api/genres/{id}", apiCfg.handleDeleteGenre).Methods("DELETE")
    r.HandleFunc("/api/genres/{id}", apiCfg.handleUpdateGenre).Methods("PUT")

    // add genre to games
    r.HandleFunc("/api/game_genres/{gameID}", apiCfg.handleAddGenreToGame).Methods("POST")

		// reservations
    r.HandleFunc("/api/reservations", apiCfg.handleCreateReservation).Methods("POST")
		r.HandleFunc("/api/reservations", apiCfg.handleGetAllReservations).Methods("GET")
		r.HandleFunc("/api/reservations/active", apiCfg.handleActiveReservations).Methods("GET")
		r.HandleFunc("/api/reservations/check/{gameID}", apiCfg.handleCheckGameAvailable).Methods("GET")
		r.HandleFunc("/api/reservations/{id}", apiCfg.handleDeleteReservation).Methods("DELETE")
		r.HandleFunc("/api/reservations/{userID}", apiCfg.handleGetReservationsForUser).Methods("GET")

    // create an admin user
    hashPass, err := auth.HashPassword("admin")
    if err != nil {
      log.Fatalf("could not hash pw")
    }
    user, err := apiCfg.db.CreateUser(context.Background(), 
      database.CreateUserParams{
        FirstName: "admin",
        LastName: "admin",
        Email: "admin@mail.com",
        HashedPassword: hashPass,
      })
    if err != nil {
      log.Println("Could not create admin user, already exists")
    }

    // set admin
    apiCfg.db.SetAdmin(context.Background(), database.SetAdminParams{
      IsAdmin: true,
      ID: user.ID,
    })
    if err != nil {
      log.Printf("Could not set user as admin: %v", err)
    }
    go apiCfg.runPeriodicReservationChecker(5 * time.Minute)

		server.ListenAndServe()
		
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // allow any origin
    w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}