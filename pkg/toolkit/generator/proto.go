package generator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func GenerateProtoFile(path string) {
	isFileExists := fileExists(path)
	if !isFileExists {
		log.Fatalf("The file %s does not exist.", path)
		return
	}

	isModel := strings.Contains(path, "proto/models/")
	cmdTemplate := `protoc \
  --proto_path=proto/models \
  --go_out=paths=source_relative:proto/exmsg/models \
  --go-grpc_out=paths=source_relative:proto/exmsg/models \
  %v`

	if !isModel {
		cmdTemplate = `protoc \
  --plugin=protoc-gen-micro=$GOPATH/bin/protoc-gen-micro \
  --proto_path=proto \
  --go_out=paths=source_relative:proto/exmsg \
  --go-grpc_out=paths=source_relative:proto/exmsg \
  --micro_out=paths=source_relative:proto/exmsg \
  services/%v`
	}

	paths := strings.Split(path, "/")
	fileName := paths[len(paths)-1]
	cmd := fmt.Sprintf(cmdTemplate, fileName)
	err := executeCommand(cmd)
	if err != nil {
		log.Fatalf("Error running command: %v", err)
	}
	fmt.Printf("Generated %s\n", fileName)
}

func executeCommand(cmd string) error {
	command := fmt.Sprintf("/bin/sh -c '%s'", cmd)
	fmt.Println("Executing command:", command)
	cmdObj := exec.Command("/bin/sh", "-c", command)
	output, err := cmdObj.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error running command: %v, output: %s", err, output)
	}
	return nil
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
