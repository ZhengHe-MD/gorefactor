package gorefactor

import (
	"bytes"
	"fmt"
	"github.com/dave/dst"
	"github.com/stretchr/testify/assert"
	"go/token"
	"strings"
	"testing"
)

func assertCodesEqual(t *testing.T, a, b string) {
	aa := strings.Join(strings.Fields(a), "")
	bb := strings.Join(strings.Fields(b), "")
	if aa != bb {
		t.Errorf("%s should be equal to %s", b, a)
	}
}

func printToBuf(df *dst.File) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	_ = FprintFile(buf, df)
	return buf
}

func TestHasArgInCallExpr(t *testing.T) {
	//t.Run("basic literals", func(t *testing.T) {
	//	var src = `
	//	package main
	//
	//	import (
	//		"fmt"
	//	)
	//
	//	func main() {
	//		fmt.Println(1, 1.1, true, "hello")
	//	}
	//	`
	//
	//	df, _ := ParseSrcFileFromBytes([]byte(src))
	//
	//	cases := []struct {
	//		arg      dst.Expr
	//		expected bool
	//	}{
	//		{&dst.BasicLit{Kind: token.INT, Value: "1"}, true},
	//		{&dst.BasicLit{Kind: token.INT, Value: "0"}, false},
	//		{&dst.BasicLit{Kind: token.FLOAT, Value: "1.1"}, true},
	//		{&dst.BasicLit{Kind: token.FLOAT, Value: "0.1"}, false},
	//		{&dst.BasicLit{Kind: token.STRING, Value: "\"hello\""}, true},
	//		{&dst.BasicLit{Kind: token.STRING, Value: "\"world\""}, false},
	//	}
	//
	//	for _, c := range cases {
	//		assert.Equal(t, c.expected, HasArgInCallExpr(df, EmptyScope, "Println", c.arg))
	//	}
	//})
	//
	//t.Run("variables", func(t *testing.T) {
	//	var src = `
	//	package main
	//
	//	import (
	//		"fmt"
	//	)
	//
	//	func main() {
	//		i := 1
	//		s := "hello"
	//		fmt.Println(i, s)
	//	}
	//	`
	//
	//	df, _ := ParseSrcFileFromBytes([]byte(src))
	//
	//	cases := []struct {
	//		arg      dst.Expr
	//		expected bool
	//	}{
	//		{dst.NewIdent("i"), true},
	//		{dst.NewIdent("s"), true},
	//		{dst.NewIdent("t"), false},
	//	}
	//
	//	for _, c := range cases {
	//		assert.Equal(t, c.expected, HasArgInCallExpr(df, EmptyScope, "Println", c.arg))
	//	}
	//})
	//
	//t.Run("structs", func(t *testing.T) {
	//	var src = `
	//	package main
	//
	//	import (
	//		"fmt"
	//	)
	//
	//	type A struct {}
	//
	//	func main() {
	//		fmt.Println(&A{})
	//	}
	//	`
	//
	//	df, _ := ParseSrcFileFromBytes([]byte(src))
	//
	//	cases := []struct {
	//		arg      dst.Expr
	//		expected bool
	//	}{
	//		{dst.NewIdent("A"), false},
	//		{&dst.UnaryExpr{
	//			Op: token.AND,
	//			X: &dst.CompositeLit{
	//				Type: dst.NewIdent("A"),
	//			},
	//		}, true},
	//	}
	//
	//	for _, c := range cases {
	//		assert.Equal(t, c.expected, HasArgInCallExpr(df, EmptyScope, "Println", c.arg))
	//	}
	//})
	//
	//t.Run("unresolved imports", func(t *testing.T) {
	//	var src = `
	//	package main
	//
	//	func main() {
	//		a := "hello"
	//		fmt.Println(1, a)
	//	}
	//	`
	//
	//	df, _ := ParseSrcFileFromBytes([]byte(src))
	//
	//	cases := []struct {
	//		arg      dst.Expr
	//		expected bool
	//	}{
	//		{&dst.BasicLit{Kind: token.INT, Value: "1"}, true},
	//		{dst.NewIdent("a"), true},
	//	}
	//
	//	for _, c := range cases {
	//		assert.Equal(t, c.expected, HasArgInCallExpr(df, EmptyScope, "Println", c.arg))
	//	}
	//})
	//
	//t.Run("selector caller", func(t *testing.T) {
	//	var src = `
	//	package main
	//
	//	type A struct {}
	//	func (a A) hello(i int) {}
	//
	//	func main() {
	//		var a A
	//		a.hello(1)
	//	}
	//	`
	//
	//	df, _ := ParseSrcFileFromBytes([]byte(src))
	//
	//	cases := []struct {
	//		arg      dst.Expr
	//		expected bool
	//	}{
	//		{&dst.BasicLit{Kind: token.INT, Value: "1"}, true},
	//		{dst.NewIdent("a"), false},
	//	}
	//
	//	for _, c := range cases {
	//		assert.Equal(t, c.expected, HasArgInCallExpr(df, EmptyScope, "hello", c.arg))
	//	}
	//})

	t.Run("with scope", func(t *testing.T) {
		var src = `
		package main
		
		func A() {
			fmt.Println("hello world")
		}
		func B() {
			fmt.Println(1)
        }
		`

		df, _ := ParseSrcFileFromBytes([]byte(src))
		cases := []struct {
			arg      dst.Expr
			scope    Scope
			expected bool
		}{
			{
				&dst.BasicLit{Kind: token.STRING, Value: "\"hello world\""},
				Scope{FuncName: "A"},
				true,
			},
			{
				&dst.BasicLit{Kind: token.STRING, Value: "\"hello world\""},
				Scope{FuncName: "B"},
				false,
			},
			{
				&dst.BasicLit{Kind: token.INT, Value: "1"},
				Scope{FuncName: "A"},
				false,
			},
			{
				&dst.BasicLit{Kind: token.INT, Value: "1"},
				Scope{FuncName: "B"},
				true,
			},
		}

		for _, c := range cases {
			assert.Equal(t, c.expected, HasArgInCallExpr(df, c.scope, "Println", c.arg))
		}
	})
}

