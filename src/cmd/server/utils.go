package main

import "github.com/wizenheimer/bloombox/pkg/emailchecker"

// countValid counts the number of valid emails in the results
func countValid(results []*emailchecker.CheckResult) int {
	count := 0
	for _, result := range results {
		if result.IsValid {
			count++
		}
	}
	return count
}

// countDisposable counts the number of disposable emails in the results
func countDisposable(results []*emailchecker.CheckResult) int {
	count := 0
	for _, result := range results {
		if result.Summary.IsDisposable {
			count++
		}
	}
	return count
}

// countFree counts the number of free emails in the results
func countFree(results []*emailchecker.CheckResult) int {
	count := 0
	for _, result := range results {
		if result.Summary.IsFree {
			count++
		}
	}
	return count
}

// countRole counts the number of role emails in the results
func countRole(results []*emailchecker.CheckResult) int {
	count := 0
	for _, result := range results {
		if result.Summary.IsRole {
			count++
		}
	}
	return count
}
