package main

import (
	"github.com/quexer/utee"

	"liuyu/stu/pkg"
)

func main() {
	a, err := pkg.New()
	utee.Chk(err)

	a.Run()
}
