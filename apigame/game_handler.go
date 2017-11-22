package apigame

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Bobsar0/Guess-A-Letter/tools"

	"github.com/go-chi/chi"
)

// GameHandler represents an HTTP API handler for games.
type GameHandler struct {
	mux *chi.Mux
	//redirectURL string
	Session *Session
	//game        *Game
	Logger *log.Logger
}

// NewGameHandler returns a new instance of GameHandler.
//mv GameService
func NewGameHandler(s *Session) *GameHandler {
	h := &GameHandler{
		mux:     chi.NewRouter(),
		Logger:  log.New(os.Stderr, "", log.LstdFlags),
		Session: s,
	}

	h.mux.Get("/games", s.Validate(h.listGamesHandler))
	h.mux.Get("/games/start", s.Validate(h.StartHandler))
	h.mux.Post("/games/guessletter", s.Validate(h.GuessLetterHandler))
	h.mux.Post("/games/guess", s.Validate(h.GuessHandler))
	h.mux.Get("/games/save/:id", s.Validate(h.SaveHandler))
	//h.mux.Get("/games/page", s.Validate(h.PageHandler))
	//h.mux.Post("/games/test", s.Validate(h.test))
	// h.mux.Post("/games", s.Validate(h.createHandler))
	// h.mux.Get("/games/:id", s.Validate(h.detailHandler))
	// h.mux.Post("/games/:id", s.Validate(h.updateHandler))
	// h.mux.Post("/games/delete/:id", s.Validate(h.deleteHandler))
	// h.mux.Get("/games/edit/:id", s.Validate(h.editFormHandler))

	return h
}

func (h *GameHandler) listGamesHandler(w http.ResponseWriter, r *http.Request) {
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		tools.Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}
	games, err := h.Session.gameService.ListGames()
	if err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	dat := struct {
		Games   []*Game
		UserGui *GameGuiService
	}{
		Games:   games,
		UserGui: h.Session.gameGuiService,
	}
	if err := tools.ListgamesTmpl.Execute(w, r, dat, user, false); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return
}

func (h *GameHandler) StartHandler(w http.ResponseWriter, r *http.Request) {
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	ctx := context.WithValue(context.Background(), gameKey, user)

	Gletter, gameID, err := h.Session.gameService.Start(ctx)
	if err != nil {
		fmt.Println("could not start game: ", err)
		tools.Error(w, http.StatusInternalServerError, h.Logger)
		return
	}
	dat := struct {
		Letter  string
		Game_id GameID
		UserGui *GameGuiService
	}{
		Letter:  Gletter,
		Game_id: gameID,
		UserGui: h.Session.gameGuiService,
	}
	if err := tools.GameTmpl.Execute(w, r, dat, user, false); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return
}

func (h *GameHandler) GuessHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n\n\n********************* In guess1 game *************************  \n\n\n\n")
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	game, err := h.gameFromForm(w, r)
	fmt.Printf("\n\n\n In guess2 game =  %v \n\n\n\n", *game)
	if err != nil {
		fmt.Println("could not save game: ", err)
		tools.Error(w, http.StatusInternalServerError, h.Logger)
		return
	}
	dat := struct {
		Letter  string
		Game_id GameID
		UserGui *GameGuiService
	}{
		Letter:  game.GuessedLetterS,
		Game_id: game.ID,
		UserGui: h.Session.gameGuiService,
	}

	if err := tools.GameTmpl.Execute(w, r, dat, user, false); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return
}

func (h *GameHandler) GuessLetterHandler(w http.ResponseWriter, r *http.Request) {
	// Decode request.
	var req request
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		tools.Error(w, tools.ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}
	fmt.Printf("\n\n\n********************* In guessLetterH game *************************  %vn\n\n\n", req)

	gm, err := h.Session.gameService.ServerGuessALetter(req.Letter, req.ID)
	if err != nil {
		fmt.Println("could not ServerGuessALetter game: ", err)
		tools.Error(w, http.StatusInternalServerError, h.Logger)
		return
	}
	if err := json.NewEncoder(w).Encode(&gm); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

func (h *GameHandler) SaveHandler(w http.ResponseWriter, r *http.Request) {
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	ctx := context.WithValue(context.Background(), gameKey, user)
	game, err := h.gameFromForm(w, r)
	if err != nil {
		fmt.Println("could not save game: ", err)
		tools.Error(w, http.StatusInternalServerError, h.Logger)
		return
	}
	game.JustStarted = "false"
	_, err = h.Session.gameService.AddGame(ctx, game)

	dat := struct {
		Game    *Game
		UserGui *GameGuiService
	}{
		Game:    game,
		UserGui: h.Session.gameGuiService,
	}

	if err := tools.GameTmpl.Execute(w, r, dat, user, false); err != nil {
		tools.Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}
	return
}

// deleteHandler deletes a given game.
func (h *GameHandler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	//check if user is authenticated before add game service
	user, err := h.Session.UserFromRequest(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	ids := chi.URLParam(r, "id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		fmt.Println("bad game id: ", err)
		tools.Error(w, http.StatusInternalServerError, h.Logger)
		return
	}
	ctx := context.WithValue(context.Background(), gameKey, user)
	err = h.Session.gameService.DeleteGame(ctx, GameID(id))
	if err != nil {
		fmt.Println("could not delete game: ", err)
		tools.Error(w, http.StatusInternalServerError, h.Logger)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return
}
func (h *GameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n\n\n ************************ In game handler for game ***********************  \n\n\n\n")
	h.mux.ServeHTTP(w, r)
}

// handleGetGame handles requests to create a new game.

// gameFromRequest retrieves a game from the database given a game ID in the
// URL's path.
func (h *GameHandler) gameFromRequest(r *http.Request) (*Game, error) {
	ids := chi.URLParam(r, "id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		fmt.Println("bad game id: ", err)
		return nil, nil
	}
	game, err := h.Session.gameService.GetGame(GameID(id))
	if err != nil {
		fmt.Println("could not save game: ", err)
		return nil, nil
	}
	return game, nil
}

func (h *GameHandler) gameFromForm(w http.ResponseWriter, r *http.Request) (*Game, error) {
	ids, _ := strconv.Atoi(r.FormValue("id"))
	wgs := []byte(r.FormValue("wronguesses"))
	game := &Game{
		PlayerID:       r.FormValue("playerid"),
		GuessedLetterS: r.FormValue("guessedletter"),
		WordSoFarS:     r.FormValue("wordsofar"),
		WrongGuesses:   wgs,
		Invalid:        r.FormValue("invalid"),
		ID:             GameID(ids),
		JustStarted:    r.FormValue("juststarted"),
		GameOver:       r.FormValue("gameover"),
	}

	fmt.Printf("\n\n\n In FromForm game =  %v \n\n\n\n", *game)
	return game, nil
}
