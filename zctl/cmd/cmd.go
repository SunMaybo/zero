//go:build !windows
// +build !windows

package cmd

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

func Execute(name, dir string, callbackResult func(lines string), args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	if outPipe, err := cmd.StdoutPipe(); err != nil {
		return err
	} else {
		if err := cmd.Start(); err != nil {
			log.Printf("Error starting command: %s......", err.Error())
			return err
		}
		cache := ""
		for {
			buf := make([]byte, 1024)
			num, err := outPipe.Read(buf)
			if err != nil && err != io.EOF {
				return err
			}
			if num > 0 {
				b := buf[:num]
				s := strings.Split(string(b), "\n")
				if len(s) <= 1 {
					cache += string(b)
				} else {
					line := strings.Join(s[:len(s)-1], "\n") //取出整行的日志
					if callbackResult != nil {
						callbackResult(fmt.Sprintf("%s%s", cache, line))
					}
					cache = s[len(s)-1]
				}
			} else {
				break
			}
		}

		if err := cmd.Wait(); err != nil {
			log.Printf("Error waiting for command execution: %s......", err.Error())
			return err
		}
	}
	return nil
}
