package utils

func Intersection(a, b []string) (res int) {
	l, r := 0, 0
	for l < len(a) && r < len(b) {
		if a[l] == a[r] {
			res++
			l++
			r++
		} else if a[l] < a[r] {
			l++
		} else {
			r++
		}
	}
	return
}

func CalcMarkLowerBound(total int, need int) int {
	for c := 1; c <= total; c++ {
		cur := float64(c) / float64(total) * 10.
		if ConvertMark(cur) >= need {
			return c
		}
	}
	return -1
}

func ConvertMark(mark float64) int {
	if mark < 3.5 {
		return 2
	} else if mark < 5.5 {
		return 3
	} else if mark < 7.5 {
		return 4
	} else {
		return 5
	}
}
