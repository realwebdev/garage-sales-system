package garage.rego

import rego.v1

role_user := "USER"

role_admin := "ADMIN"

role_all := {role_admin, role_user}

default rule_any := false

rule_any if {
    claim_roles := {role | some role in input.Roles}
    input_roles := roles_all & claim_roles
    count(input_roles) > 0
}