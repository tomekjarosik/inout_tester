package main

import testcase "github.com/tomekjarosik/inout_tester/internal/testcase"

func CountMatching(collection []testcase.CompletedTestCase, pred func(testcase.CompletedTestCase) bool) int {
	res := 0
	for _, x := range collection {
		if pred(x) {
			res++
		}
	}
	return res
}
