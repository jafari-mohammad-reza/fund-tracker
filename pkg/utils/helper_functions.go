package utils

type Item struct {
	Base   string
	Target string
}

func FindIndex(items []Item, base string, target string) int {
	for i, v := range items {
		if v.Base == base && v.Target == target {
			return i
		}
	}
	return -1
}
