package collection

import "sync"

func CountSyncMap(m *sync.Map) int {
	size := 0
	m.Range(func(key, value interface{}) bool {
		size++
		return true
	})
	return size
}

func Sum(m map[int32]int64) int64 {
	sum := int64(0)
	for _, v := range m {
		sum += v
	}
	return sum
}

func FindMax(m map[int32]int64) (int32, int64) {
	var (
		index int32
		max   int64
	)
	isInit := true
	for i, v := range m {
		if isInit {
			index = i
			max = v
			isInit = false
			continue
		}
		if v > max {
			max = v
			index = i
		}
	}
	return index, max
}

func FindMinNotZero(m map[int32]int64) (int32, int64, bool, int, []int32) {
	var (
		index        int32
		min          int64
		count        int
		notZeroCount int
		isSame       bool
		chairs       []int32
	)
	isInit := true
	for i, v := range m {
		if v <= 0 {
			continue
		}
		notZeroCount++
		chairs = append(chairs, i)
		if isInit {
			index = i
			min = v
			isInit = false
			continue
		}
		if v < min {
			min = v
			index = i
		} else if v == min {
			count++
		}
	}
	if count == notZeroCount-1 {
		isSame = true
	}
	return index, min, isSame, notZeroCount, chairs
}
