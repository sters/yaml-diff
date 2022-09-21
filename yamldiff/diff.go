package yamldiff

import "fmt"

type DiffStatus int

const (
	DiffStatusUnknown  DiffStatus = 0
	DiffStatusSame     DiffStatus = 1
	DiffStatusDiff     DiffStatus = 2
	DiffStatus1Missing DiffStatus = 3
	DiffStatus2Missing DiffStatus = 4
)

type (
	rawType      = interface{}
	rawTypeMap   = map[string]rawType
	rawTypeArray = []rawType

	diff struct {
		a        rawType
		b        rawType
		children *diffChildren

		status    DiffStatus
		diffCount int
		treeLevel int
	}

	diffChildrenArray = []*diff
	diffChildrenMap   = map[string]*diff

	diffChildren struct {
		a diffChildrenArray
		m diffChildrenMap
	}
)

func performDiff(rawA rawType, rawB rawType, level int) *diff {
	if rawA == nil || rawB == nil {
		return handlePrimitive(rawA, rawB, level)
	}

	if r := handleMap(rawA, rawB, level); r != nil {
		return r
	}

	if r := handleArray(rawA, rawB, level); r != nil {
		return r
	}

	// other case -> handle as primitive (int/float/bool/string)
	return handlePrimitive(rawA, rawB, level)
}

// TODO: golang's map is not stable
func handleMap(rawA rawType, rawB rawType, level int) *diff {
	result := &diff{
		a:         rawA,
		b:         rawB,
		treeLevel: level,
	}

	mapA, mapAok := tryMap(rawA)
	mapB, mapBok := tryMap(rawB)

	// if both are not map
	if !mapAok && !mapBok {
		return nil
	}

	// if A is not map but B is map -> it's different data
	if !mapAok || !mapBok {
		result.status = DiffStatusDiff
		result.diffCount = handlePrimitive(rawA, rawB, level).diffCount

		return result
	}

	// if both are map

	result.children = &diffChildren{
		m: diffChildrenMap{},
	}
	result.status = DiffStatusSame

	// if B is map -> check the same key children
	for keyA, valA := range mapA {
		foundKey := false
		for keyB, valB := range mapB {
			if keyA != keyB {
				continue
			}

			result.children.m[keyA] = performDiff(valA, valB, level+1)
			if result.children.m[keyA].status != DiffStatusSame {
				result.status = DiffStatusDiff // top level diff can't specify actual reason
			}

			foundKey = true

			break
		}

		if !foundKey {
			result.children.m[keyA] = performDiff(valA, nil, level+1)
			result.status = DiffStatusDiff // top level diff can't specify actual reason
		}
	}

	// finding missing keyA
	for keyB, valB := range mapB {
		foundKey := false
		for keyA := range mapA {
			if keyB != keyA {
				continue
			}

			foundKey = true

			break
		}

		if !foundKey {
			result.children.m[keyB] = performDiff(nil, valB, level+1)
			result.status = DiffStatusDiff // top level diff can't specify actual reason
		}
	}

	sum := 0
	for _, v := range result.children.m {
		sum += v.diffCount
	}
	result.diffCount = sum

	return result
}

func handleArray(rawA rawType, rawB rawType, level int) *diff {
	result := &diff{
		a:         rawA,
		b:         rawB,
		treeLevel: level,
	}

	arrayA, arrayAok := tryArray(rawA)
	arrayB, arrayBok := tryArray(rawB)

	// if both are not array
	if !arrayAok && !arrayBok {
		return nil
	}

	// if A is not array but B is array -> it's different data
	if !arrayAok || !arrayBok {
		result.status = DiffStatusDiff
		result.diffCount = handlePrimitive(rawA, rawB, level).diffCount

		return result
	}

	// if both are array

	result.children = &diffChildren{
		a: diffChildrenArray{},
	}
	result.status = DiffStatusSame

	// check each elements is same or not
	diffs := map[string]*diff{}
	foundA := map[int]struct{}{}
	foundB := map[int]struct{}{}

	for keyA, valA := range arrayA {
		for keyB, valB := range arrayB {
			key := fmt.Sprintf("%d-%d", keyA, keyB)

			diffs[key] = performDiff(valA, valB, level+1)
			if diffs[key].status == DiffStatusSame {
				// store result and mark as confirmed
				result.children.a = append(result.children.a, diffs[key])
				foundA[keyA] = struct{}{}
				foundB[keyB] = struct{}{}

				break
			}
		}
	}

	// found all elements, it's same array
	if len(foundA) == len(arrayA) && len(foundB) == len(arrayB) {
		return result
	}

	result.status = DiffStatusDiff

	// check diff elements
	for {
		// arrayA < arrayB, and all confirmed arrayA
		if len(foundA) == len(arrayA) {
			for k, v := range arrayB {
				if _, ok := foundB[k]; ok {
					continue
				}

				result.children.a = append(result.children.a, performDiff(nil, v, level+1))
			}

			break
		}

		// arrayB < arrayA, and all confirmed arrayB
		if len(foundB) == len(arrayB) {
			for k, v := range arrayA {
				if _, ok := foundA[k]; ok {
					continue
				}

				result.children.a = append(result.children.a, performDiff(v, nil, level+1))
			}

			break
		}

		smallestDiff := &diff{diffCount: 100000} // FIXME
		smallestKeyA := 0
		smallestKeyB := 0

		for keyA := range arrayA {
			if _, ok := foundA[keyA]; ok {
				continue
			}

			for keyB := range arrayB {
				if _, ok := foundB[keyB]; ok {
					continue
				}

				key := fmt.Sprintf("%d-%d", keyA, keyB)
				if diffs[key].status == DiffStatusSame {
					continue
				}

				if smallestDiff.diffCount > diffs[key].diffCount {
					smallestDiff = diffs[key]
					smallestKeyA = keyA
					smallestKeyB = keyB
				}
			}
		}

		result.children.a = append(result.children.a, smallestDiff)
		foundA[smallestKeyA] = struct{}{}
		foundB[smallestKeyB] = struct{}{}
	}

	sum := 0
	for _, v := range result.children.a {
		sum += v.diffCount
	}
	result.diffCount = sum

	return result
}

func handlePrimitive(rawA rawType, rawB rawType, level int) *diff {
	result := &diff{
		a:         rawA,
		b:         rawB,
		treeLevel: level,
	}

	strA := []rune(fmt.Sprint(rawA))
	strB := []rune(fmt.Sprint(rawB))

	switch {
	case rawA == rawB:
		result.status = DiffStatusSame
	case rawA == nil:
		result.status = DiffStatus1Missing
		result.diffCount = len(strB)
	case rawB == nil:
		result.status = DiffStatus2Missing
		result.diffCount = len(strA)
	default:
		result.status = DiffStatusDiff
	}

	// calculate diff size for diff
	if result.status == DiffStatusDiff {
		maxLen := len(strA)
		if lenB := len(strB); maxLen < lenB {
			maxLen = lenB
		}

		for nA, a := range strA {
			// lenA > lenB
			if len(strB) <= nA {
				result.diffCount = maxLen - nA

				break
			}

			// found diff in A and B strings
			if b := strB[nA]; a != b {
				result.diffCount = maxLen - nA

				break
			}
		}

		// guess lenA < lemB
		if result.diffCount == 0 {
			result.diffCount = maxLen - len(strA)
		}
	}

	return result
}

func tryMap(x rawType) (rawTypeMap, bool) {
	m, ok := x.(map[string]interface{})

	return m, ok
}

func tryArray(x rawType) (rawTypeArray, bool) {
	a, ok := x.([]interface{})

	return a, ok
}
