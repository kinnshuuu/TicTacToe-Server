package main

import (
	controller "TicTacToe-GolangServer/controller"
)

func main() {
	err := controller.Handler()
	if err != nil {
		panic("Can't start server")
	}
}
