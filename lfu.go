package main

// FreqNode represents a node in the frequencyList. It simply acts as a counter for the LFU eviction scheme. Each Node represents one distinct frequnecy value e.g 1, 4,5,7.
// All items inserted to the cache start at a frequency of 1. This node is usually created as the Head. The subsequent increase in count lead to creation of new frequency nodes.
// If an item has been accessed 4 times, it will be added to the freqNode with a frequnecy of 4, in the valuelist holding other nodes that have the same Frequency.
// The destruction of a FreqNode is tied to the items in the Nodelist. If there is nothing in the list on the eviction of the last node, then the freqNode is destroyed.
type FreqNode struct {
	Frequency int
	Prev      *FreqNode
	Next      *FreqNode
	valuelist *NodeList
}

type ItemNode struct {
	Cachedkey string
	Prev      *ItemNode
	Next      *ItemNode
	Parent    *FreqNode
}

// FrequencyList is a doubly linked list for the FreqNode
type FrequencyList struct {
	Head *FreqNode
}

///////////////////////////////////////

// NodeList is a doubly linked list for the ItemNode
type NodeList struct {
	Head *ItemNode
	Tail *ItemNode
}

func NewFrequencyList(freqNode *FreqNode) *FrequencyList {
	return &FrequencyList{
		Head: freqNode,
	}
}

func NewFrequnecyNode() *FreqNode {
	return &FreqNode{
		Frequency: 1,
		valuelist: &NodeList{},
	}
}
