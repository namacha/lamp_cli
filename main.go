package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/morikuni/aec"
	"github.com/urfave/cli/v2"
)

type Device struct {
	Id          string
	Token       string
	NickName    string
	Last_Update float64
	Recipients  []string
	Alive       bool
	Testing     bool
}

const API = "https://lamp.deta.dev"

var TOKEN string = os.Getenv("LAMP_TOKEN")

func getDevices() []Device {

	req, _ := http.NewRequest("GET", API+"/devices", nil)
	req.Header.Set("Authorization", "Bearer "+TOKEN)

	client := new(http.Client)
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var devices []Device
	err = json.NewDecoder(resp.Body).Decode(&devices)
	if err != nil {
		log.Fatal(err)
	}

	return devices
}

func showDevices(devices []Device) {
	current := time.Now().Unix()

	for _, device := range devices {
		lastSeen := float64(current) - device.Last_Update
		fmt.Print(aec.Bold)

		fmt.Print("┌")
		fmt.Print(device.NickName)

		if device.Alive {
			fmt.Print(aec.GreenF.Apply(" ✔"))
		} else {
			fmt.Print(aec.RedF.Apply(" ✗"))
		}
		fmt.Println(aec.Bold)

		fmt.Printf("└Last seen: %.1fs ago\n", lastSeen)
	}
	fmt.Printf(aec.Reset)
}

func main() {
	app := &cli.App{
		Name:  "lamp",
		Usage: "lamp - cli tool for 'lamp'",
		Action: func(c *cli.Context) error {
			showDevices(getDevices())
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
