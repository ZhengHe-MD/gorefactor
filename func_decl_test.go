package gorefactor

import (
	"bytes"
	"fmt"
	"github.com/dave/dst"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHasFieldInFuncDeclParams(t *testing.T) {
	t.Run("normal fields", func(t *testing.T) {
		var src = `
		package main

		import "context"

		type A struct {}

		func f(ctx context.Context, a int, b float, c bool, d string, e A, f interface{}) {}

		func main() {}
		`

		df, _ := ParseSrcFileFromBytes([]byte(src))

		cases := []struct {
			field    *dst.Field
			expected bool
		}{
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("a")},
				Type:  dst.NewIdent("int"),
			}, true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("b")},
				Type:  dst.NewIdent("float"),
			}, true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("c")},
				Type:  dst.NewIdent("bool"),
			}, true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("d")},
				Type:  dst.NewIdent("string"),
			}, true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("e")},
				Type:  dst.NewIdent("A"),
			}, true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("f")},
				Type:  &dst.InterfaceType{},
			}, true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("ctx")},
				Type: &dst.Ident{
					Name: "Context",
					Path: "context",
				},
			}, true},
		}

		for _, c := range cases {
			assert.Equal(t, c.expected, HasFieldInFuncDeclParams(df, "f", c.field))
		}
	})

	t.Run("decl inside decl", func(t *testing.T) {
		var src = `
		package main

		func f(a int, b int, compare func(c, d int) bool) bool { return true }
		`

		df, _ := ParseSrcFileFromBytes([]byte(src))

		cases := []struct {
			field    *dst.Field
			expected bool
		}{
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("a")},
				Type:  dst.NewIdent("int"),
			}, true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("c")},
				Type:  dst.NewIdent("int"),
			}, false},
		}

		for _, c := range cases {
			assert.Equal(t, c.expected, HasFieldInFuncDeclParams(df, "f", c.field))
		}
	})
}

