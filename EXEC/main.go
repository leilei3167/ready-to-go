package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("ls", "-lah")
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to call cmd.Run(): %v", err)
	}
}
