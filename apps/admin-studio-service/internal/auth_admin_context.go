package internal

import "github.com/delqhi/mikasmissions/platform/libs/authz"

func actorIDFromPrincipal(principal authz.Principal) string {
	if principal.ParentUserID != "" {
		return principal.ParentUserID
	}
	return "admin-system"
}
