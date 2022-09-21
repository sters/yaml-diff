package yamldiff

type DiffStatus int

const (
	DiffStatusUnknown  DiffStatus = 0
	DiffStatusSame     DiffStatus = 1
	DiffStatusDiff     DiffStatus = 2
	DiffStatus1Missing DiffStatus = 3
	DiffStatus2Missing DiffStatus = 4
)

type YamlDiff struct {
	rawA   *RawYaml
	rawB   *RawYaml
	result interface{} // should be same structure with rawA but nested YamlDiff if tree
}

type rawRaw = interface{}
type rawMap = map[string]rawRaw
type rawSlice = []rawRaw

type rawDiff struct {
	a     rawRaw
	b     rawRaw
	child *diffChildren

	status DiffStatus
}

type diffChildren struct {
	a []*rawDiff
	m map[string]*rawDiff
}

// TODO: array support
func diff(rawA rawRaw, rawB rawRaw) *rawDiff {
	result := &rawDiff{
		a: rawA,
		b: rawB,
	}
	mapA, mapAok := tryMap(rawA)
	mapB, mapBok := tryMap(rawB)

	// if A is map
	if mapAok {
		// if B is not map -> it's different data
		if !mapBok {
			result.status = DiffStatusDiff

			return result
		}

		result.child = &diffChildren{
			m: map[string]*rawDiff{},
		}

		// if B is map -> check the same key children
		for keyA, valA := range mapA {
			foundKey := false
			for keyB, valB := range mapB {
				if keyA != keyB {
					continue
				}

				result.child.m[keyA] = diff(valA, valB)
				foundKey = true

				break
			}

			if !foundKey {
				result.child.m[keyA] = &rawDiff{
					a:      valA,
					status: DiffStatus2Missing,
				}
			}
		}

		// finding missing keyA
		for keyB, valB := range mapB {
			foundKey := false
			for keyA, _ := range mapA {
				if keyB != keyA {
					continue
				}

				foundKey = true

				break
			}

			if !foundKey {
				result.child.m[keyB] = &rawDiff{
					b:      valB,
					status: DiffStatus1Missing,
				}
			}
		}

		result.status = DiffStatusUnknown

		return result
	}

	// if A is not map but B is map -> it's different data
	if !mapAok && mapBok {
		result.status = DiffStatusDiff

		return result
	}

	// if A and B is not map -> int/float/string
	if !mapAok && !mapBok {
		switch {
		case rawA == rawB:
			result.status = DiffStatusSame
		case rawA == nil:
			result.status = DiffStatus1Missing
		case rawB == nil:
			result.status = DiffStatus2Missing
		default:
			result.status = DiffStatusDiff
		}

		return result
	}

	// unexpected case
	return result
}

func tryMap(x rawRaw) (rawMap, bool) {
	mapp, ok := x.(map[string]interface{})
	return mapp, ok
}
