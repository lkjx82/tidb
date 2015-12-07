// Copyright 2015 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package plan

import (
	"fmt"
	"strings"

	"github.com/ngaut/log"
)

// Explain explains a Plan, returns description string.
func Explain(p Plan) string {
	var e explainer
	p.Accept(&e)
	return strings.Join(e.strs, "->")
}

type explainer struct {
	strs []string
}

func (e *explainer) Enter(in Plan) (Plan, bool) {
	return in, false
}

func (e *explainer) Leave(in Plan) (Plan, bool) {
	var str string
	switch x := in.(type) {
	case *TableScan:
		str = fmt.Sprintf("Table(%s)", x.Table.Name.L)
	case *IndexScan:
		str = fmt.Sprintf("Index(%s.%s)", x.Table.Name.L, x.Index.Name.L)
	case *Filter:
		str = "Filter"
	case *SelectFields:
		str = "Fields"
	case *Sort:
		if x.Bypass {
			return in, true
		}
		str = "Sort"
	case *SelectLock:
		str = "Lock"
	case *Limit:
		str = "Limit"
	default:
		log.Fatalf("Unknown plan type %T", in)
	}
	e.strs = append(e.strs, str)
	return in, true
}
