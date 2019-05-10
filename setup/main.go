package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	_, err := exec.LookPath("bower")
	if err != nil {
		fmt.Println("FATAL ERROR: could not find Bower executable\n" +
			"SEE: https://bower.io/ for installation instructions...")
		return
	}
	err = os.Chdir("./public")
	if err != nil {
		fmt.Printf("FATAL ERROR: could not descend into ./public folder: %v\n", err)
		return
	}
	cmd := exec.Command("bower", "install", "--config.interactive=false", "-f")
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("FATAL ERROR: could not fetch the Bower dependencies: %v", err)
	}
}
