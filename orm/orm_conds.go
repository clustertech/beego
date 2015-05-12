// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orm

import (
	"fmt"
	"strings"
	"time"
)

const (
	ExprSep = "__"
)

type Condition interface {
	toSQL(t *dbTables, tz *time.Location) (string, []interface{})
}

type simpleCond struct {
	exprs []string
	args  []interface{}
}

type junctionCond struct {
	conds []Condition
	op    string
}

type notCond struct {
	cond Condition
}

func (cond *junctionCond) toSQL(t *dbTables, tz *time.Location) (where string, params []interface{}) {
	for i, c := range cond.conds {
		w, ps := c.toSQL(t, tz)
		w = fmt.Sprintf("( %s) ", w)
		if i > 0 {
			where += cond.op + " "
		}
		where += w
		params = append(params, ps...)
	}
	return
}

func (cond *notCond) toSQL(t *dbTables, tz *time.Location) (where string, params []interface{}) {
	w, ps := cond.cond.toSQL(t, tz)
	where += fmt.Sprintf("NOT ( %s) ", w)
	params = append(params, ps...)
	return
}

func (cond *simpleCond) toSQL(t *dbTables, tz *time.Location) (where string, params []interface{}) {
	exprs := cond.exprs

	Q := t.base.TableQuote()
	mi := t.mi

	num := len(exprs) - 1
	operator := ""
	if operators[exprs[num]] {
		operator = exprs[num]
		exprs = exprs[:num]
	}

	index, _, fi, suc := t.parseExprs(mi, exprs)
	if suc == false {
		panic(fmt.Errorf("unknown field/column name `%s`", strings.Join(cond.exprs, ExprSep)))
	}

	if operator == "" {
		operator = "exact"
	}

	operSql, args := t.base.GenerateOperatorSql(mi, fi, operator, cond.args, tz)

	leftCol := fmt.Sprintf("%s.%s%s%s", index, Q, fi.column, Q)
	t.base.GenerateOperatorLeftCol(fi, operator, &leftCol)

	where += fmt.Sprintf("%s %s ", leftCol, operSql)
	params = append(params, args...)

	return
}

func And(conds ...Condition) Condition {
	if len(conds) == 0 {
		panic(fmt.Errorf("<Condition.And> args cannot empty"))
	}
	return &junctionCond{conds, "AND"}
}

func Or(conds ...Condition) Condition {
	if len(conds) == 0 {
		panic(fmt.Errorf("<Condition.Or> args cannot empty"))
	}
	return &junctionCond{conds, "OR"}
}

func Not(cond Condition) Condition {
	if cond == nil {
		panic(fmt.Errorf("<Condition.Not> arg cannot empty"))
	}
	return &notCond{cond}
}

func Cond(expr string, args ...interface{}) Condition {
	if expr == "" || len(args) == 0 {
		panic(fmt.Errorf("<Condition> args cannot empty"))
	}
	return &simpleCond{exprs: strings.Split(expr, ExprSep), args: args}
}