func TestDeleteArgFromCallExpr(t *testing.T) {
	t.Run("delete multiple args", func(t *testing.T) {
		var src = `
		package main

		import (
			"fmt"
		)

		func main() {
			fmt.Println(1, 2, 3)
		}
		`

		var expected1 = `
		package main

		import (
			"fmt"
		)

		func main() {
			fmt.Println(1, 2)
		}
		`

		var expected2 = `
		package main

		import (
			"fmt"
		)

		func main() {
			fmt.Println(1)
		}
		`

		var expected3 = `
		package main

		import (
			"fmt"
		)

		func main() {
			fmt.Println()
		}
		`

		df, _ := ParseSrcFileFromBytes([]byte(src))

		var buf *bytes.Buffer

		assert.True(t, DeleteArgFromCallExpr(df, "Println", &dst.BasicLit{Kind: token.INT, Value: "3"}))
		buf = printToBuf(df)
		assertCodesEqual(t, expected1, buf.String())
		assert.False(t, DeleteArgFromCallExpr(df, "Println", &dst.BasicLit{Kind: token.INT, Value: "3"}))
		buf = printToBuf(df)
		assertCodesEqual(t, expected1, buf.String())
		assert.True(t, DeleteArgFromCallExpr(df, "Println", &dst.BasicLit{Kind: token.INT, Value: "2"}))
		buf = printToBuf(df)
		assertCodesEqual(t, expected2, buf.String())
		assert.True(t, DeleteArgFromCallExpr(df, "Println", &dst.BasicLit{Kind: token.INT, Value: "1"}))
		buf = printToBuf(df)
		assertCodesEqual(t, expected3, buf.String())
	})

	t.Run("delete args from selector caller", func(t *testing.T) {
		var src = `
		package main

		type A struct {}
		func (a A) hello(ints... int) {}

		func main() {
			var a A
			a.hello(1)
		}
		`

		var expected = `
		package main

		type A struct {}
		func (a A) hello(ints... int) {}

		func main() {
			var a A
			a.hello()
		}
		`

		df, _ := ParseSrcFileFromBytes([]byte(src))

		var buf *bytes.Buffer
		assert.True(t, DeleteArgFromCallExpr(df, "hello", &dst.BasicLit{Kind: token.INT, Value: "1"}))
		buf = printToBuf(df)
		assertCodesEqual(t, expected, buf.String())
	})

	t.Run("delete args from multiple calls", func(t *testing.T) {
		var src = `
		package main

		type A struct {}
		func (a A) hello(ints... int) {}

		func main() {
			var a A
			a.hello(1, 2)
		    a.hello(2, 3)
			a.hello(3, 1)
		}
		`

		var expected = `
		package main

		type A struct {}
		func (a A) hello(ints... int) {}

		func main() {
			var a A
			a.hello(2)
		    a.hello(2, 3)
			a.hello(3)
		}
		`

		df, _ := ParseSrcFileFromBytes([]byte(src))

		var buf *bytes.Buffer
		assert.True(t, DeleteArgFromCallExpr(df, "hello", &dst.BasicLit{Kind: token.INT, Value: "1"}))
		buf = printToBuf(df)
		assertCodesEqual(t, expected, buf.String())
	})

	t.Run("only delete outside arg", func(t *testing.T) {
		var src = `
		package main

		import (
			"fmt"
		)
		
		type A struct {
			i int
		}

		func main() {
			fmt.Println(&A{1}, 1)
			fmt.Println(1, &A{1})
		}
		`

		var expected = `
		package main

		import (
			"fmt"
		)
		
		type A struct {
			i int
		}

		func main() {
			fmt.Println(&A{1})
			fmt.Println(&A{1})
		}
		`

		df, _ := ParseSrcFileFromBytes([]byte(src))

		var buf *bytes.Buffer
		assert.True(t, DeleteArgFromCallExpr(df, "Println", &dst.BasicLit{Kind: token.INT, Value: "1"}))
		buf = printToBuf(df)
		assertCodesEqual(t, expected, buf.String())
	})
}

