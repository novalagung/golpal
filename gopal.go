package gopal

import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "runtime"
import "bytes"
import "time"
import "errors"
import "strings"

var TemporaryFolderName = (func(os string) string {
	if os == "windows" {
		return "temp"
	}

	return ".temp"
}(runtime.GOOS))

var DeleteTemporaryFile = true

var rawSimple = `package main
import "fmt"
func doStuff() interface{} {
	__CMD__
}
func main() {
	__MAIN__
}`

var basePath = (func(dir string, err error) string {
	return dir
}(os.Getwd()))

var temporaryFolderPath = filepath.Join(basePath, TemporaryFolderName)

const DEFAULT_PERM = 0755

func prepareTemporaryFile() (string, *os.File, error) {
	filename := fmt.Sprintf("temp-%d.go", time.Now().UnixNano())
	fileLocation := filepath.Join(temporaryFolderPath, filename)

	folder, err := os.Open(temporaryFolderPath)
	if folder != nil {
		defer folder.Close()
	}
	if os.IsNotExist(err) {
		if err := os.Mkdir(temporaryFolderPath, DEFAULT_PERM); err != nil {
			return fileLocation, nil, err
		}
	}

	file, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY, DEFAULT_PERM)
	return fileLocation, file, err
}

func runCommand(fileLocation string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("go", "run", fileLocation)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errString := err.Error()
		if stderr.Len() > 0 {
			errString = fmt.Sprintf("%s%s", fmt.Sprintln(errString), stderr.String())
		}

		return "", errors.New(strings.TrimSpace(errString))
	}

	return strings.TrimSpace(stdout.String()), nil
}

func deleteTemporaryPath() {
	if !DeleteTemporaryFile {
		return
	}

	os.RemoveAll(temporaryFolderPath)
}

func ExecuteSimple(cmdString string) (string, error) {
	defer deleteTemporaryPath()
	cmdString = strings.TrimSpace(cmdString)

	output := ""
	fileLocation, file, err := prepareTemporaryFile()
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return output, err
	}

	hasPrint := (func() bool {
		for _, each := range []string{"Println", "Printf", "Print"} {
			if strings.Contains(cmdString, fmt.Sprintf("fmt.%s", each)) {
				return true
			}
		}

		return false
	}())

	callDoStuff := `doStuff()`

	if !hasPrint {
		callDoStuff = fmt.Sprintf(`fmt.Print(%s)`, callDoStuff)

		if !strings.HasPrefix(cmdString, "return") {
			cmdString = fmt.Sprintf("return %s", cmdString)
		}
	} else {
		cmdString = fmt.Sprintf("%sreturn true", fmt.Sprintln(cmdString))
	}

	cmdString = strings.Replace(rawSimple, `__CMD__`, cmdString, -1)
	cmdString = strings.Replace(cmdString, `__MAIN__`, callDoStuff, -1)

	if _, err := file.WriteString(cmdString); err != nil {
		return "", err
	}

	return runCommand(fileLocation)
}
