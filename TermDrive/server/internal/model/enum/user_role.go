package enum

type Role string

const (
	ROLE_ADMIN Role = "ADMIN"
	ROLE_USER  Role = "USER"
)

// IsValid checks if the role is valid within the defined constants.
// Returns true if the role is either ADMIN or USER, otherwise false.
func (r Role) IsValid() bool {
	switch r {
	case ROLE_ADMIN, ROLE_USER:
		return true
	default:
		return false
	}
}
