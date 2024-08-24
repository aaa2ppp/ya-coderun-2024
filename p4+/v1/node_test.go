package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func makeTree(keys []int) *node {
	var root *node
	for _, k := range keys {
		root = root.Insert(k)
	}
	return root
}

type node = Node[int]

func lnr(root *node) []int {
	keys := make([]int, 0, root.Size())

	var dfs func(*node)
	dfs = func(node *node) {
		if node == nil {
			return
		}
		dfs(node.left)
		keys = append(keys, node.key)
		dfs(node.right)
	}

	dfs(root)
	return keys
}

func nlr(root *node) []int {
	keys := make([]int, 0, root.Size())

	var dfs func(*node)
	dfs = func(node *node) {
		if node == nil {
			return
		}
		keys = append(keys, node.key)
		dfs(node.left)
		dfs(node.right)
	}

	dfs(root)
	return keys
}

func TestNode_Insert(t *testing.T) {
	tests := []struct {
		keys       []int
		wantLen    int
		wantHeight int
		wantLNR    []int
		wantNLR    []int
	}{
		{
			[]int{},
			0, 0,
			[]int{},
			[]int{},
		},
		{
			[]int{1},
			1, 1,
			[]int{1},
			[]int{1},
		},
		{
			[]int{2, 1, 3},
			3, 2,
			[]int{1, 2, 3},
			[]int{2, 1, 3},
		},
		{
			[]int{2, 3, 1},
			3, 2,
			[]int{1, 2, 3},
			[]int{2, 1, 3},
		},
		{
			[]int{1, 2, 3},
			3, 2,
			[]int{1, 2, 3},
			[]int{2, 1, 3},
		},
		{
			[]int{1, 1, 2},
			3, 2,
			[]int{1, 1, 2},
			[]int{1, 1, 2},
		},
		{
			[]int{2, 1, 1},
			3, 2,
			[]int{1, 1, 2},
			[]int{1, 1, 2},
		},
		{
			[]int{4, 6, 2, 5, 1, 7, 3},
			7, 3,
			[]int{1, 2, 3, 4, 5, 6, 7},
			[]int{4, 2, 1, 3, 6, 5, 7},
		},
		{
			[]int{4, 2, 6, 1, 3, 5, 7},
			7, 3,
			[]int{1, 2, 3, 4, 5, 6, 7},
			[]int{4, 2, 1, 3, 6, 5, 7},
		},
		{
			[]int{1, 2, 3, 4, 5, 6, 7},
			7, 3,
			[]int{1, 2, 3, 4, 5, 6, 7},
			[]int{4, 2, 1, 3, 6, 5, 7},
		},
		{
			[]int{7, 6, 5, 4, 3, 2, 1},
			7, 3,
			[]int{1, 2, 3, 4, 5, 6, 7},
			[]int{4, 2, 1, 3, 6, 5, 7},
		},
		{
			[]int{1, 2, 2, 4, 4, 6, 7},
			7, 3,
			[]int{1, 2, 2, 4, 4, 6, 7},
			[]int{4, 2, 1, 2, 6, 4, 7},
		},
		{
			[]int{7, 6, 4, 4, 2, 2, 1},
			7, 3,
			[]int{1, 2, 2, 4, 4, 6, 7},
			[]int{4, 2, 1, 2, 6, 4, 7},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.keys), func(t *testing.T) {
			root := makeTree(tt.keys)
			if got := root.Size(); got != tt.wantLen {
				t.Errorf("Len() = %v, want %v", got, tt.wantLen)
			}
			if got := root.height(); got != tt.wantHeight {
				t.Errorf("Height() = %v, want %v", got, tt.wantHeight)
			}
			if got := lnr(root); !reflect.DeepEqual(got, tt.wantLNR) {
				t.Errorf("lnr() = %v, want %v", got, tt.wantLNR)
			}
			if got := nlr(root); !reflect.DeepEqual(got, tt.wantNLR) {
				t.Errorf("nlr() = %v, want %v", got, tt.wantNLR)
			}
		})
	}
}

// func TestNode_Remove(t *testing.T) {

