package main

import (
	"context"
	"fmt"
	"github.com/yahoo/vssh"
	"log"
	"time"
)

func main() {
	fmt.Println("starting test")

	vs := vssh.New().Start()
	config := vssh.GetConfigUserPass("pi", "")
	vs.AddClient("", config, vssh.SetMaxSessions(4))
	vs.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := "cat README.md"
	timeout, _ := time.ParseDuration("6s")
	respChan := vs.Sftp(ctx, "/Users/kenobrien/Development/KenObie/vssh/README.md", "/home/pi", timeout, 0)

	fmt.Println("test")
	resp := <-respChan
	if err := resp.Err(); err != nil {
		log.Fatal(err)
	}

	stream := resp.GetStream()
	defer stream.Close()

	for stream.ScanStdout() {
		txt := stream.TextStdout()
		fmt.Println(txt)
	}

	cmdRespChan := vs.Run(ctx, cmd, timeout)
	rp := <-cmdRespChan
	if err := resp.Err(); err != nil {
		log.Fatal(err)
	}

	cmdstream := rp.GetStream()
	defer cmdstream.Close()

	for cmdstream.ScanStdout() {
		txt := cmdstream.TextStdout()
		fmt.Println(txt)
	}

}
