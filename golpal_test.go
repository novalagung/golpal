/*
 * Golpal - Easy to use Golang Eval Library
 * Created by Noval Agung Prayogo <caknopal@gmail.com>
 * http://novalagung.com/
 */

package golpal

import (
	"testing"
)

func TestExecuteSimple(t *testing.T) {
	cmdString := `3 + 2`
	output, err := New().ExecuteSimple(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "=>", output)
}

func TestExecuteSimpleWithReturn(t *testing.T) {
	cmdString := `return 3 + 2`
	output, err := New().ExecuteSimple(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "=>", output)
}

func TestExecuteSimpleIfElseWithReturn(t *testing.T) {
	cmdString := `
number := 3
if number == 2 {
	return "wrong"
} else {
	return "not always right"
}`

	output, err := New().ExecuteSimple(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "=>", output)
}

func TestExecuteSimpleWithOtherLibs(t *testing.T) {
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

func TestExecuteSimpleWithManualDeleteTemporaryPath(t *testing.T) {
	cmdsString := []string{"1 + 2", "2 + 3", "3 + 4", "4 + 5"}

	g := New()
	g.WillDeleteTemporaryFile = false
	g.TemporaryPath = "../testpath"

	for _, each := range cmdsString {
		output, err := g.ExecuteSimple(each)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("output of", each, "=>", output)
	}

	g.DeleteTemporaryPath()
}

func TestExecute(t *testing.T) {
	cmdString := `
func calculate(values ...int) int {
	total := 0
	for _, each := range values {
		total = total + each
	}
	return total
}

func main() {
	res := calculate(1, 2, 3, 4, 2, 3, 1)
	fmt.Printf("total : %d", res)
}`

	output, err := New().Execute(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "=>", output)
}

func TestExecuteWithOtherLibs(t *testing.T) {
	cmdString := `
func rand() string {
	timestamp := time.Now().UnixNano()
	timestampString := fmt.Sprintf("%d", timestamp)
	result := strings.Replace(timestampString, "0", "o", -1)
	return result
}

func main() {
	random := rand()
	fmt.Print(random)
}`

	output, err := New().AddLibs("time", "strings").Execute(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "=>", output)
}

func TestExecuteRaw(t *testing.T) {
	cmdString := `
package main

import "fmt"

func main() {
	fmt.Print("hello")
}`

	output, err := New().ExecuteRaw(cmdString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result", "=>", output)
}

func TestWrongMethod(t *testing.T) {
	cmdString := `
package main

import "fmt"

func main() {
	fmt.Print("hello")
}`

	if _, err := New().ExecuteSimple(cmdString); err != nil {
		t.Log(err.Error())
	}
}