// 	tests := []struct {
// 		root       *node
// 		key        int
// 		wantOk     bool
// 		wantLen    int
// 		wantHeight int
// 		wantLNR    []int
// 		wantNLR    []int
// 	}{
// 		{
// 			makeTree([]int{2, 1, 3}),
// 			2, true,
// 			2, 2,
// 			[]int{1, 3},
// 			[]int{3, 1},
// 		},
// 		{
// 			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
// 			5, true,
// 			6, 3,
// 			[]int{1, 2, 3, 4, 6, 7},
// 			[]int{4, 2, 1, 3, 6, 7},
// 		},
// 		{
// 			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
// 			8, false,
// 			7, 3,
// 			[]int{1, 2, 3, 4, 5, 6, 7},
// 			[]int{4, 2, 1, 3, 6, 5, 7},
// 		},
// 		{
// 			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
// 			0, false,
// 			7, 3,
// 			[]int{1, 2, 3, 4, 5, 6, 7},
// 			[]int{4, 2, 1, 3, 6, 5, 7},
// 		},
// 		{
// 			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
// 			7, true,
// 			6, 3,
// 			[]int{1, 2, 3, 4, 5, 6},
// 			[]int{4, 2, 1, 3, 6, 5},
// 		},
// 		{
// 			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
// 			1, true,
// 			6, 3,
// 			[]int{2, 3, 4, 5, 6, 7},
// 			[]int{4, 2, 3, 6, 5, 7},
// 		},
// 		{
// 			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
// 			3, true,
// 			6, 3,
// 			[]int{1, 2, 4, 5, 6, 7},
// 			[]int{4, 2, 1, 6, 5, 7},
// 		},
// 		{
// 			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
// 			2, true,
// 			6, 3,
// 			[]int{1, 3, 4, 5, 6, 7},
// 			[]int{4, 3, 1, 6, 5, 7},
// 		},
// 		{
// 			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
// 			6, true,
// 			6, 3,
// 			[]int{1, 2, 3, 4, 5, 7},
// 			[]int{4, 2, 1, 3, 7, 5},
// 		},
// 		{
// 			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
// 			4, true,
// 			6, 3,
// 			[]int{1, 2, 3, 5, 6, 7},
// 			[]int{5, 2, 1, 3, 6, 7},
// 		},
// 		{
// 			makeTree([]int{4, 2, 1, 3, 7, 5, 6, 8}),
// 			4, true,
// 			7, 3,
// 			[]int{1, 2, 3, 5, 6, 7, 8},
// 			[]int{5, 2, 1, 3, 7, 6, 8},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(fmt.Sprintf("%v remove %d", lnr(tt.root), tt.key), func(t *testing.T) {
// 			new_root, old_node := tt.root.Remove(tt.key)
// 			tt.root = new_root

// 			if got := old_node != nil; got != tt.wantOk {
// 				t.Errorf("gotOk = %v, want %v", got, tt.wantOk)
// 			}
// 			if got := tt.root.Size(); got != tt.wantLen {
// 				t.Errorf("Len() = %v, want %v", got, tt.wantLen)
// 			}
// 			if got := tt.root.height(); got != tt.wantHeight {
// 				t.Errorf("Height() = %v, want %v", got, tt.wantHeight)
// 			}
// 			if got := lnr(tt.root); !reflect.DeepEqual(got, tt.wantLNR) {
// 				t.Errorf("lnr = %v, want %v", got, tt.wantLNR)
// 			}
// 			if got := nlr(tt.root); !reflect.DeepEqual(got, tt.wantNLR) {
// 				t.Errorf("nlr() = %v, want %v", got, tt.wantNLR)
// 			}
// 		})
// 	}
// }

func TestNode_FindIdx(t *testing.T) {

	tests := []struct {
		root    *node
		key     int
		wantKey int
		wantIdx int
	}{
		{
			makeTree([]int{2, 1, 3}),
			2,
			2, 1,
		},
		{
			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
			5,
			5, 4,
		},
		{
			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
			8,
			0, 7,
		},
		{
			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
			0,
			1, 0,
		},
		{
			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
			7,
			7, 6,
		},
		{
			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
			1,
			1, 0,
		},
		{
			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
			3,
			3, 2,
		},
		{
			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
			2,
			2, 1,
		},
		{
			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
			6,
			6, 5,
		},
		{
			makeTree([]int{4, 2, 6, 1, 3, 5, 7}),
			4,
			4, 3,
		},
		{
			makeTree([]int{4, 2, 1, 3, 7, 5, 6, 8}),
			4,
			4, 3,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			0,
			1, 0,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			1,
			1, 0,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			2,
			3, 1,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			3,
			3, 1,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			4,
			5, 2,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			5,
			5, 2,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			6,
			7, 3,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			7,
			7, 3,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			8,
			9, 4,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			9,
			9, 4,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			10,
			11, 5,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			11,
			11, 5,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			12,
			13, 6,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			13,
			13, 6,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			14,
			15, 7,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			15,
			15, 7,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			16,
			0, 8,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			-100,
			1, 0,
		},
		{
			makeTree([]int{1, 3, 5, 7, 9, 11, 13, 15}),
			100,
			0, 8,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v get %d", lnr(tt.root), tt.key), func(t *testing.T) {
			gotIdx := tt.root.FindIdx(tt.key)

			if gotIdx != tt.wantIdx {
				t.Errorf("gotIdx = %v, want %v", gotIdx, tt.wantIdx)
			}
		})
	}
}

func FuzzNode_FindIdx(f *testing.F) {
	f.Add(int64(1))

	f.Fuzz(func(t *testing.T, seed int64) {
		rand := rand.New(rand.NewSource(seed))
		n := rand.Intn(1000000) + 1
		keys := make([]int, n)
		for i := range keys {
			keys[i] = rand.Intn(1000000) + 1
		}
		tree := makeTree(keys)
		sort.Ints(keys)
		for i := 0; i < n*2; i++ {
			key := rand.Intn(1100000) - 50000
			wantIdx := sort.Search(n, func(i int) bool {
				return keys[i] >= key
			})
			gotIdx := tree.FindIdx(key)
			if gotIdx != wantIdx {
				t.Errorf("seed=%d key=%d, got idx %d, want %d", seed, key, gotIdx, wantIdx)
				break
			}
		}
	})
}
