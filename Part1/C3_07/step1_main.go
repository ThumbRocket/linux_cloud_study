package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

// docker           run image <CMD> <ARG>
// go run main.go   run       <CMD> <ARG>

// Step1: cgroup을 사용하여 컨테이너 내 process 갯수 제한.
// 실습
// $ go run . run /bin/sh
// $ :(){ :|:& };:

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		os.Exit(1)
	}
}

func run() {
	fmt.Printf("Running: %v as %d\n", os.Args[2:], os.Getpid())
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	cmd.Run()
}

func child() {
	fmt.Printf("Running child: %v as %d\n", os.Args[2:], os.Getpid())

	cg()

	syscall.Sethostname([]byte("container"))

	syscall.Chroot("/tmp/ubuntu")
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")
	defer syscall.Unmount("proc", 0)

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	must(cmd.Run())

}

func cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, "linux_campus"), 0755)
	must(ioutil.WriteFile(filepath.Join(pids, "linux_campus/pids.max"), []byte("20"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "linux_campus/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "linux_campus/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
