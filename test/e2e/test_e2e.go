package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func getMD5(filename string) string {
	file, err := os.Open(filename)

	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		log.Panic(err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func createFile(filename string) {
	runShellCommand("dd", "if=/dev/urandom", fmt.Sprintf("of=%s", filename), "bs=64M", "count=4", "iflag=fullblock")
}

func runShellCommand(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	err := cmd.Run()
	if err != nil {
		stdOut, _ := cmd.Output()
		log.Println(string(stdOut))
		log.Panic(err)
	}
}

func main() {
	const (
		copySrc = "/src/file1"
		copyDst = "/src/file2"
		toS3    = "/src/to_s3"
		fromS3  = "/src/from_s3"
	)

	createFile(copySrc)
	createFile(toS3)

	hashCopySrc := getMD5(copySrc)
	hashS3Src := getMD5(toS3)

	// Do the jobs
	runShellCommand("/src/bin/worker-pool", "--config", "/etc/worker_pool/e2e_config.yaml")

	// check the copy file job
	hashCopyDst := getMD5(copyDst)
	if hashCopyDst != hashCopySrc {
		log.Fatalln("Copy file task failed")
	}

	// check the s3 upload job
	//aws --no-verify-ssl --endpoint-url http://localhost:9000 s3 cp s3://testbucket/just/uploaded.data downloaded.data
	// download from s3
	os.Setenv("AWS_ACCESS_KEY_ID", "minio")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "minio123")
	runShellCommand("aws", "--endpoint-url", "http://minio:9000",
		"s3", "cp", "s3://testbucket/e2e/test/uploaded.data", fromS3)
	hashS3Dst := getMD5(fromS3)
	if hashS3Dst != hashS3Src {
		log.Fatalln("Copy file task failed")
	}

	log.Print("End-to-end test passed")
}
