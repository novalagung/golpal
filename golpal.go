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
	Version                    = "v1.0.0"
	defaultDeleteTemporaryFile = true
	defaultPerm                = 0755
	rawConstLibs               = `__LIBS__`
	rawConstCmd                = `__CMD__`
)

type Golpal struct {
	WillDeleteTemporaryFile bool
	TemporaryFolderName     string

	temporaryFolderPath string
	libs                []string
}

func New() *Golpal {
	g := new(Golpal)
	g.WillDeleteTemporaryFile = defaultDeleteTemporaryFile
	g.prepareDefaultTemporaryFolderName()
	g.prepareDefaultLibs()

	return g
}

func (g *Golpal) prepareDefaultTemporaryFolderName() {
	if runtime.GOOS == "windows" {
		g.TemporaryFolderName = "temp"
	}

	g.TemporaryFolderName = ".temp"
}

func (g *Golpal) prepareDefaultLibs() {
	g.libs = []string{"fmt"}
}

func (g *Golpal) prepareTemporaryFolderPath() {
	basePath, _ := os.Getwd()
	g.temporaryFolderPath = filepath.Join(basePath, g.TemporaryFolderName)
}

func (g *Golpal) prepareTemporaryFile() (string, *os.File, error) {
	g.prepareTemporaryFolderPath()

	filename := fmt.Sprintf("temp-%d.go", time.Now().UnixNano())
	fileLocation := filepath.Join(g.temporaryFolderPath, filename)

	folder, err := os.Open(g.temporaryFolderPath)
	if folder != nil {
		defer folder.Close()
	}
	if os.IsNotExist(err) {
		if err := os.Mkdir(g.temporaryFolderPath, defaultPerm); err != nil {
			return fileLocation, nil, err
		}
	}

	file, err := os.OpenFile(fileLocation, os.O_CREATE|os.O_WRONLY, defaultPerm)
	return fileLocation, file, err
}

func (g *Golpal) deleteTemporaryPathIfAllowed() {
	if !g.WillDeleteTemporaryFile {
		return
	}

	g.DeleteTemporaryPath()
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

func (g *Golpal) containsPrint(cmdString string) bool {
	for _, each := range []string{"Println", "Printf", "Print"} {
		if strings.Contains(cmdString, fmt.Sprintf("fmt.%s", each)) {
			return true
		}
	}

	return false
}

func (g *Golpal) DeleteTemporaryPath() {
	os.RemoveAll(g.temporaryFolderPath)
}

func (g *Golpal) AddLibs(libs ...string) *Golpal {
	g.libs = append(g.libs, libs...)
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
