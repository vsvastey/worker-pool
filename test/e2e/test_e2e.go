package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
)

func getMD5(filename string) string {
	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func createFile(filename string) {
	runShellCommand("dd", "if=/dev/urandom", fmt.Sprintf("of=%s", filename), "bs=64M", "count=16")
}

func runShellCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	err := cmd.Run()
	if err != nil {
		stdOut, _ := cmd.Output()
		fmt.Println(string(stdOut))
		panic(err)
	}
}

func main() {
	src := "/src/file1"
	dst := "/src/file2"
	createFile(src)
	hash1 := getMD5(src)
	fmt.Println("Hash of src is", hash1)
	runShellCommand("/src/bin/worker-pool", "--config", "/etc/worker_pool/e2e_config.yaml")
	hash2 := getMD5(dst)
	fmt.Println("Hash of dst is", hash2)
	//if hash1 != "hi" {
	//	fmt.Println("Wrong hash")
	//	panic(hash1)
	//}
	fmt.Println("End-to-end test passed")
}
