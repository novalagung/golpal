/*
 * Golpal - Easy to use Golang Eval Library
 * Created by Noval Agung Prayogo <caknopal@gmail.com>
 * http://novalagung.com/
 */

package golpal

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func prepareCmd(cmdString string, args ...string) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(cmdString, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	return cmd, &stdout, &stderr
}

func TestPrepareCli(t *testing.T) {
	var err error

	if runtime.GOOS == "windows" {
		err = exec.Command("cmd", "/C", `install_cli.bat`).Run()
	} else {
		err = exec.Command("bash", "-c", `sh install_cli.sh`).Run()
	}

	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log("result", "=>", `golpal cli successfully installed!`)
}

func TestCliSimple(t *testing.T) {
	cmd, stdout, stderr := prepareCmd("golpal", "-content", "3 + 1")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
		return
	}

	if stderr.Len() > 0 {
		t.Fatal(stderr.String())
		return
	}

	t.Log("result", "=>", strings.TrimSpace(stdout.String()))
}

func TestCliSimpleMultiline(t *testing.T) {
	cmdString := `
number := 3
if number == 2 {
	return "wrong"
} else {
	return "not always right"
}`

	cmd, stdout, stderr := prepareCmd("golpal", "-content", cmdString)
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
		return
	}

	if stderr.Len() > 0 {
		t.Fatal(stderr.String())
		return
	}

	t.Log("result", "=>", strings.TrimSpace(stdout.String()))
}

func TestCliSimpleMultilineWithLibs(t *testing.T) {
	cmdString := `
osName := runtime.GOOS
arr := []string{"my", "operation system", "is", osName}
return strings.Join(arr, ", ")`

	cmd, stdout, stderr := prepareCmd("golpal", "-content", cmdString, "-libs", "strings, runtime")
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
		return
	}

	if stderr.Len() > 0 {
		t.Fatal(stderr.String())
		return
	}

	t.Log("result", "=>", strings.TrimSpace(stdout.String()))
}
