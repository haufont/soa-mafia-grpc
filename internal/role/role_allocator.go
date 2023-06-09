package role

import (
	"errors"
	"math/rand"
)

type RoleAllocator struct {
	roles []Role
}

func NewRoleAllocator(numberOfPlayers uint) RoleAllocator {
	roles := make(map[Role]uint)
	roles[Role_MAFIA] = numberOfPlayers / 3
	roles[Role_COMMISAR] = 1
	roles[Role_TOWNIE] = numberOfPlayers - roles[Role_MAFIA] - roles[Role_COMMISAR]

	rolesList := make([]Role, numberOfPlayers)
	pos := 0
	for role, n := range roles {
		for i := uint(0); i < n; i++ {
			rolesList[pos] = role
			pos++
		}
	}
	rand.Shuffle(int(numberOfPlayers), func(i, j int) {
		rolesList[i], rolesList[j] = rolesList[j], rolesList[i]
	})
	return RoleAllocator{rolesList}
}

func (a *RoleAllocator) Allocate() (role Role, err error) {
	if len(a.roles) == 0 {
		err = errors.New("empty role allocator")
		return
	}
	role = a.roles[0]
	a.roles = a.roles[1:]
	return
}

func (a *RoleAllocator) Rollback(role Role) {
	a.roles = append(a.roles, role)
	rand.Shuffle(len(a.roles), func(i, j int) {
		a.roles[i], a.roles[j] = a.roles[j], a.roles[i]
	})
}
