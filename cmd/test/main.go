package test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/elliotchance/ok/compiler"
	"github.com/elliotchance/ok/util"
	"github.com/elliotchance/ok/vm"
)

type Command struct {
	// Verbose will print all test names.
	Verbose bool

	// Filter is a regexp based on the test name.
	Filter string
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Description is shown in "ok -help".
func (*Command) Description() string {
	return "run tests"
}

// Run is the entry point for the "ok test" command.
func (c *Command) Run(args []string) {
	flag.Usage = func() {
		fmt.Println(c.Description())
		fmt.Println("Usage:")
		fmt.Println("")
		fmt.Println("\tok test [args]")
		fmt.Println("")
		fmt.Println("the flags are:")
		fmt.Println("")
		fmt.Println("\t-f\t\tregular expression to filter tests by name")
		fmt.Println("\t-help\t\tprint this message")
		fmt.Println("\t-v\t\tprint all test names")
		fmt.Println("")
	}
	flag.StringVar(&c.Filter, "f", "", "")
	flag.BoolVar(&c.Verbose, "v", false, "")
	check(flag.CommandLine.Parse(args))
	args = flag.Args()

	if len(args) == 0 {
		args = []string{"."}
	}

	okPath, err := util.OKPath()
	check(err)

	for _, arg := range args {
		packageName := util.PackageNameFromPath(okPath, arg)
		if arg == "." {
			packageName = "."
		}
		anonFunctionName := 0
		f, _, errs := compiler.Compile(okPath, packageName, true,
			&anonFunctionName, false)
		util.CheckErrorsWithExit(errs)

		m := vm.NewVM("no-package")
		startTime := time.Now()
		check(m.LoadFile(f))
		err := m.RunTests(c.Verbose, regexp.MustCompile(c.Filter), packageName)
		elapsed := time.Since(startTime).Milliseconds()
		check(err)

		assertWord := pluralise("assert", m.TotalAssertions)
		if m.TestsFailed > 0 {
			fmt.Printf("%s: %d failed %d passed %d %s (%d ms)\n",
				packageName, m.TestsFailed, m.TestsPass,
				m.TotalAssertions, assertWord, elapsed)
		} else {
			fmt.Printf("%s: %d passed %d %s (%d ms)\n",
				packageName, m.TestsPass,
				m.TotalAssertions, assertWord, elapsed)
		}

		if m.TestsFailed > 0 {
			os.Exit(1)
		}
	}
}

func pluralise(word string, n int) string {
	if n == 1 {
		return word
	}

	return word + "s"
}
