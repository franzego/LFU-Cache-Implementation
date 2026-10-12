package main

import "fmt"

const capacity = 1 * 1024

type Cache struct {
	Items    map[string]*ItemNode
	FreqList FrequencyList
	Capacity int
}

func NewCache(list *FrequencyList) *Cache {
	return &Cache{
		Items:    make(map[string]*ItemNode),
		FreqList: *list,
		Capacity: capacity,
	}
}

/*
 * Three Different operations will be performed on the cached data:
 * - Set: Inserts an item to the cache.
 * - Retireve: Looksup an item from the cache, simulataneouly increasing the usage count for the LFU.
 * - Evict: Delete the LFU item from the cache on dearth of memory.
 */

func appendToNodeList(nl *NodeList, node *ItemNode) {
	node.Prev = nil
	node.Next = nil
	if nl.Tail == nil {
		nl.Head = node
		nl.Tail = node
		return
	}
	node.Prev = nl.Tail
	nl.Tail.Next = node
	nl.Tail = node
}

func removeFromNodeList(nl *NodeList, node *ItemNode) {
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		nl.Head = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		nl.Tail = node.Prev
	}
	node.Prev = nil
	node.Next = nil
}

func findFreqNode(head *FreqNode, freq int) *FreqNode {
	for cur := head; cur != nil; cur = cur.Next {
		if cur.Frequency == freq {
			return cur
		}
	}
	return nil
}

func insertFreqNodeSorted(list *FrequencyList, newNode *FreqNode) {
	if list.Head == nil {
		list.Head = newNode
		return
	}
	var prev *FreqNode
	cur := list.Head
	for cur != nil && cur.Frequency < newNode.Frequency {
		prev = cur
		cur = cur.Next
	}
	newNode.Prev = prev
	newNode.Next = cur
	if prev == nil {
		list.Head = newNode
	} else {
		prev.Next = newNode
	}
	if cur != nil {
		cur.Prev = newNode
	}
}

func (c *Cache) evictLFU() error {
	head := c.FreqList.Head
	if head == nil || head.valuelist == nil || head.valuelist.Head == nil {
		return fmt.Errorf("nothing to evict")
	}
	victim := head.valuelist.Head
	removeFromNodeList(head.valuelist, victim)
	delete(c.Items, victim.Cachedkey)
	// check if it that was the last element in the FreqNode
	if head.valuelist.Head == nil {
		c.FreqList.Head = head.Next
		if head.Next != nil {
			head.Next.Prev = nil
		}
	}
	return nil
}

func (c *Cache) InsertToCache(key, val string) (string, error) {
	// check if the cache has reached its limit. If it has, we remove the items at the Head of the FreqList and if there is only one element in it, on removing the item, we also remove the FreqNode at the same time.
	fmt.Println("Starting insertion")
	if _, exists := c.Items[key]; exists {
		return "key already exists", nil
	}
	if len(c.Items) >= c.Capacity {
		if err := c.evictLFU(); err != nil {
			return "", err
		}
		return c.insertFresh(key)
	}

	// New items must be placed in frequency 1.
	if c.FreqList.Head == nil {
		c.FreqList.Head = &FreqNode{
			Frequency: 1,
			valuelist: &NodeList{},
		}
	}

	if c.FreqList.Head.Frequency != 1 {
		freqNode := &FreqNode{
			Frequency: 1,
			Next:      c.FreqList.Head,
			valuelist: &NodeList{},
		}
		c.FreqList.Head.Prev = freqNode
		c.FreqList.Head = freqNode
		return c.insertFresh(key)
	}

	if c.FreqList.Head.valuelist == nil {
		c.FreqList.Head.valuelist = &NodeList{}
	}
	return c.insertFresh(key)
}

func (c *Cache) insertFresh(key string) (string, error) {
	head := c.FreqList.Head
	if head == nil || head.Frequency != 1 {
		return "", fmt.Errorf("freq 1 head missing")
	}
	node := &ItemNode{
		Cachedkey: key,
		Parent:    head,
	}
	appendToNodeList(head.valuelist, node)
	c.Items[key] = node
	return "Inserted to existing list's Head", nil
}

func (c *Cache) AccessAndIncrement(key string) (string, *ItemNode) {
	fmt.Println("starting access")
	// look up the key in the cache
	node, ok := c.Items[key]
	if !ok || node == nil {
		return "failed to find node", nil // "cache miss"
	}
	fmt.Printf("gotten from node: %s", node.Cachedkey)

	parent := node.Parent // we locate the exact position in the Frequencylist
	if parent == nil {
		return "failed to find node", nil
	}
	newFreq := parent.Frequency + 1

	// unlink node from its old frequency list
	removeFromNodeList(parent.valuelist, node)

	var target *FreqNode
	if parent.valuelist.Head == nil {
		// old frequency node is now empty -> remove it from frequency list
		prev := parent.Prev
		next := parent.Next
		if prev != nil {
			prev.Next = next
		} else {
			c.FreqList.Head = next
		}
		if next != nil {
			next.Prev = prev
		}
		// find or create the new frequency node
		target = findFreqNode(c.FreqList.Head, newFreq)
		if target == nil {
			target = &FreqNode{
				Frequency: newFreq,
				valuelist: &NodeList{},
			}
			insertFreqNodeSorted(&c.FreqList, target)
		}
	} else {
		// parent still has items, look for next freq node
		if parent.Next != nil && parent.Next.Frequency == newFreq {
			target = parent.Next
		} else {
			target = &FreqNode{
				Frequency: newFreq,
				Prev:      parent,
				Next:      parent.Next,
				valuelist: &NodeList{},
			}
			if parent.Next != nil {
				parent.Next.Prev = target
			}
			parent.Next = target
		}
	}

	if target.valuelist == nil {
		target.valuelist = &NodeList{}
	}
	node.Parent = target
	appendToNodeList(target.valuelist, node)

	return "successful", node
}
