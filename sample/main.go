package main

import (
	"bufio"
	"fmt"
	"github.com/justinsacbibit/gocalc"
	"math"
	"os"
)

func paramResolver(param string) interface{} {
	switch param {
	case "pi":
		return math.Pi
	default:
		return nil
	}
}

func pow(x int64, y int64) int64 {
	if y == 0 {
		return 1
	} else if y == 1 {
		return x
	} else if y%2 == 0 {
		p := pow(x, y/2)
		return p * p
	} else {
		return x * pow(x, y-1)
	}
}

func funcResolver(function string, args ...func() interface{}) (interface{}, bool) {
	switch function {
	case "pow":
		if len(args) != 2 {
			panic(gocalc.EvaluationError("pow(x, y) requires two arguments"))
		}

		var x, y interface{}
		x = args[0]()
		y = args[1]()
		return pow(x.(int64), y.(int64)), true
	}

	return nil, false
}

func main() {
	for {
		fmt.Print("input: ")
		b := bufio.NewReader(os.Stdin)
		input, _, _ := b.ReadLine()

		expr, err := gocalc.NewExpr(string(input))
		if err != nil {
			fmt.Println(err)
			continue
		}

		result, err := expr.Evaluate(paramResolver, funcResolver)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("result: ", result)
	}
}
