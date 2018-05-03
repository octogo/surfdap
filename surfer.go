// Package surfdap provides primitives for navigating an LDAP directory tree.
package surfdap

import (
	"crypto/tls"
	"fmt"

	ldap "gopkg.in/ldap.v2"
)

// Limits
var (
	SizeLimit = 0
	TimeLimit = 0
)

// Filter defines the type of a search filter.
type Filter string

// Scope defines the type for an LDAP search scope.
type Scope int

// Available search scopes
const (
	Base = Scope(ldap.ScopeBaseObject)
	One  = Scope(ldap.ScopeSingleLevel)
	Sub  = Scope(ldap.ScopeWholeSubtree)
)

// OnlyAttrs defines the type of an attribute filter.
type OnlyAttrs []string

// Surfer defines the interface of an LDAP surfer.
type Surfer interface {
	Lookup(Scope, Filter, OnlyAttrs) ([]Surfer, error)
	Entry() *ldap.Entry
	Parent() Surfer
}

// S implements the Surfer interface.
type S struct {
	conn   *ldap.Conn
	entry  *ldap.Entry
	parent Surfer
}

// New returns a newly initialized Surfer.
func New(host string, port uint16, startTLS bool, baseDn, bindDn, bindPW string) (*S, error) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	if startTLS {
		err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return nil, err
		}
	}

	err = l.Bind(bindDn, bindPW)
	if err != nil {
		return nil, err
	}

	s := &S{
		conn: l,
		entry: &ldap.Entry{
			DN: baseDn,
		},
	}

	entries, err := s.lookup(Base, "(objectClass=*)", []string{})
	if err != nil {
		return nil, err
	}
	s.entry = entries[0]

	return s, nil
}

func (s S) String() string {
	out := fmt.Sprintf("dn: %s\n", s.entry.DN)

	for _, attr := range s.entry.Attributes {
		for _, value := range attr.Values {
			out += fmt.Sprintf("%s: %s\n", attr.Name, value)
		}
	}

	return out
}

func (s *S) lookup(scope Scope, f Filter, a OnlyAttrs) ([]*ldap.Entry, error) {
	searchRequest := ldap.NewSearchRequest(s.entry.DN, int(scope), ldap.NeverDerefAliases, 0, 0,
		false, string(f), a, nil)

	searchResult, err := s.conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	return searchResult.Entries, nil
}

// Lookup performs an LDAP search and returns a set of new surfers.
func (s *S) Lookup(scope Scope, filter Filter, attrs OnlyAttrs) ([]Surfer, error) {
	if attrs == nil {
		attrs = []string{}
	}

	out := []Surfer{}

	entries, err := s.lookup(scope, filter, attrs)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		out = append(out, &S{
			conn:  s.conn,
			entry: entry,
		})
	}

	return out, nil
}

// Entry returns the *ldap.Entry of this node.
func (s *S) Entry() *ldap.Entry {
	return s.entry
}

// Parent returns the parent node of this node.
func (s *S) Parent() Surfer {
	return s.parent
}
