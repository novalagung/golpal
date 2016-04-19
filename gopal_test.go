package gopal

import (
	"testing"
)

func TestSimpleCommand(t *testing.T) {
	cmdString := `3 + 2`
	output, err := New().ExecuteSimple(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "=>", output)
}

func TestSimpleCommandWithReturn(t *testing.T) {
	cmdString := `return 3 + 2`
	output, err := New().ExecuteSimple(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "=>", output)
}

func TestSimpleCommandWithFmt(t *testing.T) {
	cmdString := `fmt.Println(3 + 2)`
	output, err := New().ExecuteSimple(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "=>", output)
}

func TestSimpleCommandIfElse(t *testing.T) {
	cmdString := `
temp := 3
if temp == 2 {
	fmt.Println("wrong")
} else {
	fmt.Println("right")
}`

	output, err := New().ExecuteSimple(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "=>", output)
}

func TestSimpleCommandWithOtherLib(t *testing.T) {
	cmdString := `
osName := runtime.GOOS
arr := []string{"my", "operation system", "is", osName}
return strings.Join(arr, ", ")`

	output, err := New().AddLibs("strings", "runtime").ExecuteSimple(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "=>", output)
}

func TestSimpleCommandWithManualDeleteTemporaryPath(t *testing.T) {
	cmdsString := []string{"1 + 2", "2 + 3", "3 + 4", "4 + 5"}

	g := New()
	g.WillDeleteTemporaryFile = false

	for _, each := range cmdsString {
		output, err := g.ExecuteSimple(each)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("output of", each, "=>", output)
	}

	g.DeleteTemporaryPath()
}
