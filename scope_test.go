package container

import (
	"testing"
)

func TestInsert(t *testing.T) {
	s := NewScope()
	items := []ScopeItem{
		{100, 400}, {600, 200}, {50, 40}, {10, 40}, {550, 20}, {570, 30}, {1000, 300}, {1300, 200}, {900, 100},
	}
	itemsShould := []ScopeItem{
		{10, 80}, {100, 400}, {550, 250}, {900, 600},
	}

	size := uint64(0)
	for _, v := range items {
		s.Insert(v.Address, v.Size)
		size += v.Size
	}

	should(t, s, itemsShould, size)
}

func should(t *testing.T, s *Scope, items []ScopeItem, size uint64) {
	if s.size != size {
		t.Errorf("total size error, should: %v, get: %v\n", size, s.size)
	} else if s.items.Len() != len(items) {
		t.Errorf("items count error, should: %v, get: %v\n", len(items), s.items.Len())
	} else {
		e := s.items.Front()
		for i, v := range items {
			cur := e.Value.(ScopeItem)
			if cur.Address != v.Address || cur.Size != v.Size {
				t.Errorf("%v addr: %v size: %v, should addr: %v, size: %v\n", i, cur.Address, cur.Size, v.Address, v.Size)
			}
			e = e.Next()
		}
	}
}

func TestRemove(t *testing.T) {
	s := NewScope()
	items := []ScopeItem{
		{550, 250}, {900, 600}, {10, 80}, {100, 400},
	}

	size := uint64(0)
	for _, v := range items {
		s.Insert(v.Address, v.Size)
		size += v.Size
	}

	//删除元素的开始
	s.Remove(900, 200)
	s.Remove(10, 20)
	s.Remove(550, 50)
	s.Remove(100, 1)
	shouldItems := []ScopeItem{
		{30, 60}, {101, 399}, {600, 200}, {1100, 400},
	}
	should(t, s, shouldItems, 1059)

	//删除掉元素
	s.Remove(101, 399)
	shouldItems = append(shouldItems[:1], shouldItems[2:]...)
	should(t, s, shouldItems, 660)

	s.Remove(500, 700)
	shouldItems = []ScopeItem{{30, 60}, {1200, 300}}
	should(t, s, shouldItems, 360)

	s.Remove(500, 1000)
	shouldItems = []ScopeItem{{30, 60}}
	should(t, s, shouldItems, 60)

	s.Remove(80, 10)
	shouldItems = []ScopeItem{{30, 50}}
	should(t, s, shouldItems, 50)

	s.Remove(40, 10)
	shouldItems = []ScopeItem{{30, 10}, {50, 30}}
	should(t, s, shouldItems, 40)
}

func TestAlloc(t *testing.T) {
	s := NewScope()
	items := []ScopeItem{{550, 250}, {900, 600}, {10, 80}, {100, 400}}

	size := uint64(0)
	for _, v := range items {
		s.Insert(v.Address, v.Size)
		size += v.Size
	}

	//{10, 80}, {100, 400}, {550, 250}, {900, 600},
	items = []ScopeItem{{10, 85}, {100, 400}, {550, 250}, {900, 600}}
	size += 5
	shouldAlloc(t, s, items, size, 5, 90)

	size += 5
	items = []ScopeItem{{10, 490}, {550, 250}, {900, 600}}
	shouldAlloc(t, s, items, size, 5, 95)

	size += 90
	items = []ScopeItem{{10, 490}, {550, 340}, {900, 600}}
	shouldAlloc(t, s, items, size, 90, 800)

	size += 1000
	items = []ScopeItem{{10, 490}, {550, 340}, {900, 1600}}
	shouldAlloc(t, s, items, size, 1000, 1500)
}

func shouldAlloc(t *testing.T, s *Scope, items []ScopeItem, size, sizeAlloc, addrShould uint64) {
	if addr := s.Alloc(sizeAlloc); addr != addrShould {
		t.Errorf("Alloc Scope %v should %v, return %v\n", sizeAlloc, addrShould, addr)
	} else {
		should(t, s, items, size)
	}
}

func TestIsFree(t *testing.T) {
	s := NewScope()
	items := []ScopeItem{{10, 80}, {100, 400}, {550, 250}, {900, 600}}

	for _, v := range items {
		s.Insert(v.Address, v.Size)
	}

	type node struct {
		addr, size uint64
		free       bool
	}
	nodes := []node{
		{5, 4, true}, {5, 5, true}, {10, 80, false}, {90, 10, true}, {95, 2, true},
		{90, 30, false}, {120, 400, false}, {90, 450, false}, {1500, 10, true}, {1510, 10, true},
	}

	for _, v := range nodes {
		free := s.IsFree(v.addr, v.size)
		if free != v.free {
			t.Errorf("%v %v free should be: %v, be: %v\n", v.addr, v.size, v.free, free)
			break
		}
	}
}
