package internal

type Repository interface {
	CreateProfile(parentID, displayName, ageBand, avatar string) (ChildProfile, error)
	ListProfilesByParent(parentID string) ([]ChildProfile, error)
	FindProfile(id string) (ChildProfile, bool, error)
}
