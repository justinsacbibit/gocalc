// Package gocalc is for evaluating mathematical expressions at runtime.
//
// See https://github.com/justinsacbibit/gocalc-cl-sample for a sample
// console app that uses gocalc.
//
//   func main() {
//     expression, err := gocalc.NewExpr("1 + 2")
//     if err != nil {
//       fmt.Println(err)
//       return
//     }
//
//     result, err := expression.Evaluate(nil, nil)
//     if err != nil {
//       fmt.Println(err)
//       return
//     }
//
//     fmt.Println(result)
//   }
package gocalc
