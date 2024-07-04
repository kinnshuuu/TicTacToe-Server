package main

import (
	controller "TicTacToe-Server/controller"
)

func main() {
	err := controller.Handler()
	if err != nil {
		panic("Can't start server")
	}
}
