package gorefactor

import (
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/goast"
	"github.com/dave/dst/decorator/resolver/guess"
	"go/token"
	"io"
	"io/ioutil"
	"os"
)

var defaultFileSet = token.NewFileSet()

// given go src file, in the form of bytes, and parse it into *dst.File
func ParseSrcFileFromBytes(src []byte) (df *dst.File, err error) {
	dec := decorator.NewDecoratorWithImports(
		defaultFileSet,
		"main",
		goast.WithResolver(guess.New()))
	return dec.Parse(src)
}

// given go src filename, in the form of valid path, and parse it into *dst.File
func ParseSrcFile(filename string) (df *dst.File, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	src, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	return ParseSrcFileFromBytes(src)
}

// given a io.Writer and *dst.File, write the *dst.File out to io.Writer
func FprintFile(out io.Writer, df *dst.File) error {
	dec := decorator.NewRestorerWithImports("main", guess.New())
	return dec.Fprint(out, df)
}
