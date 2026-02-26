package internal

import (
	"fmt"
	"net/http"

	"github.com/delqhi/mikasmissions/platform/libs/authz"
	"github.com/delqhi/mikasmissions/platform/libs/profileclient"
)

func parentUserIDFromRequest(r *http.Request, requested string) (string, error) {
	principal, ok := authz.PrincipalFrom(r.Context())
	if !ok || principal.Role != "parent" {
		return requested, nil
	}
	if principal.ParentUserID == "" {
		return "", fmt.Errorf("missing parent identity")
	}
	if requested != "" && requested != principal.ParentUserID {
		return "", fmt.Errorf("parent_user_id does not match authenticated parent")
	}
	return principal.ParentUserID, nil
}

func ensureChildOwnership(r *http.Request, parentUserID, childProfileID string) error {
	if parentUserID == "" || childProfileID == "" {
		return nil
	}
	principal, ok := authz.PrincipalFrom(r.Context())
	if !ok || principal.Role != "parent" {
		return nil
	}
	reader, err := profileclient.NewFromEnv()
	if err != nil {
		return err
	}
	owned, err := reader.IsOwnedByParent(r.Context(), parentUserID, childProfileID)
	if err != nil {
		return err
	}
	if !owned {
		return fmt.Errorf("child profile access is not allowed for current parent")
	}
	return nil
}
