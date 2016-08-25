/*
 * Golpal - Easy to use Golang Eval Library
 * Created by Noval Agung Prayogo <caknopal@gmail.com>
 * http://novalagung.com/
 */

package golpal

import "fmt"
import "os"
import "os/exec"
import "path/filepath"
import "runtime"
import "bytes"
import "time"
import "errors"
import "strings"

const rawExecuteSimple = `package main
import (
	__LIBS__
)
func doStuff() interface{} {
	__CMD__
}
func main() {
	fmt.Println(doStuff())
}`

const rawExecute = `package main
import (
	__LIBS__
)

__CMD__`

const (
	Version                               = "v1.0.0"
	defaultWillDeleteTemporaryPath        = true
	defaultPerm                           = 0755
	rawConstLibs                          = `__LIBS__`
	rawConstCmd                           = `__CMD__`
	temporaryFolderPathSubFolder   string = "golpal-files"
)

type Golpal struct {
	WillDeleteTemporaryPath bool
	TemporaryPath           string

	libs []string
}

func New() *Golpal {
	g := new(Golpal).init()

	return g
}

func (g *Golpal) init() *Golpal {
	g.WillDeleteTemporaryPath = defaultWillDeleteTemporaryPath

	// ===== prepare default temporary folder path
	basePath, _ := os.Getwd()
	if runtime.GOOS == "windows" {
		g.TemporaryPath = filepath.Join(basePath, "temp")
	}
	g.TemporaryPath = filepath.Join(basePath, ".temp")

	// ===== prepare default libs
	g.libs = []string{"fmt"}

	return g
}

func (g *Golpal) prepareTemporaryFile() (string, *os.File, error) {
	filename := fmt.Sprintf("temp-%d.go", time.Now().UnixNano())
	fileLocation := filepath.Join(g.getTemporaryFolderExactPath(), filename)

	folder, err := os.Open(g.getTemporaryFolderExactPath())
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(g.getTemporaryFolderExactPath(), defaultPerm); err != nil {
				return fileLocation, nil, err
			}
		}
	}
	defer folder.Close()

	file, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY, defaultPerm)
	return fileLocation, file, err
}

func (g *Golpal) getTemporaryFolderExactPath() string {
	return filepath.Join(g.TemporaryPath, temporaryFolderPathSubFolder)
}

func (g *Golpal) deleteTemporaryPathIfAllowed() {
	if !g.WillDeleteTemporaryPath {
		return
	}

	g.DeleteTemporaryPath()
}

func (g *Golpal) containsPrint(cmdString string) bool {
	for _, each := range []string{"Println", "Printf", "Print"} {
		if strings.Contains(cmdString, fmt.Sprintf("fmt.%s", each)) {
			return true
		}
	}

	return false
}

func (g *Golpal) renderLibs(cmdString string) string {
	quotedLibString := []string{}

	for _, each := range g.libs {
		if !strings.HasPrefix(`"`, each) {
			each = fmt.Sprintf(`"%s`, each)
		}
		if !strings.HasSuffix(`"`, each) {
			each = fmt.Sprintf(`%s"`, each)
		}

		quotedLibString = append(quotedLibString, each)
	}

	libsString := strings.Join(quotedLibString, "\n")
	res := strings.Replace(cmdString, rawConstLibs, libsString, -1)

	return res
}

func (g *Golpal) runCommand(fileLocation string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var cmd *exec.Cmd

	goRun := fmt.Sprintf("go run %s", fileLocation)

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", goRun)
	} else {
		cmd = exec.Command("bash", "-c", goRun)
	}

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errString := err.Error()
		if stderr.Len() > 0 {
			errString = fmt.Sprintf("%s\n%s", errString, stderr.String())
		}

		return "", errors.New(strings.TrimSpace(errString))
	}

	return strings.TrimSpace(stdout.String()), nil
}

// ===== Public Funcs

func (g *Golpal) DeleteTemporaryPath() *Golpal {
	os.RemoveAll(g.getTemporaryFolderExactPath())
	return g
}

func (g *Golpal) AddLibs(libs ...string) *Golpal {
	for _, lib := range libs {
		lib = strings.TrimSpace(lib)

		if lib == "fmt" || lib == "" {
			continue
		}

		g.libs = append(g.libs, lib)
	}

	return g
}

func (g *Golpal) ExecuteSimple(cmdString string) (string, error) {
	defer g.deleteTemporaryPathIfAllowed()
	cmdString = strings.TrimSpace(cmdString)

	if strings.HasPrefix(cmdString, "package") {
		return "", errors.New("Use `ExecuteRaw()` to exec complete golang file")
	}

	if strings.HasPrefix(cmdString, "func") {
		return "", errors.New("Use `Execute()` to exec code which contains `main()` func")
	}

	if g.containsPrint(cmdString) {
		return "", errors.New("The code cannot contain fmt.Println(), fmt.Printf(), or fmt.Print()")
	}

	fileLocation, file, err := g.prepareTemporaryFile()
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return "", err
	}

	if len(strings.Split(cmdString, "\n")) == 1 {
		if !strings.HasPrefix(cmdString, "return") {
			cmdString = fmt.Sprintf("return %s", cmdString)
		}
	}

	cmdString = strings.Replace(rawExecuteSimple, rawConstCmd, cmdString, -1)
	cmdString = g.renderLibs(cmdString)

	if _, err := file.WriteString(cmdString); err != nil {
		return "", err
	}

	return g.runCommand(fileLocation)
}

func (g *Golpal) Execute(cmdString string) (string, error) {
	defer g.deleteTemporaryPathIfAllowed()
	cmdString = strings.TrimSpace(cmdString)

	if strings.HasPrefix(cmdString, "package") {
		return "", errors.New("Use `ExecuteRaw()` to exec complete golang file")
	}

	fileLocation, file, err := g.prepareTemporaryFile()
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return "", err
	}

	cmdString = strings.Replace(rawExecute, rawConstCmd, cmdString, -1)
	cmdString = g.renderLibs(cmdString)

	if _, err := file.WriteString(cmdString); err != nil {
		return "", err
	}

	return g.runCommand(fileLocation)
}

func (g *Golpal) ExecuteRaw(cmdString string) (string, error) {
	defer g.deleteTemporaryPathIfAllowed()
	cmdString = strings.TrimSpace(cmdString)

	fileLocation, file, err := g.prepareTemporaryFile()
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return "", err
	}

	if _, err := file.WriteString(cmdString); err != nil {
		return "", err
	}

	return g.runCommand(fileLocation)
}
