package internal

import (
	"context"
	"io"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/urfave/cli/v2"
	"golang.org/x/net/proxy"
)

func CreateRipCli() *cli.App {
	ripCli := cli.NewApp()

	setDocFields(ripCli)

	ripCli.Action = activateRefresh

	return ripCli
}

func setDocFields(cli *cli.App) {
	cli.Name = "rip"
	cli.Usage = "RIP is a CLI utility to refreshing your IP address using the Tor."
	cli.Version = "0.0.1"
}

func activateRefresh(c *cli.Context) error {
	log.Println("Starting RIP...")
	ctx, cancel := context.WithTimeout(c.Context, time.Minute*5)

	go func() {
		for {
			time.Sleep(30 * time.Second)
			log.Println("Restarting Tor...")
			cmd := exec.Command("sv", "restart", "tor")
			if err := cmd.Run(); err != nil {
				cancel()
			}
			verifyIP()
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("RIP finished")
		return nil
	}
}

func verifyIP() {
	dialer, err := proxy.SOCKS5("tcp", "127.0.0.1:9050", nil, proxy.Direct)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: dialer.Dial,
		},
	}

	res, err := client.Get("https://api.ipify.org")
	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Your new IP: %s\n", string(bodyBytes))
}
