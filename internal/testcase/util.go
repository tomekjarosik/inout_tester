package testcase

func CountMatching(collection []CompletedTestCase, pred func(CompletedTestCase) bool) int {
	res := 0
	for _, x := range collection {
		if pred(x) {
			res++
		}
	}
	return res
}
