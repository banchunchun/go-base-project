package util

import (
	"bufio"
	"com.banxiaoxiao.server/logger"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func ExecCommand(timeout time.Duration, name string, arg ...string) ([]string, error) {
	var r []string

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			r = append(r, scanner.Text())
		}
	}()
	err = cmd.Run()
	stdout.Close()
	return r, err
}

func ExecCmd(timeout time.Duration, name string, arg ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, arg...)
	contents, err := cmd.Output()
	if err != nil {
		return "", err
	}
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Command timed out")
		return "", ctx.Err()
	}
	return string(contents), nil
}

func ListFSCmd(fs string) (map[string]string, error) {
	cmd := "mount -l | awk '$(NF-1)~/^(" + fs + ")/{print $1,$3}'"
	r, err := ExecCommand(10*time.Second, "bash", "-c", cmd)
	m := make(map[string]string)

	if err != nil {
		return m, err
	}

	for _, str := range r {
		d := strings.Split(str, " ")
		m[d[1]] = d[0]
	}
	return m, nil
}

func MountCmd(c string) error {
	cmd := "mount -t " + c
	r, err := ExecCommand(10*time.Second, "bash", "-c", cmd)
	logger.Log().Debug(r)
	return err
}

func RunBash(c string) (string, error) {
	return ExecCmd(10*time.Second, "bash", "-c", c)
}

func ExecBash(c string) error {
	r, err := ExecCommand(10*time.Second, "bash", "-c", c)
	logger.Log().Debug(r)
	return err
}

func UnmountCmd(c string) error {
	cmd := "umount -l " + c
	r, err := ExecCommand(10*time.Second, "bash", "-c", cmd)
	logger.Log().Debug(r)
	return err
}
