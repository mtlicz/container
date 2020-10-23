package container

import (
	"container/list"
)

type ScopeItem struct {
	Address, Size uint64
}

type Scope struct {
	items *list.List
	size  uint64
}

func (p *Scope) Alloc(size uint64) (addr uint64) {
	items := p.items
	if items == nil || items.Len() == 0 {
		p.Insert(addr, size)
	} else {
		var e *list.Element
		for e = items.Front(); e != nil; {
			next := e.Next()
			if next != nil {
				i, _ := e.Value.(ScopeItem)
				j, _ := next.Value.(ScopeItem)
				iend := i.Address + i.Size

				if j.Address-iend > size {
					addr = iend
					i.Size += size
					e.Value = i
					break
				} else if j.Address-iend == size {
					addr = iend
					i.Size += size + j.Size
					e.Value = i
					items.Remove(next)
					break
				}
			}

			e = next
		}

		if e == nil {
			e = items.Back()
			i := e.Value.(ScopeItem)
			addr = i.Address + i.Size

			i.Size += size
			e.Value = i
		}
		p.size += size
	}

	return
}

func (p *Scope) Insert(addr, size uint64) {
	if p.items == nil {
		p.items = list.New()
	}

	items := p.items
	item := ScopeItem{addr, size}
	if items.Len() == 0 {
		items.PushBack(item)
		p.size += size
	} else {
		var e *list.Element
		var i ScopeItem

		for e = items.Front(); e != nil; e = e.Next() {
			i = e.Value.(ScopeItem)
			if addr < i.Address {
				break
			}
		}

		if e == nil { //should push back
			e = items.Back()
			if i.Address+i.Size >= addr {
				p.size += addr + size - (i.Address + i.Size)
				i.Size = addr + size - i.Address
				e.Value = i
			} else {
				items.PushBack(item)
				p.size += size
			}
		} else {
			if addr+size >= i.Address { //可以和后面的合并
				p.size += i.Address - addr
				i.Size = i.Address + i.Size - addr
				i.Address = addr
				e.Value = i
			} else {
				e = items.InsertBefore(item, e)
				i = item
				p.size += size
			}

			front := items.Front()
			if e != front { //看是否能和前面的合并
				prev := e.Prev()
				j := prev.Value.(ScopeItem)

				if j.Address+j.Size >= addr { //可以和前面的合并
					j.Size = i.Address + i.Size - j.Address
					prev.Value = j
					items.Remove(e)
				}
			}
		}
	}
}

func (p *Scope) Remove(addr, size uint64) {
	items := p.items
	if items != nil && items.Len() > 0 {
		end := addr + size

		for e := items.Front(); e != nil && addr < end; {
			i := e.Value.(ScopeItem)
			ie := i.Address + i.Size
			needNext := true

			if addr >= ie {
			} else if addr >= i.Address && addr < ie {
				if end <= ie {
					p.size -= end - addr
					needNext = false

					if addr == i.Address {
						if end == ie {
							items.Remove(e)
						} else {
							i.Address = end
							i.Size = ie - end
							e.Value = i
						}
					} else if end == ie {
						i.Size = addr - i.Address
						e.Value = i
					} else { //addr > i.Address, end < ie
						i.Size = addr - i.Address
						e.Value = i

						i.Address = end
						i.Size = ie - end
						items.InsertAfter(i, e)
					}
					break
				} else if addr == i.Address {
					p.size -= i.Size
					c := e
					e = e.Next()
					items.Remove(c)
					needNext = false
					addr = ie
				} else {
					p.size -= ie - addr
					i.Size = addr - i.Address
					e.Value = i
					addr = ie
				}
			} else if addr < i.Address && end >= ie {
				c := e
				e = e.Next()
				items.Remove(c)
				needNext = false
				p.size -= i.Size
				addr = ie
			} else if end > i.Address && end <= ie {
				addr = end
				if end == ie {
					c := e
					e = e.Next()
					items.Remove(c)
					needNext = false
					p.size -= i.Size
				} else {
					p.size -= end - i.Address
					i.Size = ie - end
					i.Address = end
					e.Value = i
				}
			}

			if needNext {
				e = e.Next()
			}
		}
	}
}

func (p *Scope) IsFree(addr, size uint64) bool {
	ret := true
	items := p.items

	if items != nil && items.Len() > 0 {
		end := addr + size
		e := items.Front()
		i := e.Value.(ScopeItem)
		if end <= i.Address {
			return true
		}

		eend := items.Back()
		j := eend.Value.(ScopeItem)
		if addr >= j.Address+j.Size {
			return true
		}

		ret = false
		for e != eend {
			next := e.Next()
			j = next.Value.(ScopeItem)

			if addr >= i.Address+i.Size && end <= j.Address {
				ret = true
				break
			}

			i = j
			e = next
			if end <= j.Address {
				break
			}
		}
	}

	return ret
}

func isOverlap(a, b, c, d uint64) bool {
	return (a == c && b == d) || (a > c && a < d) || (b > c && b < d) || (c > a && c < b) || (d > a && d < b)
}

func (p *Scope) Size() uint64 {
	return p.size
}

func NewScope() *Scope {
	return &Scope{}
}
