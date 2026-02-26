package internal

import contractsapi "github.com/delqhi/mikasmissions/platform/libs/contracts-api"

type Repository interface {
	CreateParent(email, country, lang, passwordHash string) (Parent, error)
	FindParentByEmail(email string) (Parent, bool, error)
	UpdateParentLastLogin(parentID string) error
	FindAdminByEmail(email string) (AdminUser, bool, error)
	UpdateAdminLastLogin(adminUserID string) error
	VerifyConsent(parentID, method string) (Consent, error)
	ParentExists(parentID string) (bool, error)
	GetControls(childProfileID string) (contractsapi.ParentalControls, error)
	SetControls(childProfileID string, controls contractsapi.ParentalControls) error
	SaveGateToken(childProfileID, gateToken string) error
	IsValidGateToken(childProfileID, gateToken string) (bool, error)
	ConsumeGateToken(childProfileID, gateToken string) (bool, error)
	CreateGateChallenge(parentUserID, childProfileID, method string) (ParentGateChallenge, error)
	ConsumeGateChallenge(challengeID, parentUserID, childProfileID string) (bool, error)
}
