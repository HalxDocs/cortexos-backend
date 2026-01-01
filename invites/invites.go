package invites

import "time"

type Invite struct {
	Code      string
	ExpiresAt time.Time
	Used      bool
}

var store = map[string]*Invite{}

func CreateInvite(code string, ttl time.Duration) {
	store[code] = &Invite{
		Code:      code,
		ExpiresAt: time.Now().Add(ttl),
		Used:      false,
	}
}

func ValidateInvite(code string) bool {
	invite, ok := store[code]
	if !ok {
		return false
	}

	if invite.Used {
		return false
	}

	if time.Now().After(invite.ExpiresAt) {
		return false
	}

	invite.Used = true
	return true
}
