package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// docker           run image <CMD> <ARG>
// go run main.go   run       <CMD> <ARG>

// Step3: 새로운 UTS 설정 추가. Clone flag 추가 NEW UTS namespace. hostname 수동 변경 실습.
// 실습
// $ go run . run /bin/sh
// $ hostname box

func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		os.Exit(1)
	}
}

func run() {
	fmt.Printf("Running: %v\n", os.Args[2:])
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
	/*
		syscall.CLONE_NEWUTS는 리눅스의 네임스페이스 중 UTS(Unix Timesharing System) 네임스페이스를 새로 생성하겠다는 의미입니다.
		UTS 네임스페이스는 시스템의 **호스트 이름(hostname)**과 **도메인 이름(domain name)**을 분리된 상태로 유지할 수 있게 해줍니다.
		즉, 이 플래그를 설정하면 새로 생성된 프로세스는 독립된 UTS 네임스페이스에서 실행되며, 부모 프로세스와는 다른 호스트 이름이나 도메인 이름을 설정할 수 있습니다.
	*/

	must(cmd.Run())
}

// error 발생했을 때 에러 확인 후 프로그램 종료하도록 하는 함수
func must(err error) {
	if err != nil {
		panic(err)
	}
}
