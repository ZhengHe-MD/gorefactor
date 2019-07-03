# GoRefactor

this module provides some basic utilities for you to do code migration on large code base written in golang.

[![Build Status](https://travis-ci.org/ZhengHe-MD/gorefactor.svg?branch=master)](https://travis-ci.org/ZhengHe-MD/gorefactor)
[![Go Report Card](https://goreportcard.com/badge/github.com/ZhengHe-MD/gorefactor)](https://goreportcard.com/report/github.com/ZhengHe-MD/gorefactor)
[![Coverage Status](https://coveralls.io/repos/github/ZhengHe-MD/gorefactor/badge.svg?branch=master)](https://coveralls.io/github/ZhengHe-MD/gorefactor?branch=master)
[![golang](https://img.shields.io/badge/Language-Go-green.svg?style=flat)](https://golang.org)
[![godoc](https://godoc.org/github.com/ZhengHe-MD/gorefactor?status.svg)](https://godoc.org/github.com/ZhengHe-MD/gorefactor)

## Installation

```sh
$ go get -u github.com/ZhengHe-MD/gorefactor
```

## Examples

* [insert context](/examples/insert_context.go)

## API

### parse src

```
ParseSrcFile(filename string) (df *dst.File, err error)
ParseSrcFileFromBytes(src []byte) (df *dst.File, err error)
```

### write src

```
FprintFile(out io.Writer, df *dst.File) error
```

### function body utilities

```
HasStmtInsideFuncBody(df *dst.File, funcName string, stmt dst.Stmt) (ret bool)
DeleteStmtFromFuncBody(df *dst.File, funcName string, stmt dst.Stmt) (modified bool)
AddStmtToFuncBody(df *dst.File, funcName string, stmt dst.Stmt, pos int) (modified bool)
AddStmtToFuncBodyStart(df *dst.File, funcName string, stmt dst.Stmt) (modified bool)
AddStmtToFuncBodyEnd(df *dst.File, funcName string, stmt dst.Stmt) (modified bool)
AddNodeToFuncBodyBefore(df *dst.File, funcName string, stmt, refStmt dst.Stmt) (modified bool) 
AddNodeToFuncBodyAfter(df *dst.File, funcName string, stmt, refStmt dst.Stmt) (modified bool)
```

### function call utilities

```
HasArgInCallExpr(df *dst.File, funcName string, arg dst.Expr) (ret bool)
DeleteArgFromCallExpr(df *dst.File, funcName string, arg dst.Expr) (modified bool)
AddArgToCallExpr(df *dst.File, funcName string, arg dst.Expr, pos int) (modified bool)
```

### function declaration utilities

```
HasFieldInFuncDeclParams(df *dst.File, funcName string, field *dst.Field) (ret bool)
DeleteFieldFromFuncDeclParams(df *dst.File, funcName string, field *dst.Field) (modified bool)
AddFieldToFuncDeclParams(df *dst.File, funcName string, field *dst.Field, pos int) (modified bool)
```
