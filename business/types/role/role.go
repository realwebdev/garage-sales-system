// Package role represents the role type in the system.
package role

import (
	"encoding/json"
	"fmt"
)

// The set of roles that can be used
var (
	Admin = newRole("Admin")
	User  = newRole("User")
)

// =============================================================================

// Set of known roles
var roles = make(map[string]Role)

// Role represents a role in the system
type Role struct {
	value string
}

func newRole(role string) Role {
	r := Role{role}
	roles[role] = r
	return r
}

// String returns the name of the role.
func (r Role) String() string {
	return r.value
}

// Equal provides the support for the go-cmp package and testing.
func (r Role) Equal(r2 Role) bool {
	return r.value == r2.value
}

// MarshalText provides support for logging and any marshal needs.
func (r Role) MarshalText() ([]byte, error) {
	return []byte(r.value), nil
}

// MarshalJson provides support for json encoding.
func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.value)
}

// UnmarshalJSON provides support for decoding
func (r *Role) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	role, err := Parse(s)
	if err != nil {
		return err
	}

	*r = role
	return nil
}

// Parse parses the string value and returns a role if one exists.
func Parse(value string) (Role, error) {
	role, exists := roles[value]
	if !exists {
		return Role{}, fmt.Errorf("invalid role %q", value)
	}

	return role, nil
}

// MustParse parses the string value and returns a role if one
// exists. If an error occurs the function panics.
func MustParse(value string) Role {
	role, err := Parse(value)
	if err != nil {
		panic(err)
	}

	return role
}

// ParseToString takes a collection of user roles and converts
// them to a slice of string
func ParseToString(usrRoles []Role) []string {
	roles := make([]string, len(usrRoles))
	for i, role := range usrRoles {
		roles[i] = role.String()
	}

	return roles
}

// ParseMany takes a collection of strings and converts them
// to a slice of roles.
func ParseMany(roles []string) ([]Role, error) {
	useRoles := make([]Role, len(roles))
	for i, roleStr := range roles {
		role, err := Parse(roleStr)
		if err != nil {
			return nil, err
		}
		useRoles[i] = role
	}

	return useRoles, nil
}

// Set returns a copy of the set of known roles.
func Set()[]Role {
	var role []Role
	for _, r := range roles {
		role = append(role, r)
	}

	return role
}