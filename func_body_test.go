package gorefactor

import (
	"fmt"
	"github.com/dave/dst"
	"github.com/stretchr/testify/assert"
	"go/token"
	"testing"
)

func TestHasStmtInsideFuncBody(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var src = `
		package main

		import "fmt"

		func main() {
			a := 1
			b := 1
			fmt.Println(a+b)
		}
		`

		cases := []struct {
			stmt     dst.Stmt
			expected bool
		}{
			{&dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent("a")},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{&dst.BasicLit{
					Kind:  token.INT,
					Value: "1",
				}},
			}, true},
			{&dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent("b")},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{&dst.BasicLit{
					Kind:  token.INT,
					Value: "1",
				}},
			}, true},
			{&dst.ExprStmt{
				X: &dst.CallExpr{
					Fun: &dst.Ident{
						Name: "Println",
						Path: "fmt",
					},
					Args: []dst.Expr{&dst.BinaryExpr{
						X:    dst.NewIdent("a"),
						Op:   token.ADD,
						Y:    dst.NewIdent("b"),
						Decs: dst.BinaryExprDecorations{},
					}},
				},
			}, true},
			{&dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent("a")},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{&dst.BasicLit{
					Kind:  token.INT,
					Value: "2",
				}},
			}, false},
		}

		df, _ := ParseSrcFileFromBytes([]byte(src))

		for _, c := range cases {
			assert.Equal(t, c.expected, HasStmtInsideFuncBody(df, "main", c.stmt))
		}
	})
}

func TestDeleteStmtFromFuncBody(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var src = `
		package main

		import (
			"fmt"
		)

		func main() {
			a := 1
			b := 1
			fmt.Println(a+b)
		}
		`

		var expectedTemplate = `
		package main
		
		func main() {
			%s
		}
		`

		cases := []struct {
			stmt             dst.Stmt
			expectedModified bool
			expectedBody     string
		}{
			{&dst.ExprStmt{
				X: &dst.CallExpr{
					Fun: &dst.Ident{
						Name: "Println",
						Path: "fmt",
					},
					Args: []dst.Expr{&dst.BinaryExpr{
						X:    dst.NewIdent("a"),
						Op:   token.ADD,
						Y:    dst.NewIdent("b"),
						Decs: dst.BinaryExprDecorations{},
					}},
				},
			}, true, "a := 1 \n b := 1"},
			{&dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent("b")},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{&dst.BasicLit{
					Kind:  token.INT,
					Value: "1",
				}},
			}, true, "a := 1"},
			{&dst.AssignStmt{
				Lhs: []dst.Expr{dst.NewIdent("a")},
				Tok: token.DEFINE,
				Rhs: []dst.Expr{&dst.BasicLit{
					Kind:  token.INT,
					Value: "1",
				}},
			}, true, ""},
		}

		df, _ := ParseSrcFileFromBytes([]byte(src))

		for _, c := range cases {
			assert.Equal(t, c.expectedModified, DeleteStmtFromFuncBody(df, "main", c.stmt))
			buf := printToBuf(df)
			assertCodesEqual(t, fmt.Sprintf(expectedTemplate, c.expectedBody), buf.String())
		}
	})

	t.Run("multiple", func(t *testing.T) {
		var src = `
		package main

		import (
			"fmt"
		)

		func main() {
			fmt.Println()
			fmt.Println()
		}
		`

		var expected = `
		package main

		func main() {
		}
		`

		printStmt := &dst.ExprStmt{
			X: &dst.CallExpr{
				Fun: &dst.Ident{
					Name: "Println",
					Path: "fmt",
				},
			},
		}

		df, _ := ParseSrcFileFromBytes([]byte(src))

		assert.Equal(t, true, DeleteStmtFromFuncBody(df, "main", printStmt))
		buf := printToBuf(df)
		assertCodesEqual(t, expected, buf.String())
	})
}

