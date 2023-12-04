package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("bac command")
	}
}

func run() {
	fmt.Printf("running %v as %d\n", os.Args[2:], os.Getpid())
	// 調用自身命令列
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		// 希望container具有new ns(namespace) 但不與主機共享
		Unshareflags: syscall.CLONE_NEWNS,
	}
	must(cmd.Run())
}

func child() {
	fmt.Printf("running %v as %d\n", os.Args[2:], os.Getpid())
	cg()
	syscall.Sethostname([]byte("container"))
	// chroot must need /bin and other
	must(syscall.Chroot("/usr/"))
	syscall.Chdir("/")
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	must(cmd.Run())

	syscall.Unmount("/proc", 0)
}

func cg() {
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups, "pids")
	err := os.Mkdir(filepath.Join(pids, "liz"), 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	//限制容器可用最大量
	must(os.WriteFile(filepath.Join(pids, "liz/pids.max"), []byte("20"), 0700))
	//remove the new cgroup in place after the container exist
	must(os.WriteFile(filepath.Join(pids, "liz/notify_on_release"), []byte("1"), 0700))
	must(os.WriteFile(filepath.Join(pids, "liz/cgroup.procs"), []byte(strconv.Itoa(os.Getegid())), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
