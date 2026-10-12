# O(1) LFU Cache in Go

A Go implementation of the paper
[An O(1) algorithm for implementing the LFU cache eviction scheme](https://arxiv.org/abs/2110.11602).

## Overview

A Least Frequently Used (LFU) cache evicts the item with the lowest access count
when it runs out of space. Naive implementations (a min-heap or sorted list) cost
O(log n) or O(n) per access, because an item must be repositioned every time its
count changes. This implementation performs insert, access, and eviction in O(1).

## How it works

The design combines three pieces:

- **Frequency list:** a doubly linked list of frequency nodes, kept in ascending
  order. Each node represents one distinct access count.
- **Item lists:** each frequency node owns its own doubly linked list of the items
  currently at that frequency, ordered by recency so ties are broken by evicting
  the oldest item.
- **Hash map:** maps each key directly to its item node, so lookups never traverse
  a list. Each item also holds a pointer to its parent frequency node.

Because the lists are only used to maintain order, never to search, every operation
reduces to a fixed number of pointer updates:

| Operation | How it stays O(1) |
|-----------|-------------------|
| Get       | Map lookup, then move the item to the adjacent frequency node (count + 1) |
| Put       | Map insert; new items start at frequency 1 at the head of the frequency list |
| Evict     | Remove the oldest item from the lowest-frequency node, at the head of the list |
