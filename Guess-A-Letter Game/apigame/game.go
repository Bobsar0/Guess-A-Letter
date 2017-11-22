//Step 1

package apigame

//Game contains the words as fields

type GameID int

type Game struct {
	ID             GameID
	CorrectWord    string
	GuessedWord    string
	GuessedLetterS string
	WordSoFar      []byte
	WordSoFarS     string
	WrongGuesses   []byte
	Word           string
	Playername     string
	PlayerID       string
	Count          int
	WinOrLose      bool
	ImageURL       string
	JustStarted    string
	Invalid        string
	GameOver       string
}

//GameService specifies methods that the system will use to operate the game
// type GameService interface {
// 	Start() (*Game string)
// 	Guess(*Game) (Correct bool)
// 	Result() (winorloss bool)
// }
