# GoRefactor

this module provides some basic utilities for you to do code migration on large code base written in golang.

[![Go Report Card](https://goreportcard.com/badge/github.com/ZhengHe-MD/gorefactor)](https://goreportcard.com/report/github.com/ZhengHe-MD/gorefactor)
[![Coverage Status](https://coveralls.io/repos/github/ZhengHe-MD/gorefactor/badge.svg?branch=master)](https://coveralls.io/github/ZhengHe-MD/gorefactor?branch=master)
[![golang](https://img.shields.io/badge/Language-Go-green.svg?style=flat)](https://golang.org)

## Installation

```sh
$ go get -u github.com/ZhengHe-MD/gorefactor
```

## Examples

for full example, check examples directory

### insert context

say we want to refactor the following code

```go
package main

func f() {}

func main() {
    f()
    f()
}
```

to 

```go
package main

import "context"

func f(ctx context.Context) {}

func main() {
    f(context.TODO())
    f(context.TODO())
}
```

[examples/insert_context.go](/examples/insert_context.go)

## API

### parse src

```go
func ParseSrcFile(filename string) (df *dst.File, err error)
func ParseSrcFileFromBytes(src []byte) (df *dst.File, err error)
```

### write src

```go
func FprintFile(out io.Writer, df *dst.File) error
```

### function body utilities

```go
func HasStmtInsideFuncBody(df *dst.File, funcName string, stmt dst.Stmt) (ret bool)
func DeleteStmtFromFuncBody(df *dst.File, funcName string, stmt dst.Stmt) (modified bool)
func AddStmtToFuncBody(df *dst.File, funcName string, stmt dst.Stmt, pos int) (modified bool)
func AddStmtToFuncBodyStart(df *dst.File, funcName string, stmt dst.Stmt) (modified bool)
func AddStmtToFuncBodyEnd(df *dst.File, funcName string, stmt dst.Stmt) (modified bool)
func AddNodeToFuncBodyBefore(df *dst.File, funcName string, stmt, refStmt dst.Stmt) (modified bool) 
func AddNodeToFuncBodyAfter(df *dst.File, funcName string, stmt, refStmt dst.Stmt) (modified bool)
```

### function call utilities

```go
func HasArgInCallExpr(df *dst.File, funcName string, arg dst.Expr) (ret bool)
func DeleteArgFromCallExpr(df *dst.File, funcName string, arg dst.Expr) (modified bool)
func AddArgToCallExpr(df *dst.File, funcName string, arg dst.Expr, pos int) (modified bool)
```

### function declaration utilities

```go
func HasFieldInFuncDeclParams(df *dst.File, funcName string, field *dst.Field) (ret bool)
func DeleteFieldFromFuncDeclParams(df *dst.File, funcName string, field *dst.Field) (modified bool)
func AddFieldToFuncDeclParams(df *dst.File, funcName string, field *dst.Field, pos int) (modified bool)
```
