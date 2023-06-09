package role

func RoleToString(r Role) string {
	return Role_name[int32(r)]
}