func TestDeleteFieldFromFuncDeclParams(t *testing.T) {

	t.Run("normal fields", func(t *testing.T) {
		var src = `
		package main

		import "context"

		type A struct {}

		func f(ctx context.Context, a int, b float, c bool, d string, e A, f interface{}) {}

		func main() {}
		`

		var expectedTemplate = `
		package main

		import "context"

		type A struct {}

		func f(%s) {}

		func main() {}
		`

		df, _ := ParseSrcFileFromBytes([]byte(src))

		cases := []struct {
			field            *dst.Field
			expectedFargs    string
			expectedModified bool
		}{
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("a")},
				Type:  dst.NewIdent("int"),
			}, "ctx context.Context, b float, c bool, d string, e A, f interface{}", true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("b")},
				Type:  dst.NewIdent("float"),
			}, "ctx context.Context, c bool, d string, e A, f interface{}", true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("c")},
				Type:  dst.NewIdent("bool"),
			}, "ctx context.Context, d string, e A, f interface{}", true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("d")},
				Type:  dst.NewIdent("string"),
			}, "ctx context.Context, e A, f interface{}", true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("e")},
				Type:  dst.NewIdent("A"),
			}, "ctx context.Context, f interface{}", true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("f")},
				Type:  &dst.InterfaceType{},
			}, "ctx context.Context", true},
		}

		for _, c := range cases {
			var buf *bytes.Buffer
			assert.True(t, c.expectedModified, DeleteFieldFromFuncDeclParams(df, "f", c.field))
			buf = printToBuf(df)
			assertCodesEqual(t, fmt.Sprintf(expectedTemplate, c.expectedFargs), buf.String())
		}

		var finalExpect = `
		package main

		type A struct {}

		func f() {}

		func main() {}
		`

		var buf *bytes.Buffer
		assert.True(t, true, DeleteFieldFromFuncDeclParams(df, "f", &dst.Field{
			Names: []*dst.Ident{dst.NewIdent("ctx")},
			Type: &dst.Ident{
				Name: "Context",
				Path: "context",
			},
		}))
		buf = printToBuf(df)
		assertCodesEqual(t, finalExpect, buf.String())

	})

	t.Run("decl inside decl", func(t *testing.T) {
		var src = `
		package main

		func f(a int, b int, compare func(b int, d int) bool) bool { return true }
		`

		var expectedTemplate = `
		package main

		func f(%s) bool { return true }
		`

		cases := []struct {
			field            *dst.Field
			expectedArgs     string
			expectedModified bool
		}{
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("a")},
				Type:  dst.NewIdent("int"),
			}, "b int, compare func(b int, d int) bool", true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("b")},
				Type:  dst.NewIdent("int"),
			}, "a int, compare func(b int, d int) bool", true},
			{&dst.Field{
				Names: []*dst.Ident{dst.NewIdent("d")},
				Type:  dst.NewIdent("int"),
			}, "a int, b int, compare func(b int, d int) bool", false},
		}

		for _, c := range cases {
			var buf *bytes.Buffer
			df, _ := ParseSrcFileFromBytes([]byte(src))
			assert.Equal(t, c.expectedModified, DeleteFieldFromFuncDeclParams(df, "f", c.field))
			buf = printToBuf(df)
			assertCodesEqual(t, fmt.Sprintf(expectedTemplate, c.expectedArgs), buf.String())
		}
	})

	t.Run("decl inside decl2", func(t *testing.T) {
		var src = `
		package main

		func rpc(ctx context.Context, hashKey string, timeout time.Duration, fn func(*AccountServiceClient) error) error {
		return clientThrift.RpcWithContext(ctx, hashKey, timeout, func(c interface{}) error {
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

		func rpc(ctx context.Context, hashKey string, timeout time.Duration) error {
		return clientThrift.RpcWithContext(ctx, hashKey, timeout, func(c interface{}) error {
				ct, ok := c.(*AccountServiceClient)
				if ok {
					return fn(ct)
				} else {
					return fmt.Errorf("reflect client thrift error")
				}
			})
		}
		`

		df, _ := ParseSrcFileFromBytes([]byte(src))
		fn := &dst.Field{
			Names: []*dst.Ident{dst.NewIdent("fn")},
			Type: &dst.FuncType{
				Params: &dst.FieldList{
					List: []*dst.Field{
						{
							Type: &dst.StarExpr{
								X: dst.NewIdent("AccountServiceClient"),
							},
						},
					},
				},
				Results: &dst.FieldList{
					List: []*dst.Field{
						{
							Type: dst.NewIdent("error"),
						},
					},
				},
			},
		}
		modified := DeleteFieldFromFuncDeclParams(df, "rpc", fn)
		assert.Equal(t, true, modified)
		var buf *bytes.Buffer
		buf = printToBuf(df)
		assertCodesEqual(t, expected, buf.String())
	})
}

