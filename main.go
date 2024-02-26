package main

import (
	"math/rand"
)

func main() {

	qs := []Q{
		Q{},
	}

	QS := BuildQStore(qs)

	QS.PickQSet(map[int]int{
		1: 10,
		2: 2,
		3: 2,
	})

}

type QStore struct {
	ByTypes map[int][]*Q
}

func BuildQStore(Qs []Q) (ret QStore) {

	for _, q := range Qs {
		ret.ByTypes[q.Type] = append(ret.ByTypes[q.Type], &q)
	}

	return
}

// number of qs per type in opt
func (QS *QStore) PickQSet(opt map[int]int) (ret map[int][]*Q) {

	for qType, qs := range QS.ByTypes {
		ret[qType] = append(ret[qType], pickQ(qs, opt[qType])...)
	}

	return
}

func pickQ(qs []*Q, nums int) (ret []*Q) {

	r := rand.Intn(len(qs) - 1)

	for i := r + 1; i < len(qs); i = (i + 1) % len(qs) {
		if nums == 0 {
			return
		}
		if i == r {
			q := qs[rand.Intn(len(qs))]
			ret = append(ret, q)
			q.PickedCount++
			nums--

			continue
		}

		if qs[i].PickedCount == 0 {
			ret = append(ret, qs[i])
			qs[i].PickedCount++
			nums--
		}

	}
	return
}
