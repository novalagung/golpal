# Golpal

<img align="right" width="400" alt="Golpal logo" src="https://cloud.githubusercontent.com/assets/982868/14644972/14fd4faa-067e-11e6-97db-d9b18cdf6c4b.png">

Easy to use Golang Exec Library

> Stable version: v1.0.0

## Introduction

Golpal is simple library to allow developer exec-ing golang source codes.

Technically Golang doesn't provide API to do some exec, so we use temporary file to achieve that.

[![Build Status](https://travis-ci.org/novalagung/golpal.png?branch=master)](https://travis-ci.org/novalagung/golpal)

## Installation

Just go get the lib is enough

```bash
go get -u github.com/novalagung/golpal
```

Run test

```bash
cd $GOPATH/src/github.com/novalagung/golpal
go test *.go -v
```

## Example

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

#### Another Example

```go
cmdString := `
	number := 3
	if number == 2 {
		fmt.Println("wrong")
	} else {
		fmt.Println("right")
	}
`

output, err := golpal.New().ExecuteSimple(cmdString)
if err != nil {
	fmt.Println(err)
}
fmt.Println("result", "=>", output)
```

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

#### More Example

For more examples please take a look at the [`golpal_test.go` file](https://github.com/novalagung/golpal/blob/master/golpal_test.go).

## API Reference

#### Func of `golpal`

There are only one func available, `golpal.New()` which return object that type is `*golpal.Golpal`

| Func        | Usage          |
| ----------- | :------------- |
| golpal.New() | instantiate new `*golpal.Golpal` object |

#### Property of `*golpal.Golpal`

| Property    | Type | Usage          |
| ----------- | ---- | :------------- |
| .WillDeleteTemporaryFile | `bool` | Determine if temporary path will be deleted or not after executing the source codes (default is `true`) |
| .TemporaryFolderName | `string` | Name of temporary folder used to store all `*.go` temp files (default is `.temp` for *\*nix* / \*d*rwin*, and `temp` for *w\*ndows*) | Run golang source codes. The code will be placed inside virtual `main()` func |

#### Methods of `*golpal.Golpal`

| Method      | Usage          |
| ----------- | :------------- |
| .AddLibs(libs ...string) | Add other libraries, by default only `fmt` is included |
| .ExecuteSimple(cmdString&nbsp;string) | Run golang source codes. The code will be placed inside virtual `main()` func |
| .Execute(cmdString string) | Run golang source codes which contains `main()` func |
| .ExecuteRaw(cmdString string) | Run complete golang source code |
| .DeleteTemporaryPath() | Force delete temporary path which used to do the exec process |


## Contribution

Feel free to contribute

`fork` -> `commit` -> `push` -> `pull request`

## License

MIT License

## Author

Noval Agung Prayogo - [http://novalagung.com/](http://novalagung.com)
