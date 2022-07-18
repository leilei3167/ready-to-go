package main

import (
	"fmt"
	"log"
	"os"
	"ready-to-go/Test/TDD/Practice/http-server/poker"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	store, close, err := poker.FileSysStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()
	game := poker.NewGame(poker.BlindAlerterFunc(poker.StdOutAlerter), store)
	cli := poker.NewCLI(os.Stdin, os.Stdout, game)
	cli.PlayPoker()
}
