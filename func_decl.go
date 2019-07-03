package gorefactor

import (
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
)

func HasFieldInFuncDeclParams(df *dst.File, funcName string, field *dst.Field) (ret bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				funcType := nn.Type
				for _, ff := range funcType.Params.List {
					if nodesEqual(ff, field) {
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

func DeleteFieldFromFuncDeclParams(df *dst.File, funcName string, field *dst.Field) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				funcType := nn.Type
				var newList []*dst.Field
				for _, ff := range funcType.Params.List {
					if !nodesEqual(ff, field) {
						newList = append(newList, ff)
					} else {
						modified = true
					}
				}
				funcType.Params.List = newList
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

func AddFieldToFuncDeclParams(df *dst.File, funcName string, field *dst.Field, pos int) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				funcType := nn.Type
				fieldList := funcType.Params.List
				pos = normalizePos(pos, len(fieldList))
				funcType.Params.List = append(
					fieldList[:pos],
					append([]*dst.Field{field}, fieldList[pos:]...)...)
				modified = true
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}
