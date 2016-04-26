/*
 * Golpal - Easy to use Golang Eval Library
 * Created by Noval Agung Prayogo <caknopal@gmail.com>
 * http://novalagung.com/
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/novalagung/golpal"
)

const (
	ModeSimple string = "simple"
	ModeNormal string = "normal"
	ModeRaw    string = "raw"
)

func main() {
	flagFile := flag.String("file", "", "location file to eval")
	flagContent := flag.String("content", "", "content to eval")
	flagMode := flag.String("mode", "simple", "eval mode (simple, normal, raw). default is simple")
	flagLibs := flag.String("libs", "fmt", `include packages. ex: -libs="fmt, io/ioutil, bytes"`)

	flag.Parse()

	file := strings.TrimSpace(*flagFile)
	content := strings.TrimSpace(*flagContent)
	mode := strings.TrimSpace(*flagMode)
	libs := []string{}

	for _, lib := range strings.Split(*flagLibs, ",") {
		lib = strings.TrimSpace(lib)
		libs = append(libs, lib)
	}

	var gpl = golpal.New().AddLibs(libs...)
	var output string
	var err error

	if file == "" && content == "" {
		fmt.Fprintln(os.Stderr, "-content or -file is required")
		return
	}

	if file != "" {
		byts, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		content = strings.TrimSpace(string(byts))
	}

	switch mode {
	case ModeSimple:
		output, err = gpl.ExecuteSimple(content)
	case ModeNormal:
		output, err = gpl.Execute(content)
	case ModeRaw:
		output, err = gpl.ExecuteRaw(content)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fmt.Fprintln(os.Stdout, output)
}
