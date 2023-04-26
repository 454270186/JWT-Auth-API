package domain

import "strings"

type RolePermissions struct {
	rolePermissions map[string][]string
}

func GetRolePermissions() RolePermissions {
	return RolePermissions{map[string][]string{
		"admin": {"GetAllCustomers", "GetCustomer", "NewAccount", "NewTransaction"},
		"user": {"GetCutsomer", "NewTransaction"},
	}}
}

func (p RolePermissions) IsAuthorizedFor(role string, routeName string) bool {
	perms := p.rolePermissions[role]

	for _, v := range perms {
		if v == strings.TrimSpace(routeName) {
			return true
		}
	}

	return false
}
