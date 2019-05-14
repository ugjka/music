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
	_, err = os.Stat("./node_modules/.bin/polymer")
	if err != nil {
		fmt.Println("ERROR: Could not find Polymer CLI")
		fmt.Println("*** INSTALLING POLYMER-CLI ***")
		cmd := exec.Command("npm", "install", "--prefix=./", "polymer-cli")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("FATAL ERROR: Could not install Polymer-CLI: %v\n", err)
			return
		}
		fmt.Println("*** SUCCESS ***")
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
		fmt.Printf("FATAL ERROR: could not return to the main folder: %v\n", err)
		return
	}
	fmt.Println("*** BUILDING POLYMER SOURCES ***")
	cmd = exec.Command("./node_modules/.bin/polymer", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("FATAL ERROR: could not build Polymer sources: %v", err)
		return
	}
	fmt.Println("*** SUCCESS ***")
}
