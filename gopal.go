/*
 * Gopal - Easy to use Golang Exec Library
 * Created by Noval Agung Prayogo <caknopal@gmail.com>
 * http://novalagung.com/
 */

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

const rawSimple = `package main
import (
	__LIBS__
)
func doStuff() interface{} {
	__CMD__
}
func main() {
	__MAIN__
}`

const rawAdvance = `package main
import (
	__LIBS__
)

__CMD__`

const (
	defaultDeleteTemporaryFile = true
	defaultPerm                = 0755
	rawConstLibs               = `__LIBS__`
	rawConstCmd                = `__CMD__`
	RrawConstMain              = `__MAIN__`
)

type Gopal struct {
	WillDeleteTemporaryFile bool
	TemporaryFolderName     string

	temporaryFolderPath string
	libs                []string
}

func New() *Gopal {
	g := new(Gopal)
	g.WillDeleteTemporaryFile = defaultDeleteTemporaryFile
	g.prepareDefaultTemporaryFolderName()
	g.prepareDefaultLibs()

	return g
}

func (g *Gopal) prepareDefaultTemporaryFolderName() {
	if runtime.GOOS == "windows" {
		g.TemporaryFolderName = "temp"
	}

	g.TemporaryFolderName = ".temp"
}

func (g *Gopal) prepareDefaultLibs() {
	g.libs = []string{"fmt"}
}

func (g *Gopal) prepareTemporaryFolderPath() {
	basePath, _ := os.Getwd()
	g.temporaryFolderPath = filepath.Join(basePath, g.TemporaryFolderName)
}

func (g *Gopal) prepareTemporaryFile() (string, *os.File, error) {
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

func (g *Gopal) DeleteTemporaryPathIfAllowed() {
	if !g.WillDeleteTemporaryFile {
		return
	}

	g.DeleteTemporaryPath()
}

func (g *Gopal) DeleteTemporaryPath() {
	os.RemoveAll(g.temporaryFolderPath)
}

func (g *Gopal) runCommand(fileLocation string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("go", "run", fileLocation)
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

func (g *Gopal) renderLibs(cmdString string) string {
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

func (g *Gopal) AddLibs(libs ...string) *Gopal {
	g.libs = append(g.libs, libs...)
	return g
}

func (g *Gopal) ExecuteSimple(cmdString string) (string, error) {
	defer g.DeleteTemporaryPathIfAllowed()
	cmdString = strings.TrimSpace(cmdString)

	if strings.HasPrefix(cmdString, "package") {
		return "", errors.New("Use `ExecuteRaw()` to exec complete golang file")
	}

	if strings.HasPrefix(cmdString, "func") {
		return "", errors.New("Use `Execute()` to exec code which contains `main()` func")
	}

	fileLocation, file, err := g.prepareTemporaryFile()
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		return "", err
	}

	hasPrint := (func() bool {
		for _, each := range []string{"Println", "Printf", "Print"} {
			if strings.Contains(cmdString, fmt.Sprintf("fmt.%s", each)) {
				return true
			}
		}

		return false
	}())

	callDoStuffString := `doStuff()`

	if !hasPrint {
		callDoStuffString = fmt.Sprintf(`fmt.Print(%s)`, callDoStuffString)

		if !strings.Contains(cmdString, "return") {
			cmdPart := strings.Split(cmdString, "\n")
			cmdPart[len(cmdPart)-1] = fmt.Sprintf("return %s", cmdString)
			cmdString = strings.Join(cmdPart, "\n")
		}
	} else {
		cmdString = fmt.Sprintf("%s\nreturn true", cmdString)
	}

	cmdString = strings.Replace(rawSimple, rawConstCmd, cmdString, -1)
	cmdString = strings.Replace(cmdString, RrawConstMain, callDoStuffString, -1)
	cmdString = g.renderLibs(cmdString)

	if _, err := file.WriteString(cmdString); err != nil {
		return "", err
	}

	return g.runCommand(fileLocation)
}

func (g *Gopal) Execute(cmdString string) (string, error) {
	defer g.DeleteTemporaryPathIfAllowed()
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

	cmdString = strings.Replace(rawAdvance, rawConstCmd, cmdString, -1)
	cmdString = g.renderLibs(cmdString)

	if _, err := file.WriteString(cmdString); err != nil {
		return "", err
	}

	return g.runCommand(fileLocation)
}

func (g *Gopal) ExecuteRaw(cmdString string) (string, error) {
	defer g.DeleteTemporaryPathIfAllowed()
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
