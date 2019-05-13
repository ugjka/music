package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	_, err := exec.LookPath("npm")
	if err != nil {
		fmt.Println("FATAL ERROR: could not find NPM\n" +
			"SEE: https://nodejs.org/ for installation instructions...")
		return
	}
	_, err = exec.LookPath("polymer")
	if err != nil {
		fmt.Println("FATAL ERROR: could not find Polymer CLI\n" +
			"SEE: https://polymer-library.polymer-project.org/3.0/docs/tools/polymer-cli for installation instructions...")
		return
	}
	err = os.Chdir("./public")
	if err != nil {
		fmt.Printf("FATAL ERROR: could not descend into ./public folder: %v\n", err)
		return
	}
	fmt.Println("*** INSTALLING WEBCOMPONENTS ***")
	cmd := exec.Command("npm", "install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("FATAL ERROR: could not fetch webcomponents: %v", err)
		return
	}
	fmt.Println("*** SUCCESS ***")
	err = os.Chdir("../")
	if err != nil {
		fmt.Printf("FATAL ERROR: could not descend into main folder: %v\n", err)
		return
	}
	fmt.Println("*** BUILDING POLYMER SOURCES ***")
	cmd = exec.Command("polymer", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("FATAL ERROR: could not build Polymer sources: %v", err)
		return
	}
	fmt.Println("*** SUCCESS ***")
}
