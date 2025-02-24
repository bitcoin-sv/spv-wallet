package record

type addressInfo struct {
	vouts map[uint32]struct{}
}

type addresses map[string]addressInfo

func (a addresses) append(addr string, vout uint32) {
	info, ok := a[addr]
	if !ok {
		info = addressInfo{vouts: make(map[uint32]struct{})}
		a[addr] = info
	}

	info.vouts[vout] = struct{}{}
}

func (a addresses) contains(addr string, vout uint32) bool {
	info, ok := a[addr]
	if !ok {
		return false
	}

	_, ok = info.vouts[vout]
	return ok
}

func (a addresses) remove(addr string, vout uint32) {
	info, ok := a[addr]
	if !ok {
		return
	}

	delete(info.vouts, vout)
	if len(info.vouts) == 0 {
		delete(a, addr)
	}
}
