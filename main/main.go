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
	config := vssh.GetConfigUserPass("", "")
	vs.AddClient("", config, vssh.SetMaxSessions(4))

	vs.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//cmd := "cat README.md"
	timeout, _ := time.ParseDuration("6s")
	respChan := vs.Sftp(ctx, "/Users/kenobrien/Development/KenObie/vssh/README.md", "/home/pi", timeout, 0)
	respChan2 := vs.Sftp(ctx, "/Users/kenobrien/Development/KenObie/vssh/Contributing.md", "/home/pi", timeout, 0)

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

	resp2 := <-respChan2
	if err := resp2.Err(); err != nil {
		log.Fatal(err)
	}

	stream2 := resp.GetStream()
	defer stream2.Close()

	for stream2.ScanStdout() {
		txt := stream2.TextStdout()
		fmt.Println(txt)
	}

}