func TestAddArgToCallExpr(t *testing.T) {
	t.Run("add param to empty call", func(t *testing.T) {
		var src = `
		package main
		
		import (
			"fmt"
		)

		func main() {
			fmt.Println()
		}
		`

		var expected = `
		package main
		
		import (
			"fmt"
		)

		func main() {
			fmt.Println(1)
		}
		`

		cases := []struct {
			pos int
		}{
			{-1},
			{0},
			{1},
			{2},
			{3},
		}

		for _, c := range cases {
			df, _ := ParseSrcFileFromBytes([]byte(src))

			var buf *bytes.Buffer
			assert.True(t, AddArgToCallExpr(df, "Println", &dst.BasicLit{Kind: token.INT, Value: "1"}, c.pos))
			buf = printToBuf(df)
			assertCodesEqual(t, expected, buf.String())
		}
	})

	t.Run("add param to pos", func(t *testing.T) {
		var srcTemplate = `
		package main
		
		import (
			"fmt"
		)

		func main() {
			fmt.Println(%s)
		}
		`

		var expectedTemplate = `
		package main
		
		import (
			"fmt"
		)

		func main() {
			fmt.Println(%s)
		}
		`

		cases := []struct {
			srcArgs      string
			expectedArgs string
			arg          dst.Expr
			pos          int
		}{
			{"1, 2, 3", "0, 1, 2, 3", &dst.BasicLit{Kind: token.INT, Value: "0"}, 0},
			{"1, 2, 3", "1, 0, 2, 3", &dst.BasicLit{Kind: token.INT, Value: "0"}, 1},
			{"1, 2, 3", "1, 2, 0, 3", &dst.BasicLit{Kind: token.INT, Value: "0"}, 2},
			{"1, 2, 3", "1, 2, 3, 0", &dst.BasicLit{Kind: token.INT, Value: "0"}, 3},
			{"1, 2, 3", "1, 2, 3, 0", &dst.BasicLit{Kind: token.INT, Value: "0"}, -1},
		}

		for _, c := range cases {
			var buf *bytes.Buffer
			src := fmt.Sprintf(srcTemplate, c.srcArgs)
			expected := fmt.Sprintf(expectedTemplate, c.expectedArgs)

			df, _ := ParseSrcFileFromBytes([]byte(src))
			assert.True(t, AddArgToCallExpr(df, "Println", c.arg, c.pos))
			buf = printToBuf(df)
			assertCodesEqual(t, expected, buf.String())
		}
	})

	t.Run("structs", func(t *testing.T) {
		var srcTemplate = `
		package main

		import (
			"fmt"
		)

		type A struct {}

		func main() {
			fmt.Println(%s)
		}
		`

		var expectedTemplate = `
		package main

		import (
			"fmt"
		)

		type A struct {}

		func main() {
			fmt.Println(%s)
		}
		`

		cases := []struct {
			srcArgs      string
			expectedArgs string
			arg          dst.Expr
			pos          int
		}{
			{"&A{}", "1, &A{}", &dst.BasicLit{Kind: token.INT, Value: "1"}, 0},
			{"&A{}", "&A{}, 1", &dst.BasicLit{Kind: token.INT, Value: "1"}, -1},
		}

		for _, c := range cases {
			var buf *bytes.Buffer
			src := fmt.Sprintf(srcTemplate, c.srcArgs)
			expected := fmt.Sprintf(expectedTemplate, c.expectedArgs)

			df, _ := ParseSrcFileFromBytes([]byte(src))
			assert.True(t, AddArgToCallExpr(df, "Println", c.arg, c.pos))
			buf = printToBuf(df)
			assertCodesEqual(t, expected, buf.String())
		}
	})

	t.Run("complex case 1", func(t *testing.T) {
		var src = `
		package main

		func rpc(ctx context.Context, hashKey string, timeout time.Duration, fn func(*AccountServiceClient) error) error {
			return clientThrift.RpcWithContextV2(ctx, hashKey, timeout, func(c interface{}) error {
				ct, ok := c.(*AccountServiceClient)
				if ok {
					return fn(ct)
				} else {
					return fmt.Errorf("reflect client thrift error")
				}
			})
		}
 		`

		var expected = `
		package main

		func rpc(ctx context.Context, hashKey string, timeout time.Duration, fn func(*AccountServiceClient) error) error {
			return clientThrift.RpcWithContextV2(ctx, hashKey, timeout, func(c interface{}) error {
				ct, ok := c.(*AccountServiceClient)
				if ok {
					return fn(fctx, ct)
				} else {
					return fmt.Errorf("reflect client thrift error")
				}
			})
		}
		`

		df, _ := ParseSrcFileFromBytes([]byte(src))
		var buf *bytes.Buffer
		assert.True(t, true, AddArgToCallExpr(df, "fn", dst.NewIdent("fctx"), 0))
		buf = printToBuf(df)
		assertCodesEqual(t, expected, buf.String())
	})
}
