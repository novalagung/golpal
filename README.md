# Gopal

Easy to use Golang Exec Library


### Introduction

Gopal is simple library to allow developer exec-ing golang source codes.

Technically Golang doesn't provide API to do some exec, so we use temporary file to achieve that.

[![Build Status](https://travis-ci.org/novalagung/gopal.png?branch=master)](https://travis-ci.org/novalagung/gopal)
[![Version](http://img.shields.io/cocoapods/v/gopal.svg)](http://cocoadocs.org/docsets/gopal)

### Installation

Just go get the lib is enough

```bash
go get -u github.com/novalagung/gopal
```

### Example

##### Simple Example

```go
package main

import "github.com/novalagung/gopal"
import "fmt"

func main() {
	cmdString := `3 + 2`
	output, err := gopal.New().ExecuteSimple(cmdString)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("result", "=>", output)
}
```

##### Another Example

```go
cmdString := `
	number := 3
	if number == 2 {
		fmt.Println("wrong")
	} else {
		fmt.Println("right")
	}
`

output, err := gopal.New().ExecuteSimple(cmdString)
if err != nil {
	fmt.Println(err)
}
fmt.Println("result", "=>", output)
```

##### Example which use `strings` and `runtime`

```go
cmdString := `
	osName := runtime.GOOS
	arr := []string{"my", "operation system", "is", osName}
	return strings.Join(arr, ", ")
`

output, err := gopal.New().AddLibs("strings", "runtime").ExecuteSimple(cmdString)
if err != nil {
	fmt.Println(err)
}
fmt.Println("result", "=>", output)
```

##### Example Not Simple (again, NOT SIMPLE, by using `Execute()` func)

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

output, err := gopal.New().Execute(cmdString)
if err != nil {
	fmt.Println(err)
}
fmt.Println("result", "=>", output)
```

##### Example Execute Raw (use `ExecuteRaw()` func)

```go
cmdString := `
	package main

	import "fmt"

	func main() {
		fmt.Print("hello")
	}
`

output, err := gopal.New().ExecuteRaw(cmdString)
if err != nil {
	fmt.Println(err)
}
fmt.Println("result", "=>", output)
```

##### More Example

For more examples please take a look at the [`gopal_test.go` file](https://github.com/novalagung/gopal/blob/master/gopal_test.go).

### API Reference



### Contribution

Feel free to contribute

`fork` -> `commit` -> `push` -> `pull request`


### License

MIT License