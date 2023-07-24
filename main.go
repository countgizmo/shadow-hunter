package main

import (
	"ziggytwister.com/shadow-hunter/transmitter"
	"ziggytwister.com/shadow-hunter/ui"
)

func main() {
	edn := transmitter.GetAppDB("localhost", "5555")
	ui.Render(edn)
}
