package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-getter"
)

var url = "https://api.github.com/repos/mikimou/mc-launcher/releases/latest"

func main() {
	time.Sleep(3 * time.Second)
	update, latest, current := checkUpdate()
	fmt.Println("Current:", current+"\nLatest:  "+latest, update)
	if update == true {
		fmt.Println("Updating..")
		installUpdate(latest)
	}
	//go exec.Command("./launcher.exe").Start()
	time.Sleep(3 * time.Second)
	os.Exit(1)
}

func checkUpdate() (bool, string, string) {
	var version string
	ver, err := os.ReadFile("launcher.version")
	if ver != nil {
		version = string(ver)
	} else {
		version = "error"
		fmt.Println("Error loading current version!")
	}
	if err != nil {
		version = "error"
		fmt.Println("Error loading current version!")
	}

	var Client = &http.Client{Timeout: 5 * time.Second}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("No internet connection!")
		}
	}()
	resp, err := Client.Get(url)
	if err != nil {
		fmt.Println("Server not responding")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	var jsonMap map[string]string
	json.Unmarshal([]byte(body), &jsonMap)

	if jsonMap["name"] != version {
		return true, jsonMap["name"], version
	}
	if jsonMap["name"] == "" {
		return false, "", version
	}
	return false, "", version
}

func installUpdate(ver string) {
	client := getter.Client{DisableSymlinks: true}
	client.Dst = "."
	client.Dir = true
	client.Src = "https://github.com/mikimou/mc-launcher/releases/download/" + ver + "/update.zip"
	fmt.Println(client.Src)
	err := client.Get()
	if err != nil {
		fmt.Println("Error")
	}
	if err == nil {
		fmt.Println("Finished")
	}
}
