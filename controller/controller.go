package controller

import (
	"TicTacToe-Server/game"
	"log"
	"net/http"
)

func Handler() error {
	game := game.NewGame()
	go game.MatchOnlinePlayer()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		game.ServeClient(w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))

	return nil
}
