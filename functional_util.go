package main

func CountMatching(collection []TestCase, pred func(TestCase) bool) int {
	res := 0
	for _, x := range collection {
		if pred(x) {
			res++
		}
	}
	return res
}
