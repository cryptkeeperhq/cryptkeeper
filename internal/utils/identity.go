package utils

type Identity interface {
	GetID() string
	GetUsername() string
	GetAuthType() string
	GetGroups() []string
}

type UserIdentity struct {
	UserID   string
	Username string
	Groups   []string
}

func (u UserIdentity) GetID() string {
	return u.UserID
}

func (u UserIdentity) GetUsername() string {
	return u.Username
}

func (u UserIdentity) GetGroups() []string {
	return u.Groups
}

func (u UserIdentity) GetAuthType() string {
	return "user"
}

type AppRoleIdentity struct {
	AppRoleID string
	RoleName  string
}

func (a AppRoleIdentity) GetID() string {
	return a.AppRoleID
}

func (a AppRoleIdentity) GetUsername() string {
	return a.RoleName
}

func (a AppRoleIdentity) GetGroups() []string {
	return []string{}
}

func (a AppRoleIdentity) GetAuthType() string {
	return "approle"
}
