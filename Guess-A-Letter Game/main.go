package main

import (
	"Game/apiapp"
	"Game/apigame"
	"Game/apiuser"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	addr := os.Getenv("PORT")
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) //relays incoming signals to sigs
	//[Start User Session]
	//Implemented User in apiuser and currently the only option
	var dbUser = make(apiuser.DBType)
	var dbGame = make(apigame.DBType)

	gameGui := apigame.NewGameGuiService(addr)

	us := apiuser.NewSession(dbUser)
	gs := apigame.NewSession(dbGame, us, gameGui)

	gh := apigame.NewGameHandler(gs)
	uh := apiuser.NewUserHandler(us)
	_ = uh.Session.Userservicesess.AddUser(&apiuser.User{Username: "bobsar0", Password: "Atib0b00", Level: "Admin"})
	//initialize the handler
	h := apiapp.NewHandler(uh, gh)
	//open apiuser server
	server := apiapp.NewServer(addr, h)
	if err := server.Open(done, sigs); err != nil {
		log.Fatal(err)
	}
	//    fresh      gukjj\
	fmt.Println("Listening on: ", server.Port())
	<-done
	fmt.Println("exiting")
}
