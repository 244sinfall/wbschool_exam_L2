package main

import (
	"flag"
	"fmt"
	"man_cut/pkg/cut"
	"os"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель (дефолт - \t)
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	f := cut.Fields{}
	var d string
	var s bool
	flag.Var(&f, "f", "fields to be on output. Use comma as delimiter")
	flag.StringVar(&d, "d", "\t", "set delimiter (default is \\t)")
	flag.BoolVar(&s, "s", false, "show only separated lines")
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		_, _ = os.Stderr.WriteString("you provided no file path\n")
		os.Exit(1)
	}
	file, err := os.Open(args[0])
	if err != nil {
		_, _ = os.Stderr.WriteString("could not open the file:" + err.Error() + "\n")
		os.Exit(1)
	}
	process := cut.Cut{
		Source: file,
		F:      f,
		D:      d,
		S:      s,
	}
	err = process.Print()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
