package apigame

import (
	"Game/tools"
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
)

type GameGuiService struct {
	url  *url.URL
	Port string
}

func NewGameGuiService(Pt string) *GameGuiService {
	a := &GameGuiService{
		url:  &url.URL{},
		Port: Pt,
	}
	return a
}

type request struct {
	Letter string
	ID     GameID
}

func (g *GameGuiService) GuessALetter(letter string, gmID GameID) (*Game, error) {
	//Validate arguments.
	if gmID == 0 {
		return nil, tools.ErrGameRequired
	}
	req := request{
		Letter: letter,
		ID:     gmID,
	}

	b := new(bytes.Buffer)
	urlstr := fmt.Sprintf("%s", "http://localhost:"+g.Port+"/games/guessletter")

	json.NewEncoder(b).Encode(&req)
	resp, err := http.Post(urlstr, "application/json; charset=utf-8", b)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Decode response into JSON.
	var respBody Game
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, err
	}
	return &respBody, nil
}

func (g *GameGuiService) MustName(Name string) string {
	if Name == "" {
		return "Guess A Letter Game"
	}
	return Name
}

//for listgames.html
var gct int

func (g *GameGuiService) ListCount(m int) bool {
	fmt.Println("\ncout = ", gct)
	gct++
	if math.Mod(float64(m), 4) == 0 {
		return true
	}
	return false
}
