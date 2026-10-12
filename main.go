package main

import "fmt"

func main() {
	freqNode := NewFrequnecyNode()
	list := NewFrequencyList(freqNode)
	cache := NewCache(list)
	k1 := "David"
	v1 := "Engineer"
	k2 := "musa"
	v2 := "Doctor"

	val1, err := cache.InsertToCache(k1, v1)
	if err != nil {
		fmt.Println(err)
	}
	valu2, err := cache.InsertToCache(k2, v2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(val1)
	fmt.Println(valu2)

	answer, node := cache.AccessAndIncrement(k1)
	fmt.Print(answer)
	fmt.Println(node.Cachedkey)
}