func TestDeleteSelectorExprFromFuncBody(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var src = `
		package main

		func A() {
			fmt.Printf("%s", "foo")
		}

		func main() {
			fmt.Println()
			fmt.Printf("%s", "bar")
			fmt.Printf("%s", "baz")
			fmt.Println()
		}
		`

		var expected = `
		package main

		func A() {
			fmt.Printf("%s", "foo")
		}

		func main() {
			fmt.Println()
			fmt.Println()
		}
		`

		printSelector := &dst.SelectorExpr{
			X:    dst.NewIdent("fmt"),
			Sel:  dst.NewIdent("Printf"),
		}

		df, _ := ParseSrcFileFromBytes([]byte(src))

		assert.Equal(t, true, DeleteSelectorExprFromFuncBody(df, "main", printSelector))
		buf := printToBuf(df)
		assertCodesEqual(t, expected, buf.String())
	})

	t.Run("complex", func(t *testing.T) {
		var src = `
		package main

		func (m *TestServiceImpl) Hello(req *HelloReq) (r *HelloRes, err error) {
			fun := "TestServiceImpl.Hello -->"

			st := stime.NewTimeStat()
			defer func() {
				dur := st.Duration()
				log.Infof("%s req:%v tm:%d", fun, req, dur)
				monitor.Stat("RPC-Hello", dur)
			}()

			return logic.HandleTest.Hello(req), nil
		}
		`

		var expected = `
		package main

		func (m *TestServiceImpl) Hello(req *HelloReq) (r *HelloRes, err error) {
			fun := "TestServiceImpl.Hello -->"

			st := stime.NewTimeStat()
			defer func() {
				dur := st.Duration()
				monitor.Stat("RPC-Hello", dur)
			}()

			return logic.HandleTest.Hello(req), nil
		}
		`

		printSelector := &dst.SelectorExpr{
			X:    dst.NewIdent("log"),
			Sel:  dst.NewIdent("Infof"),
		}

		df, _ := ParseSrcFileFromBytes([]byte(src))

		assert.Equal(t, true, DeleteSelectorExprFromFuncBody(df, "Hello", printSelector))
		buf := printToBuf(df)
		assertCodesEqual(t, expected, buf.String())
	})
}

func TestAddStmtToFuncBody(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var src = `
		package main

		func main() {
			a := 1
			b := 2
		}
		`

		var expectedTemplate = `
		package main

		func main() {
			%s
		}
		`

		cstmt := &dst.AssignStmt{
			Lhs: []dst.Expr{dst.NewIdent("c")},
			Tok: token.DEFINE,
			Rhs: []dst.Expr{&dst.BasicLit{
				Kind:  token.INT,
				Value: "3",
			}},
		}

		cases := []struct {
			pos          int
			expectedBody string
		}{
			{0, "c := 3\na := 1\nb := 2"},
			{1, "a := 1\nc := 3\nb := 2"},
			{2, "a := 1\nb := 2\nc := 3"},
			{-1, "a := 1\nb := 2\nc := 3"},
		}

		for _, c := range cases {
			df, _ := ParseSrcFileFromBytes([]byte(src))
			assert.True(t, AddStmtToFuncBody(df, "main", cstmt, c.pos))
			buf := printToBuf(df)
			assertCodesEqual(t, fmt.Sprintf(expectedTemplate, c.expectedBody), buf.String())
		}
	})
}

func TestAddStmtToFuncBodyRelativeTo(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var src = `
		package main
		
		func main() {
			a := 1
			b := 2
		}
		`

		var expectedTemplate = `
		package main

		func main() {
			%s
		}
		`

		astmt := &dst.AssignStmt{
			Lhs: []dst.Expr{dst.NewIdent("a")},
			Tok: token.DEFINE,
			Rhs: []dst.Expr{&dst.BasicLit{
				Kind:  token.INT,
				Value: "1",
			}},
		}

		bstmt := &dst.AssignStmt{
			Lhs: []dst.Expr{dst.NewIdent("b")},
			Tok: token.DEFINE,
			Rhs: []dst.Expr{&dst.BasicLit{
				Kind:  token.INT,
				Value: "2",
			}},
		}

		cstmt := &dst.AssignStmt{
			Lhs: []dst.Expr{dst.NewIdent("c")},
			Tok: token.DEFINE,
			Rhs: []dst.Expr{&dst.BasicLit{
				Kind:  token.INT,
				Value: "3",
			}},
		}

		cases := []struct {
			direction    int
			refStmt      dst.Stmt
			expectedBody string
		}{
			{relativeDirectionBefore, astmt, "c := 3\na := 1\nb := 2"},
			{relativeDirectionAfter, astmt, "a := 1\nc := 3\nb := 2"},
			{relativeDirectionBefore, bstmt, "a := 1\nc := 3\nb := 2"},
			{relativeDirectionAfter, bstmt, "a := 1\nb := 2\nc := 3"},
		}

		for _, c := range cases {
			df, _ := ParseSrcFileFromBytes([]byte(src))
			assert.True(t, addStmtToFuncBodyRelativeTo(df, "main", cstmt, c.refStmt, c.direction))
			buf := printToBuf(df)
			assertCodesEqual(t, fmt.Sprintf(expectedTemplate, c.expectedBody), buf.String())
		}
	})
}
