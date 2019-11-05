package main

import (
	"github.com/ZhengHe-MD/gorefactor"
	"github.com/dave/dst"
	"go/token"
	"log"
	"os"
)

func main() {
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

	func RegisterGuest(ctx context.Context, req *RegisterGuestReq) (r *RegisterGuestRes) {
		tctx := trace.NewThriftUtilContextFromContext(ctx)
		err := rpc(ctx, strconv.Itoa(rand.Int()), time.Millisecond*1000,
			func(c *AccountServiceClient) (er error) {
				r, er = c.RegisterGuest(req, tctx)
				return er
			},
		)
		if err != nil {
			r = &RegisterGuestRes{
				Errinfo: tiperr.NewInternalErr(-1001, fmt.Sprintf("call serice:%s proc:%s m:RegisterGuest err:%s", service_key, proc_thrift, err)),
			}
		}
		return
	}
	`

	//var expected = `
	//package main
	//
	//func rpc(ctx context.Context, hashKey string, timeout time.Duration, fn func(*AccountServiceClient) error) error {
	//	return clientThrift.RpcWithContextV2(ctx, hashKey, timeout, func(fctx context.Context, c interface{}) error {
	//		ct, ok := c.(*AccountServiceClient)
	//		if ok {
	//			return fn(fctx, ct)
	//		} else {
	//			return fmt.Errorf("reflect client thrift error")
	//		}
	//	})
	//}
	//
	//func RegisterGuest(ctx context.Context, req *RegisterGuestReq) (r *RegisterGuestRes) {
	//	err := rpc(ctx, strconv.Itoa(rand.Int()), time.Millisecond*1000,
	//		func(fctx context.Context, c *AccountServiceClient) (er error) {
	//			tctx := trace.NewThriftUtilContextFromContext(fctx)
	//			r, er = c.RegisterGuest(req, tctx)
	//			return er
	//		},
	//	)
	//	if err != nil {
	//		r = &RegisterGuestRes{
	//			Errinfo: tiperr.NewInternalErr(-1001, fmt.Sprintf("call serice:%s proc:%s m:RegisterGuest err:%s", service_key, proc_thrift, err)),
	//		}
	//	}
	//	return
	//}
	//`

	df, err := gorefactor.ParseSrcFileFromBytes([]byte(src))
	if err != nil {
		log.Println(err)
		return
	}

	gorefactor.SetMethodOnReceiver(df, "clientThrift", "RpcWithContext", "RpcWithContextV2")
	gorefactor.AddArgToCallExpr(df, "fn", dst.NewIdent("fctx"), 0)
	gorefactor.AddFieldToFuncLitParams(df, &dst.Field{
		Names: []*dst.Ident{dst.NewIdent("fctx")},
		Type:  dst.NewIdent("context.Context"),
	}, 0)
	gorefactor.DeleteStmtFromFuncBody(df, "RegisterGuest", &dst.AssignStmt{
		Lhs:  []dst.Expr{
			dst.NewIdent("tctx"),
		},
		Tok: token.DEFINE,
		Rhs:  []dst.Expr{
			&dst.CallExpr{
				Fun: &dst.SelectorExpr{
					X:    dst.NewIdent("trace"),
					Sel:  dst.NewIdent("NewThriftUtilContextFromContext"),
				},
				Args:    []dst.Expr{
					dst.NewIdent("ctx"),
				},
			},
		},
	})
	gorefactor.AddStmtToFuncLitBody(df, &dst.AssignStmt{
		Lhs:  []dst.Expr{
			dst.NewIdent("tctx"),
		},
		Tok: token.DEFINE,
		Rhs:  []dst.Expr{
			&dst.CallExpr{
				Fun: &dst.SelectorExpr{
					X:    dst.NewIdent("trace"),
					Sel:  dst.NewIdent("NewThriftUtilContextFromContext"),
				},
				Args:    []dst.Expr{
					dst.NewIdent("tctx"),
				},
			},
		},
	}, 0)
	gorefactor.DeleteStmtFromFuncBody(df, "rpc", &dst.AssignStmt{
		Lhs:  []dst.Expr{
			dst.NewIdent("tctx"),
		},
		Tok: token.DEFINE,
		Rhs:  []dst.Expr{
			&dst.CallExpr{
				Fun: &dst.SelectorExpr{
					X:    dst.NewIdent("trace"),
					Sel:  dst.NewIdent("NewThriftUtilContextFromContext"),
				},
				Args:    []dst.Expr{
					dst.NewIdent("tctx"),
				},
			},
		},
	})

	err = gorefactor.FprintFile(os.Stdout, df)
	if err != nil {
		log.Println(err)
		return
	}
}