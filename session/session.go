// Package session provides sessions for the Netflow version 9 and IPFIX
// decoders that need to track templates bound to a session.
package session

import "sync"

type Template interface {
	TID() uint16
	OID() uint32
}

type Session interface {
	Lock()
	Unlock()

	// To keep track of maximum record sizes per template
	GetRecordSize(uint16, uint32) (size int, found bool)
	SetRecordSize(uint16, uint32, int)

	// To keep track of templates
	AddTemplate(Template)
	GetTemplate(uint16, uint32) (t Template, found bool)
}

type basicSession struct {
	mutex     *sync.Mutex
	templates map[uint64]Template
	sizes     map[uint64]int
}

func New() *basicSession {
	return &basicSession{
		mutex:     &sync.Mutex{},
		templates: make(map[uint64]Template, 65536),
		sizes:     make(map[uint64]int, 65536),
	}
}

func (s *basicSession) Lock() {
	s.mutex.Lock()
}

func (s *basicSession) Unlock() {
	s.mutex.Unlock()
}

func (s *basicSession) combinedID(tid uint16, oid uint32) uint64 {
	return uint64(tid)<<32 | uint64(oid)
}

func (s *basicSession) GetRecordSize(tid uint16, oid uint32) (size int, found bool) {
	size, found = s.sizes[s.combinedID(tid, oid)]
	return
}

func (s *basicSession) SetRecordSize(tid uint16, oid uint32, size int) {
	id := s.combinedID(tid, oid)
	if s.sizes[id] < size {
		s.sizes[id] = size
	}
}

func (s *basicSession) AddTemplate(t Template) {
	id := s.combinedID(t.TID(), t.OID())
	s.templates[id] = t
}

func (s *basicSession) GetTemplate(tid uint16, oid uint32) (t Template, found bool) {
	id := s.combinedID(tid, oid)
	t, found = s.templates[id]
	return
}

// Test if basicSession is compliant
var _ Session = (*basicSession)(nil)
