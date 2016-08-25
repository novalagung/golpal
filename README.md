<img align="right" alt="Golpal logo" src="https://raw.githubusercontent.com/novalagung/golpal/master/logo.png">

# Golpal

Easy to use Golang Eval Library

## Introduction

Golpal is simple library to allow developer to do **eval** operation on golang source codes. Technically Golang doesn't provide API to do some eval, so we use temporary file to achieve that.

[![Build Status](https://travis-ci.org/novalagung/golpal.png?branch=master)](https://travis-ci.org/novalagung/golpal)

## Table of contents

 - <a href="#user-content-installation")>Installation</a>
 - <a href="#user-content-cli-installation")>CLI Installation</a>
 - <a href="#user-content-example")>Example</a>
 - <a href="#user-content-api-reference")>API Reference</a>
 - <a href="#user-content-contributing")>Contributing</a>
 - <a href="#user-content-license")>License</a>
 - <a href="#user-content-author")>Author</a>


<h2 id="installation">Installation</h2>

> Stable version: v1.0.0

Just go get the lib is enough

```bash
go get -u github.com/novalagung/golpal
```

Run test

```bash
cd $GOPATH/src/github.com/novalagung/golpal
go test *.go -v
```


<h2 id="cli-installation">CLI Installation</h2>

By using cli, golpal will run faster. Follow these instructions to install golpal CLI.

#### Linux

```bash
sh install_cli.sh
golpal -content "3 + 4"
```

#### Windows

```bash
install_cli.bat
golpal.exe -content "3 + 4"
```


<h2 id="example">Example</h2>

#### Simple Example

```go
package main

import "github.com/novalagung/golpal"
import "fmt"

func main() {
	cmdString := `3 + 2`
	output, err := golpal.New().ExecuteSimple(cmdString)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("result", "=>", output)
}
```

> For one line statement using `ExecuteSimple()`, `return` keyword is optional

#### Another Example

```go
cmdString := `
	number := 3
	if number == 2 {
		return "wrong"
	} else {
		return "right"
	}
`

output, err := golpal.New().ExecuteSimple(cmdString)
if err != nil {
	fmt.Println(err)
}
fmt.Println("result", "=>", output)
```

> For multiline statement using `ExecuteSimple()`, `return` must be defined

#### Example which use `strings` and `runtime`

```go
cmdString := `
	osName := runtime.GOOS
	arr := []string{"my", "operation system", "is", osName}
	return strings.Join(arr, ", ")
`

output, err := golpal.New().AddLibs("strings", "runtime").ExecuteSimple(cmdString)
if err != nil {
	fmt.Println(err)
}
fmt.Println("result", "=>", output)
```

#### Example Not Simple (again, NOT SIMPLE, by using `Execute()` func)

```go
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
	}
`

output, err := golpal.New().Execute(cmdString)
if err != nil {
	fmt.Println(err)
}
fmt.Println("result", "=>", output)
```

#### Example Execute Raw (use `ExecuteRaw()` func)

```go
cmdString := `
	package main

	import "fmt"

	func main() {
		fmt.Print("hello")
	}
`

output, err := golpal.New().ExecuteRaw(cmdString)
if err != nil {
	fmt.Println(err)
}
fmt.Println("result", "=>", output)
```

#### CLI Example

```go
func prepareCmd(cmdString string, args ...string) (*exec.Cmd, *bytes.Buffer, *bytes.Buffer) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	if runtime.GOOS == "windows" {
		cmdString = fmt.Sprintf(".exe", cmdString)
	}

	cmd := exec.Command(cmdString, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	return cmd, &stdout, &stderr
}

func main() {
	cmd, stdout, stderr := prepareCmd("golpal", "-content", "3 + 1")
	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
		return
	}

	if stderr.Len() > 0 {
		fmt.Println(stderr.String())
		return
	}

	fmt.Println("result", "=>", strings.TrimSpace(stdout.String()))
}
```

#### More Example

For more examples please take a look at the [`golpal_test.go`](https://github.com/novalagung/golpal/blob/master/golpal_test.go) and  [`cli_test.go`](https://github.com/novalagung/golpal/blob/master/cli_test.go).


<h2 id="api-reference">API Reference</h2>

#### Func of `golpal`

There are only one func available, `golpal.New()` which return object that type is `*golpal.Golpal`

| Func        | Usage          |
| ----------- | :------------- |
| golpal.New() | instantiate new `*golpal.Golpal` object |

#### Properties of `*golpal.Golpal`

| Property    | Type | Usage          |
| ----------- | ---- | :------------- |
| .WillDeleteTemporaryFile | `bool` | Determine if temporary path will be deleted or not after executing the source codes (default is `true`) |
| .TemporaryPath | `string` | Path of temporary folder used to store all `*.go` temp files (default is `.temp` for *\*nix* / \*d*rwin*, and `temp` for *w\*ndows*) | Run golang source codes. The code will be placed inside virtual `main()` func |

#### Methods of `*golpal.Golpal`

| Method      | Usage          |
| ----------- | :------------- |
| .AddLibs(libs ...string) | Add other libraries, by default only `fmt` is included |
| .ExecuteSimple(cmdString&nbsp;string) | Run golang source codes. The code will be placed inside virtual `main()` func. This function doesn't allow `fmt.Print*()`. Also for multiline statement, `return` must be defined |
| .Execute(cmdString string) | Run golang source codes which contains `main()` func |
| .ExecuteRaw(cmdString string) | Run complete golang source code |
| .DeleteTemporaryPath() | Force delete temporary path which used to do the exec process |


<h2 id="cli-commands">CLI Commands</h2>

Golpal CLI (`golpal` / `golpal.exe`) has several arguments

| Args      | Usage          | Example |
| --------- | :------------- | :------ |
| -content | Eval from string | `-content="3 + 4"`
| -file | Eval from file | `-file="/path/to/file"`
| -mode | Golpal mode (`simple`, `normal`, `raw`). Default is `simple` | `-mode="raw"`
| -libs | Include libs | `-libs="io/ioutil, bytes, strings"`



<h2 id="contributing">Contributing</h2>

Feel free to contribute

`fork` -> `commit` -> `push` -> `pull request`


<h2 id="license">License</h2>

MIT License


<h2 id="author">Author</h2>

Noval Agung Prayogo - [http://novalagung.com/](http://novalagung.com)
