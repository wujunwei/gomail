package sortlist

import (
	"fmt"
	"sort"
)

type Compare[T comparable] func(a, b T) int

var IntCompare Compare[int] = func(a, b int) int {
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}

// StringCompare string compare default
var StringCompare Compare[string] = func(a, b string) int {
	if a > b {
		return 1
	}
	if a < b {
		return -1
	}
	return 0
}

type ForEach[T comparable] func(index int, a T)

// PrintEach for debug usage
func PrintEach[T comparable](index int, a T) {
	fmt.Printf("index: %d, value: %v \n", index, a)
}

// BisectRight bisect search right most element
func BisectRight[T comparable](l []T, c Compare[T], target T) int {
	return sort.Search(len(l), func(i int) bool {
		return c(l[i], target) > 0
	})
}

// BisectLeft bisect search left most element
func BisectLeft[T comparable](l []T, c Compare[T], target T) int {
	return sort.Search(len(l), func(i int) bool {
		return c(l[i], target) >= 0
	})
}

// InSort insert element at the suitable location in a sorted slice
func InSort[T comparable](l []T, c Compare[T], a T) []T {
	var zeroValue T
	index := BisectRight(l, c, a)
	if index == len(l) {
		return append(l, a)
	}
	l = append(l, zeroValue)
	copy(l[index+1:], l[index:len(l)-1])
	l[index] = a
	return l
}

// RemoveSort remove element at the suitable location in a sorted slice
func RemoveSort[T comparable](l []T, c Compare[T], a T) ([]T, bool) {
	index := BisectRight(l, c, a)
	if index == 0 {
		return l, false
	}
	if l[index-1] == a {
		return Remove[T](l, index-1), true
	}
	return l, false
}

func Remove[T comparable](l []T, index int) []T {
	copy(l[index:], l[index+1:])
	return l[:len(l)-1]
}
