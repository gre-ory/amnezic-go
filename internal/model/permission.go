package model

type Permission int

const (
	_ Permission = iota

	Permission_Theme
	Permission_Music
	Permission_User
	Permission_Session
	Permission_File

	__Permission_last
)

func AllPermissions() []Permission {
	permissions := make([]Permission, 0, int(__Permission_last-1))
	for i := 1; i < int(__Permission_last); i++ {
		permissions = append(permissions, Permission(i))
	}
	return permissions
}

func (p Permission) String() string {
	switch p {
	case Permission_Theme:
		return "theme"
	case Permission_Music:
		return "music"
	case Permission_User:
		return "user"
	case Permission_Session:
		return "session"
	case Permission_File:
		return "file"
	}
	return ""
}

func ToPermission(value string) Permission {
	if value == "" {
		return 0
	}
	for i := 1; i < int(__Permission_last); i++ {
		if Permission(i).String() == value {
			return Permission(i)
		}
	}
	return 0
}
