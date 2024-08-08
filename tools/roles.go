package tools

const (
	UsualRole = "usual"
	AdminRole = "admin"
)

// HasPermission проверяет есть ли роль в списке
func HasPermission(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
