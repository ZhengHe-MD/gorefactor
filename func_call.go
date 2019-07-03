package gorefactor

import (
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
)

func HasArgInCallExpr(df *dst.File, funcName string, arg dst.Expr) (ret bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.CallExpr:
			var found bool
			nn := node.(*dst.CallExpr)
			if ie, ok := nn.Fun.(*dst.Ident); ok && ie.Name == funcName {
				found = true
			}

			if se, ok := nn.Fun.(*dst.SelectorExpr); ok && se.Sel.Name == funcName {
				found = true
			}

			if found {
				for _, cArg := range nn.Args {
					if nodesEqual(arg, cArg) {
						ret = true
					}
				}
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

func DeleteArgFromCallExpr(df *dst.File, funcName string, arg dst.Expr) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.CallExpr:
			var found bool
			nn := node.(*dst.CallExpr)
			if ie, ok := nn.Fun.(*dst.Ident); ok && ie.Name == funcName {
				found = true
			}

			if se, ok := nn.Fun.(*dst.SelectorExpr); ok && se.Sel.Name == funcName {
				found = true
			}

			if found {
				var newArgs []dst.Expr
				for _, cArg := range nn.Args {
					if !nodesEqual(arg, cArg) {
						newArgs = append(newArgs, cArg)
					} else {
						modified = true
					}
				}
				nn.Args = newArgs
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

func AddArgToCallExpr(df *dst.File, funcName string, arg dst.Expr, pos int) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.CallExpr:
			nn := node.(*dst.CallExpr)

			var ce *dst.CallExpr

			if ie, ok := nn.Fun.(*dst.Ident); ok && ie.Name == funcName {
				ce = nn
			}

			if se, ok := nn.Fun.(*dst.SelectorExpr); ok && se.Sel.Name == funcName {
				ce = nn
			}

			if ce != nil {
				args := ce.Args
				pos = normalizePos(pos, len(args))
				ce.Args = append(
					args[:pos],
					append([]dst.Expr{dst.Clone(arg).(dst.Expr)}, args[pos:]...)...)
				modified = true
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}
