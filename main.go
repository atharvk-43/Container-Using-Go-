package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

var pidFile = "container.pid"
var logFile = "container.log"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: run | ps | logs | stop | rm")
		return
	}

	switch os.Args[1] {

	case "run":
		run()

	case "child":
		child()

	case "ps":
		ps()

	case "logs":
		logs()

	case "stop":
		stop()

	case "rm":
		remove()

	default:
		fmt.Println("Unknown command")
	}
}

func run() {
	fmt.Printf("Running %v\n", os.Args[2:])

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	file, _ := os.Create(logFile)

	cmd.Stdout = file
	cmd.Stderr = file
	cmd.Stdin = os.Stdin

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	must(cmd.Start())

	pid := cmd.Process.Pid
	os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)

	fmt.Println("Container started with PID:", pid)

	cmd.Wait()
}

func child() {
	fmt.Printf("Running child %v\n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	must(cmd.Run())
}

func ps() {
	data, err := os.ReadFile(pidFile)

	if err != nil {
		fmt.Println("No running container")
		return
	}

	fmt.Println("Container PID:", strings.TrimSpace(string(data)))
}

func logs() {
	data, err := os.ReadFile(logFile)

	if err != nil {
		fmt.Println("No logs found")
		return
	}

	fmt.Println(string(data))
}

func stop() {
	data, err := os.ReadFile(pidFile)

	if err != nil {
		fmt.Println("No running container")
		return
	}

	pid, _ := strconv.Atoi(strings.TrimSpace(string(data)))

	process, err := os.FindProcess(pid)

	if err != nil {
		fmt.Println("Process not found")
		return
	}

	process.Kill()

	fmt.Println("Container stopped")
}

func remove() {
	os.Remove(pidFile)
	os.Remove(logFile)

	fmt.Println("Container metadata removed")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}