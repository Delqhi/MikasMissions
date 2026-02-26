package authz

type Principal struct {
	ParentUserID string
	Role         string
}

func (p Principal) IsParent() bool {
	return p.Role == "parent"
}
