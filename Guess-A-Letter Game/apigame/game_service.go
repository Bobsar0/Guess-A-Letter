//STEP 2
package apigame

import (
	"Game/apiuser"
	"Game/tools"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

var nextID GameID

func init() {
	nextID = 1
}

type key int

var gameKey key = 1

type GameService struct {
	session *Session
}
type DBType map[GameID]*Game

func (g *GameService) ServerGuessALetter(letter string, gameID GameID) (*Game, error) {

	gm, err := g.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	gm.GuessedLetterS = strings.ToUpper(letter)
	fmt.Println("Letter is: ", gm.GuessedLetterS)
	fmt.Printf("\n\n\n********************* In ServerGuessALetter0 *************************  %vn\n\n\n", gm)

	if gm.GuessedLetterS == "" {
		gm.JustStarted = "false"
		fmt.Printf("\n\n\n********************* In ServerGuessALetter1 (EMPTY GUESS) *************************  %vn\n\n\n", gm)
		if err := g.UpdateGame(gm); err != nil {
			return nil, err
		}
		return gm, nil
	}
	if !strings.Contains("ABCDEFGHIJKLMNOPQRSTUVWXYZ", gm.GuessedLetterS) {
		gm.Invalid = "true"
		gm.JustStarted = "true"
		if err := g.UpdateGame(gm); err != nil {
			return gm, err
		}
		return gm, nil
	}

	fmt.Printf("\n\n\n********************* In ServerGuessALetter2 *************************  %vn\n\n\n", gm)

	gm.Invalid = "false"
	gm.JustStarted = "true"

	if !strings.Contains(gm.CorrectWord, gm.GuessedLetterS) {
		fmt.Printf("\n\n\n********************* In ServerGuessALetter3 (wrong guess)*************************  %vn\n\n\n", gm)
		//gm.WrongGuesses = append(gm.WrongGuesses, gm.GuessedLetterS[0])
		gm.WrongGuesses[gm.Count] = gm.GuessedLetterS[0]
		gm.Count++
	} else {
		fmt.Printf("\n\n\n********************* In ServerGuessALetter4 (if right guess*************************  %vn\n\n\n", gm)
		for i := 0; i < len(gm.CorrectWord); i++ {
			if gm.CorrectWord[i] == gm.GuessedLetterS[0] {
				gm.WordSoFar[i] = gm.GuessedLetterS[0]
			}
		}
		gm.WordSoFarS = string(gm.WordSoFar)
		println(gm.WordSoFarS)
		fmt.Printf("\n\n\n********************* In ServerGuessALetter6 (After Updating WordSoFarS) *************************  %vn\n\n\n", gm)
	}

	fmt.Printf("\n\n\n********************* In ServerGuessALetter7 *************************  %vn\n\n\n", gm)

	gm.WinOrLose = g.Win(gm)

	if gm.Count == 7 || gm.WinOrLose == true {
		gm.GameOver = "true"
	}
	if err := g.UpdateGame(gm); err != nil {
		return gm, err
	}
	return gm, nil
}

func (g *GameService) Win(gm *Game) bool {
	gm.Word = gm.CorrectWord
	a := strings.ToLower(gm.WordSoFarS)
	b := strings.ToLower(gm.CorrectWord)
	return reflect.DeepEqual(a, b)
}

//mthods
func (g *GameService) ServeRender(letter string, gameID GameID) (string, error) {
	fmt.Println("ServeRender **********", letter, gameID)

	//return tools.Response{Hash: "#"}, nil
	return "#", nil
}
func (g *GameService) Start(ctx context.Context) (string, GameID, error) {
	//Retrieve current session user
	u, ok := ctx.Value(gameKey).(*apiuser.User)
	if !ok {
		return "", 0, tools.ErrUnauthorized
	}
	//retrieve word
	var bB []byte
	resp, _ := http.PostForm("http://watchout4snakes.com/wo4snakes/Random/RandomWord", url.Values{})
	if resp != nil {
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			defer resp.Body.Close()
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return "", 0, err
			}
			bB = bodyBytes
		}
	} else {
		bB = []byte("catarrh")
	}
	gm := &Game{
		CorrectWord: string(bB),
		PlayerID:    u.Username,
		JustStarted: "false",
		Count:       0,
	}
	length := len(gm.CorrectWord)
	gm.WordSoFar = make([]byte, length)
	gm.CorrectWord = strings.ToUpper(gm.CorrectWord)
	for i := 0; i < length; i++ {
		gm.WordSoFar[i] = '-'
	} //replace each letter with a hyphen
	gm.WrongGuesses = make([]byte, 7)
	for i := 0; i < 7; i++ {
		gm.WrongGuesses[i] = '-'
	} //replace each letter with a hyphen
	fmt.Println(gm)
	id, err := g.AddGame(ctx, gm)
	if err != nil {
		log.Fatal(err)
	}
	return gm.GuessedLetterS, id, nil

} //end start

//******************************************************************************/

func (g *GameService) AddGame(ctx context.Context, game *Game) (GameID, error) {
	// Retrieve current session user.
	u, ok := ctx.Value(gameKey).(*apiuser.User)
	if !ok {
		return 0, tools.ErrUnauthorized
	}
	db, ok := g.session.db.(DBType)
	if !ok {
		return 0, tools.ErrGamDbUnreachable
	}
	if game.PlayerID != u.Username {
		return game.ID, tools.ErrUnauthorized
	}
	db[nextID] = game
	game.ID = nextID
	nextID++
	return game.ID, nil
}

//*************************************************************************//
func (g *GameService) GetGame(gameid GameID) (*Game, error) {
	db, ok := g.session.db.(DBType)
	if !ok {
		return nil, tools.ErrGamDbUnreachable
	}
	if gameid != 0 {
		game := db[gameid]
		if game == nil {
			return nil, tools.ErrGameNotFound
		}
		return game, nil
	}
	return nil, tools.ErrGameIDRequired
}

//*****************************************************************************//
func (g *GameService) DeleteGame(ctx context.Context, gameid GameID) error {
	u, ok := ctx.Value(gameKey).(apiuser.User)
	if !ok {
		return tools.ErrUnauthorized
	}
	db, ok := g.session.db.(DBType)
	if !ok {
		return tools.ErrGamDbUnreachable
	}
	if gameid != 0 {
		game := db[gameid]
		if game == nil {
			return tools.ErrGameNotFound
		}
		if game.ID != gameid && game.PlayerID == u.Username {
			return tools.ErrUnauthorized
		}
		delete(db, gameid) //howto delete
	}
	return nil
}
func (g *GameService) ListGames() ([]*Game, error) {
	dbGames, ok := g.session.db.(DBType)
	if !ok {
		return nil, tools.ErrGamDbUnreachable
	}
	var db []*Game
	for _, games := range dbGames {
		db = append(db, games)
	}
	return db, nil
}

//***********************************************
func (m *GameService) UpdateGame(game *Game) error {

	db, ok := m.session.db.(DBType)
	if !ok {
		return tools.ErrGamDbUnreachable
	}
	// Only allow player to update Game.
	gameInDB, ok := db[game.ID]
	if !ok {
		return fmt.Errorf("memory db: game not found with ID %v", game.ID)
	} else if gameInDB.ID != game.ID { /* && gam.PlayerID == u.Username*/
		return tools.ErrUnauthorized
		//return fmt.Errorf("memory db: Non player not allowed to update Product %v", b.ID)
	}
	if game.ID == 0 {
		return errors.New("memory db: games with unassigned ID passed into updateProduct")
	}
	db[game.ID] = game
	return nil
}
