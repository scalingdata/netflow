package session

import (
	"testing"
)

type MockTemplate struct {
	tid uint16
	oid uint32
}

func (m *MockTemplate) ID() uint64 {
	return uint64(m.tid)<<32 | uint64(m.oid)
}

func (m *MockTemplate) TID() uint16 {
	return m.tid
}

func (m *MockTemplate) OID() uint32 {
	return m.oid
}

func TestSession(t *testing.T) {
	m := &MockTemplate{tid: uint16(1), oid: uint32(2)}
	s := New()
	s.AddTemplate(m)
	template, found := s.GetTemplate(m.TID(), m.OID())
	if !found {
		t.Fatal("Couldn't find expected template")
	}
	if template.TID() != uint16(1) || template.OID() != uint32(2) {
		t.Fatalf("Incorrect tid (%d) or oid (%d)", template.TID(), template.OID())
	}
}

func TestSessionSameTidDistinctOid(t *testing.T) {
	m1 := &MockTemplate{tid: uint16(1), oid: uint32(2)}
	m2 := &MockTemplate{tid: uint16(1), oid: uint32(3)}
	s := New()
	s.AddTemplate(m1)
	s.AddTemplate(m2)

	template1, found := s.GetTemplate(m1.TID(), m1.OID())
	if !found {
		t.Fatal("Couldn't find expected template")
	}
	template2, found := s.GetTemplate(m2.TID(), m2.OID())
	if !found {
		t.Fatal("Couldn't find expected template")
	}
	if template1.TID() != template2.TID() {
		t.Fatal("Expected template tids to match")
	}
	if template1.OID() == template2.OID() {
		t.Fatal("Expected template oids to differ")
	}
	if s.combinedID(template1.TID(), template1.OID()) == s.combinedID(template2.TID(), template2.OID()) {
		t.Fatal("Expected session combined IDs to differ")
	}
}
