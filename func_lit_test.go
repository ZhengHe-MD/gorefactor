package gorefactor

import (
	"github.com/dave/dst"
	"github.com/stretchr/testify/assert"
	"testing"
)

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

		var expectedWithoutScope = `
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

		{
			df, _ := ParseSrcFileFromBytes([]byte(src))
			assert.True(t, AddFieldToFuncLitParams(df, EmptyScope, &dst.Field{
				Names: []*dst.Ident{dst.NewIdent("fctx")},
				Type:  dst.NewIdent("context.Context"),
			}, 0))
			buf := printToBuf(df)
			assertCodesEqual(t, expectedWithoutScope, buf.String())
		}

		{
			df, _ := ParseSrcFileFromBytes([]byte(src))
			assert.False(t, AddFieldToFuncLitParams(df, Scope{FuncName: "NonExist"}, &dst.Field{
				Names: []*dst.Ident{dst.NewIdent("fctx")},
				Type:  dst.NewIdent("context.Context"),
			}, 0))
			buf := printToBuf(df)
			assertCodesEqual(t, src, buf.String())
		}
	})
}