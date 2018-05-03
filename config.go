package surfdap

// Config defines a node configuration.
type Config struct {
	Host           string // LDAP host to connect to
	Port           uint16 // LDAP port to connect to
	UseStartTLS    bool   // Use StartTLS when connecting
	BindDN, BindPW string // BindDN and BindPW (optional)
	BaseDN         string // search-base DN
}
