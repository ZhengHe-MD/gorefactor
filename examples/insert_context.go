package main

import (
	"github.com/ZhengHe-MD/gorefactor"
	"github.com/dave/dst"
	"log"
	"os"
)

/*
say we want to refactor the following code

==================================
package main

func f() {}

func main() {
    f()
    f()
}
==================================

to

==================================
package main

import "context"

func f(ctx context.Context) {}

func main() {
    f(context.TODO())
    f(context.TODO())
}
==================================

*/

func main() {
	var src = `
    	package main

		func f() {}

		func main() {
			f()
			f()
		}
		`

	df, err := gorefactor.ParseSrcFileFromBytes([]byte(src))
	if err != nil {
		log.Println(err)
		return
	}

	gorefactor.AddFieldToFuncDeclParams(df, "f", &dst.Field{
		Names: []*dst.Ident{dst.NewIdent("ctx")},
		Type: &dst.Ident{
			Name: "Context",
			Path: "context",
		},
	}, 0)

	gorefactor.AddArgToCallExpr(df, gorefactor.EmptyScope, "f", &dst.CallExpr{
		Fun: &dst.Ident{
			Name: "TODO",
			Path: "context",
		},
	}, 0)

	err = gorefactor.FprintFile(os.Stdout, df)
	if err != nil {
		log.Println(err)
		return
	}
}
