package gorefactor

import (
	"github.com/dave/dst"
	"github.com/dave/dst/dstutil"
)

func HasStmtInsideFuncBody(df *dst.File, funcName string, stmt dst.Stmt) (ret bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				funcBody := nn.Body
				for _, ss := range funcBody.List {
					if nodesEqual(ss, stmt) {
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

func DeleteStmtFromFuncBody(df *dst.File, funcName string, stmt dst.Stmt) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				funcBody := nn.Body

				var newList []dst.Stmt
				for _, ss := range funcBody.List {
					if !nodesEqual(ss, stmt) {
						newList = append(newList, ss)
					} else {
						modified = true
					}
				}

				funcBody.List = newList
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

func AddStmtToFuncBody(df *dst.File, funcName string, stmt dst.Stmt, pos int) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				stmtList := nn.Body.List
				pos = normalizePos(pos, len(stmtList))

				nn.Body.List = append(
					stmtList[:pos],
					append([]dst.Stmt{dst.Clone(stmt).(dst.Stmt)}, stmtList[pos:]...)...)
				modified = true
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

func AddStmtToFuncBodyStart(df *dst.File, funcName string, stmt dst.Stmt) (modified bool) {
	return AddStmtToFuncBody(df, funcName, stmt, 0)
}

func AddStmtToFuncBodyEnd(df *dst.File, funcName string, stmt dst.Stmt) (modified bool) {
	return AddStmtToFuncBody(df, funcName, stmt, -1)
}

const (
	relativeDirectionBefore = iota
	relativeDirectionAfter
)

func addStmtToFuncBodyRelativeTo(df *dst.File, funcName string, stmt, refStmt dst.Stmt, relDirection int) (modified bool) {
	pre := func(c *dstutil.Cursor) bool {
		node := c.Node()

		switch node.(type) {
		case *dst.FuncDecl:
			if nn := node.(*dst.FuncDecl); nn.Name.Name == funcName {
				funcBody := nn.Body
				var newStmtList []dst.Stmt
				for _, ss := range funcBody.List {
					if nodesEqual(ss, refStmt) {
						switch relDirection {
						case relativeDirectionBefore:
							newStmtList = append(newStmtList, stmt, ss)
							modified = true
						case relativeDirectionAfter:
							newStmtList = append(newStmtList, ss, stmt)
							modified = true
						}
					} else {
						newStmtList = append(newStmtList, ss)
					}
				}
				funcBody.List = newStmtList
				return false
			}
		}
		return true
	}

	dstutil.Apply(df, pre, nil)
	return
}

func AddNodeToFuncBodyBefore(df *dst.File, funcName string, stmt, refStmt dst.Stmt) (modified bool) {
	return addStmtToFuncBodyRelativeTo(df, funcName, stmt, refStmt, relativeDirectionBefore)
}

func AddNodeToFuncBodyAfter(df *dst.File, funcName string, stmt, refStmt dst.Stmt) (modified bool) {
	return addStmtToFuncBodyRelativeTo(df, funcName, stmt, refStmt, relativeDirectionAfter)
}
