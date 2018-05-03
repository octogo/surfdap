package surfdap

import (
	"fmt"
	"os"

	ldap "gopkg.in/ldap.v2"
)

// DisplaySearchResults writes search results to an os.File.
func DisplaySearchResults(scope SearchScope, results *ldap.SearchResult, w *os.File) {
	write := func(lines ...string) {
		for _, line := range lines {
			w.Write([]byte(line))
		}
		w.Write([]byte("\n"))
	}

	switch scope {
	case ScopeBase:
		for _, entry := range results.Entries {
			write("dn:", entry.DN)
			for _, attr := range entry.Attributes {
				for _, value := range attr.Values {
					write(fmt.Sprintf("%s: %s", attr.Name, value))
				}
			}
		}
		write("")

	case ScopeOne, ScopeSub:
		for _, entry := range results.Entries {
			write("dn:", entry.DN)
			for _, attr := range entry.Attributes {
				for _, value := range attr.Values {
					write(fmt.Sprintf("%s: %s", attr.Name, value))
				}
			}
			write("")
		}
		write("")
	}
}
