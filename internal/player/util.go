package player

import (
	"errors"
	"fmt"
	"mafia-grpc/internal/role"
	"mafia-grpc/internal/util"
	"strings"

	"github.com/google/uuid"
)

func ValidatePlayerName(username string) error {
	if username == "" {
		return errors.New("the user name cannot be empty")
	}
	if strings.Contains(username, "#") {
		return errors.New("the name must not contain #")
	}
	return nil
}

func CopyPlayer(p *Player, hideUuid, hideRole bool) *Player {
	result := &Player{
		Name: p.Name,
		Uuid: p.Uuid,
		Role: p.Role,
		Dead: p.Dead,
	}
	if hideUuid {
		result.Uuid = nil
	}
	if hideRole && !p.Revealed {
		result.Role = role.Role_UNKNOWN
	}
	return result
}

func PlayerToString(p *Player, uuid uuid.UUID) string {
	var result []string
	result = append(result, fmt.Sprintf("Name: %s", p.Name))
	if util.BytesToUUID(p.Uuid) == uuid {
		result = append(result, "(you)")
	}
	result = append(result, fmt.Sprintf(", Role: %s", role.RoleToString(p.Role)))
	if p.Dead {
		result = append(result, ", is dead")
	}
	return strings.Join(result, "")
}
