package gopal

import (
	"os"
	"testing"
)

func TestRmdir(t *testing.T) {
	if err := os.RemoveAll(TemporaryFolderName); err != nil {
		t.Error(err)
	}
}

func TestSimpleCommand(t *testing.T) {
	cmdString := `3 + 2`
	output, err := ExecuteSimple(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("output", "=>", output)
}

func TestSimpleCommandWithReturn(t *testing.T) {
	cmdString := `return 3 + 2`
	output, err := ExecuteSimple(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("output", "=>", output)
}

func TestSimpleCommandWithFmt(t *testing.T) {
	cmdString := `fmt.Println(3 + 2)`
	output, err := ExecuteSimple(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("output", "=>", output)
}

func TestSimpleCommandIfElse(t *testing.T) {
	cmdString := `temp := 3
if temp == 2 {
	fmt.Println("wrong")
} else {
	fmt.Println("right")
}`

	output, err := ExecuteSimple(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("output", "=>", output)
}
