// Package SurfDAP provides primitives for navigating an LDAP directory tree.
//
// The package is mainly built around an interface called `Node`, which basically wraps an LDAP
// object. Every node has a DN and functions for easily accessing its parent and children.

package surfdap

import (
	"crypto/tls"
	"fmt"

	"github.com/pkg/errors"
	ldap "gopkg.in/ldap.v2"
)

// Limits
var (
	SizeLimit = 0
	TimeLimit = 0
)

// SearchScope defines the type for an LDAP search scope.
type SearchScope int

// Search scopes
const (
	ScopeBase = SearchScope(ldap.ScopeBaseObject)
	ScopeOne  = SearchScope(ldap.ScopeSingleLevel)
	ScopeSub  = SearchScope(ldap.ScopeWholeSubtree)
)

// Node defines the interface of a SurfDAP node.
type Node interface {
	Conn() *ldap.Conn
	Root() Node
	Parent() Node
	DN() string
	Attributes() map[string][]string
	Children(filter string, attrs []string) ([]Node, error)
	Search(scope SearchScope, filter string, attrs []string) ([]Node, error)
}

// N implements the Node interface.
type N struct {
	dn     string
	conn   *ldap.Conn
	root   *N
	parent *N
}

// New takes a Config{} and returns a newly initialized SurfDAP root node.
func New(config Config) (Node, error) {
	host := fmt.Sprintf("%s:%d", config.Host, int(config.Port))
	l, err := ldap.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	if config.UseStartTLS {
		err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return nil, err
		}
	}

	if config.BindDN != "" {
		err = l.Bind(config.BindDN, config.BindPW)
		if err != nil {
			return nil, err
		}
	}

	newNode := &N{
		dn:     config.BaseDN,
		conn:   l,
		root:   nil,
		parent: nil,
	}

	newNode.Attributes()

	return newNode, nil
}

// String implements the fmt.Stringer
func (n N) String() string {
	out := fmt.Sprintf("%s\n", n.DN())

	for name, values := range n.Attributes() {
		for _, value := range values {
			out += fmt.Sprintf("%s: %s\n", name, value)
		}
	}

	out += "\n"
	return out
}

// Conn returns the underlaying *ldap.Conn.
func (n *N) Conn() *ldap.Conn { return n.conn }

// Root returns the root node that was used to originally bind to the server.
func (n *N) Root() Node { return n.root }

// Parent returns the parent node of this node.
func (n *N) Parent() Node { return n.parent }

// DN returns this node's distinguishable name.
func (n *N) DN() string {
	return n.dn
}

// Attributes returns a map[string][]string with all atributes of this nodes LDAP object.
func (n *N) Attributes() map[string][]string {
	out := map[string][]string{}

	searchResult, err := n.conn.Search(
		ldap.NewSearchRequest(n.DN(), int(ScopeBase), ldap.NeverDerefAliases, SizeLimit, TimeLimit,
			false, "(objectClass=*)", []string{}, nil))
	if err != nil {
		panic(errors.WithStack(err))
	}

	if len(searchResult.Entries) != 1 {
		panic(fmt.Errorf("unable to lookup: %s", n.DN()))
	}

	for _, attr := range searchResult.Entries[0].Attributes {
		out[attr.Name] = attr.Values
	}

	return out
}

func (n *N) newChild(dn string) Node {
	return &N{
		dn:     dn,
		conn:   n.conn,
		root:   n.root,
		parent: n,
	}
}

// Search takes a SearchScope, filter string and attrs []string and performs an LDAP search based
// on their values.
func (n *N) Search(scope SearchScope, filter string, attrs []string) ([]Node, error) {
	out := []Node{}

	searchResult, err := n.conn.Search(
		ldap.NewSearchRequest(n.DN(), int(scope), ldap.NeverDerefAliases, SizeLimit, TimeLimit,
			false, filter, attrs, nil))
	if err != nil {
		return nil, err
	}

	for _, entry := range searchResult.Entries {
		node := n.newChild(entry.DN)
		out = append(out, node)
	}

	return out, err
}

// Children takes a filter string and an attrs []string and performs an LDAP search based on their
// values.
func (n *N) Children(filter string, attrs []string) ([]Node, error) {
	var out = []Node{}
	searchResult, err := n.conn.Search(
		ldap.NewSearchRequest(n.DN(), int(ScopeBase), ldap.NeverDerefAliases, SizeLimit, TimeLimit,
			false, filter, attrs, nil,
		))
	if err != nil {
		return nil, err
	}
	for _, entry := range searchResult.Entries {
		out = append(out, n.newChild(entry.DN))
	}
	return out, nil
}