func TestAddFieldToFuncDeclParams(t *testing.T) {
	t.Run("add field to empty decl", func(t *testing.T) {
		var src = `
		package main

		func f() {}
		`

		var expected = `
		package main
		func f(i int) {}
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

		intField := &dst.Field{
			Names: []*dst.Ident{dst.NewIdent("i")},
			Type:  dst.NewIdent("int"),
		}

		for _, c := range cases {
			df, _ := ParseSrcFileFromBytes([]byte(src))
			var buf *bytes.Buffer
			assert.True(t, AddFieldToFuncDeclParams(df, "f", intField, c.pos))
			buf = printToBuf(df)
			assertCodesEqual(t, expected, buf.String())
		}
	})

	t.Run("add field to pos", func(t *testing.T) {
		var src = `
		package main

		func f(b float, c bool, d string) {}
		`

		var expectedTemplate = `
		package main

		func f(%s) {}
		`

		intField := &dst.Field{
			Names: []*dst.Ident{dst.NewIdent("a")},
			Type:  dst.NewIdent("int"),
		}

		cases := []struct {
			expectedFields string
			pos            int
		}{
			{"a int, b float, c bool, d string", 0},
			{"b float, a int, c bool, d string", 1},
			{"b float, c bool, a int, d string", 2},
			{"b float, c bool, d string, a int", 3},
			{"b float, c bool, d string, a int", -1},
		}

		for _, c := range cases {
			var buf *bytes.Buffer
			expected := fmt.Sprintf(expectedTemplate, c.expectedFields)
			df, _ := ParseSrcFileFromBytes([]byte(src))
			assert.True(t, AddFieldToFuncDeclParams(df, "f", intField, c.pos))
			buf = printToBuf(df)
			assertCodesEqual(t, expected, buf.String())
		}
	})
}

func TestSetMethodOnReceiver(t *testing.T) {
	t.Run("case1: fmt", func(t *testing.T) {
		var src = `
		package main

		import "fmt"
		
		func main() {
			fmt.Println()
		}
		`

		var expected = `
		package main

		import "fmt"

		func main() {
			fmt.Panicln()
		}
		`

		df, _ := ParseSrcFileFromBytes([]byte(src))
		var buf *bytes.Buffer
		assert.True(t, SetMethodOnReceiver(df, "fmt", "Println", "Panicln"))
		buf = printToBuf(df)
		assertCodesEqual(t, expected, buf.String())
	})

	t.Run("case2", func(t *testing.T) {
		var src = `
		package main

		func rpc(ctx context.Context, hashKey string, timeout time.Duration, fn func(*AccountServiceClient) error) error {
			return clientThrift.RpcWithContext(ctx, hashKey, timeout, func(c interface{}) error {
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
					return fn(ct)
				} else {
					return fmt.Errorf("reflect client thrift error")
				}
			})
		}
		`

		df, _ := ParseSrcFileFromBytes([]byte(src))
		var buf *bytes.Buffer
		assert.True(t, SetMethodOnReceiver(df, "clientThrift", "RpcWithContext", "RpcWithContextV2"))
		buf = printToBuf(df)
		assertCodesEqual(t, expected, buf.String())
	})
}

func TestAddFieldToFuncLitParams(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		var src = `
		package main

		func rpc(ctx context.Context, hashKey string, timeout time.Duration, fn func(*AccountServiceClient) error) error {
			return clientThrift.RpcWithContext(ctx, hashKey, timeout, func(c interface{}) error {
				ct, ok := c.(*AccountServiceClient)
				if ok {
					return fn(ct)
				} else {
					return fmt.Errorf("reflect client thrift error")
				}
			})
		}

		func UpdateUserStrInfo(ctx context.Context, req *UserStrInfoRequest) (r *UpdateUserInfoRes, err error) {
			tctx := trace.NewThriftUtilContextFromContext(ctx)
			if req == nil {
				return nil, fmt.Errorf("req nil")
			}
			err = rpc(ctx, req.Item, time.Millisecond*3000,
				func(c *AccountServiceClient) error {
					r, err = c.UpdateUserStrInfo(req, tctx)
					return err
				},
			)
			return
		}
		`

		var expected = `
		package main

		func rpc(ctx context.Context, hashKey string, timeout time.Duration, fn func(*AccountServiceClient) error) error {
			return clientThrift.RpcWithContext(ctx, hashKey, timeout, func(fctx context.Context, c interface{}) error {
				ct, ok := c.(*AccountServiceClient)
				if ok {
					return fn(ct)
				} else {
					return fmt.Errorf("reflect client thrift error")
				}
			})
		}
		
		func UpdateUserStrInfo(ctx context.Context, req *UserStrInfoRequest) (r *UpdateUserInfoRes, err error) {
			tctx := trace.NewThriftUtilContextFromContext(ctx)
			if req == nil {
				return nil, fmt.Errorf("req nil")
			}
			err = rpc(ctx, req.Item, time.Millisecond*3000,
				func(fctx context.Context, c *AccountServiceClient) error {
					r, err = c.UpdateUserStrInfo(req, tctx)
					return err
				},
			)
			return
		}
		`

		var buf *bytes.Buffer
		df, _ := ParseSrcFileFromBytes([]byte(src))
		assert.True(t, AddFieldToFuncLitParams(df, &dst.Field{
			Names: []*dst.Ident{dst.NewIdent("fctx")},
			Type:  dst.NewIdent("context.Context"),
		}, 0))
		buf = printToBuf(df)
		assertCodesEqual(t, expected, buf.String())
	})
}
