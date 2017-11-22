//STEP 3
//last
package apigame

import (
	"time"

	"Game/apiuser"
)

/*Session has the database handle the services can reference them. By making the
GameService a non-pointer field we reduce the allocations required when creating
a new session.*/
type Session struct {
	db             interface{}
	gameService    GameService
	gameGuiService *GameGuiService
	now            time.Time
	*apiuser.Session
}

func NewSession(uDB DBType, us *apiuser.Session, gamegui *GameGuiService) *Session {
	s := &Session{
		db:             uDB,
		Session:        us,
		gameGuiService: gamegui,
	}
	s.gameService.session = s
	return s
}
