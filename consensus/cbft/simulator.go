package cbft

type Simulator struct {
	kinds []uint16
}

func NewSimulator(kinds []uint16) *Simulator {
	return &Simulator{
		kinds: kinds,
	}
}

func (s *Simulator) SL1001() bool {
	return s.contains(1001)
}

func (s *Simulator) SL1002() bool {
	return s.contains(1002)
}

func (s *Simulator) SL1003() bool {
	return s.contains(1003)
}

func (s *Simulator) SL1004() bool {
	return s.contains(1004)
}

func (s *Simulator) SL1005() bool {
	return s.contains(1005)
}

func (s *Simulator) SL1006() bool {
	return s.contains(1006)
}

func (s *Simulator) SL2001() bool {
	return s.contains(2001)
}

func (s *Simulator) SL3001() bool {
	return s.contains(3001)
}

func (s *Simulator) SL3002() bool {
	return s.contains(3002)
}

func (s *Simulator) SL3003() bool {
	return s.contains(3003)
}

func (s *Simulator) SL3004() bool {
	return s.contains(3004)
}

func (s *Simulator) SL3005() bool {
	return s.contains(3005)
}

func (s *Simulator) SL3006() bool {
	return s.contains(3006)
}

func (s *Simulator) SL4001() bool {
	return s.contains(4001)
}

func (s *Simulator) SL5001() bool {
	return s.contains(5001)
}

func (s *Simulator) contains(kind uint16) bool {
	for _, k := range s.kinds {
		if k == kind {
			return true
		}
	}
	return false
}
