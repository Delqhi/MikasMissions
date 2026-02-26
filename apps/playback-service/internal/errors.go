package internal

import "errors"

var ErrParentGateRequired = errors.New("parent gate verification required")
var ErrSessionLimitTooHigh = errors.New("session limit exceeds strict policy")
var ErrSessionCapReached = errors.New("session cap reached for current profile")
var ErrEntitlementRequired = errors.New("active entitlement required")
var ErrAutoplayBlocked = errors.New("autoplay blocked by strict safety mode")
var ErrChildProfileForbidden = errors.New("child profile access is not allowed for current principal")
