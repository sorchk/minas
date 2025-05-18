// Package el 提供表达式语言支持
package el

import (
	"fmt"

	"github.com/expr-lang/expr"
)

func Evaluate(el string, data map[string]interface{}) (any, error) {
	fmt.Printf("Evaluating expression: %s\n", el)
	fmt.Printf("Using environment: %v\n", data)
	program, err := expr.Compile(el, expr.Env(data))
	if err != nil {
		panic(err)
	}
	output, err := expr.Run(program, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("output: %v type: %T\n", output, output)
	return output, nil
}
