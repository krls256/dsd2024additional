package auth

import (
	"github.com/samber/lo"
)

func NewPermissionManager() *PermissionManager {
	return &PermissionManager{Permissions: map[string]string{}}
}

type PermissionManager struct {
	Permissions map[string]string
}

func (a *PermissionManager) All() []string {
	return lo.Keys(a.Permissions)
}

func (a *PermissionManager) AllWithDescription() map[string]string {
	return a.Permissions
}

func (a *PermissionManager) Add(permission, description string) bool {
	if _, ok := a.Permissions[permission]; ok {
		return false
	}

	a.Permissions[permission] = description

	return true
}
