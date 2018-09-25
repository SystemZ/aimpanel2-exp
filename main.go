package main

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"time"
)

func main() {
	cmd := exec.Command("bash", "fake-server.sh")

	stdout, _ := cmd.StdoutPipe()
	stdin, _ := cmd.StdinPipe()

	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			io.WriteString(stdin, "test\r\n")
			time.Sleep(4 * time.Second)
		}
	}()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		log.Println(scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		log.Println(err)
	}
}
