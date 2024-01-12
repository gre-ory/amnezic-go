package model

type Permission string

const (
	Permission_Theme Permission = "theme"
	Permission_User  Permission = "user"
)

func (p Permission) String() string {
	return string(p)
}

func ToPermission(value string) Permission {
	switch value {
	case Permission_Theme.String():
		return Permission_Theme
	case Permission_User.String():
		return Permission_User
	}
	return ""
}
