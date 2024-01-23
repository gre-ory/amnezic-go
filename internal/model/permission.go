package model

type Permission string

const (
	Permission_Theme   Permission = "theme"
	Permission_Music   Permission = "music"
	Permission_User    Permission = "user"
	Permission_Session Permission = "session"
)

func (p Permission) String() string {
	return string(p)
}

func ToPermission(value string) Permission {
	switch value {
	case Permission_Theme.String():
		return Permission_Theme
	case Permission_Music.String():
		return Permission_Music
	case Permission_User.String():
		return Permission_User
	case Permission_Session.String():
		return Permission_Session
	}
	return ""
}
