package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	_, err := exec.LookPath("bower")
	if err != nil {
		fmt.Println("FATAL ERROR: could not find Bower\n" +
			"SEE: https://bower.io/ for installation instructions...")
		return
	}
	err = os.Chdir("./public")
	if err != nil {
		fmt.Printf("FATAL ERROR: could not descend into ./public folder: %v\n", err)
		return
	}
	fmt.Println("*** INSTALLING WEBCOMPONENTS ***")
	cmd := exec.Command("bower", "install", "--config.interactive=false", "-f")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("FATAL ERROR: could not fetch webcomponents: %v", err)
		return
	}
	fmt.Println("*** SUCCESS ***")
}
