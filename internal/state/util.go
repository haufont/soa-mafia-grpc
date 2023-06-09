package state

import (
	"fmt"
	"mafia-grpc/internal/player"
	"strings"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

func StateToString(s *State, uuid uuid.UUID) string {
	var result []string

	switch s.PartOfTheDay {
	case PartOfTheDay_DAY:
		result = append(result, color.YellowString("It's daytime now"))
	case PartOfTheDay_NIGHT:
		result = append(result, color.BlackString("It's night now"))
	case PartOfTheDay_UNKNOWN:
	}

	result = append(result, "Players: ")

	for i := range s.Players {
		result = append(result, player.PlayerToString(s.Players[i], uuid))
	}

	if len(s.Voices) != 0 {
		result = append(result, "Voices: ")
		for name, vote := range s.Voices {
			if vote == "" {
				result = append(result, fmt.Sprintf("\"%s\": apstained", name))
			} else {
				result = append(result, fmt.Sprintf("\"%s\": voted for \"%s\"", name, vote))
			}
		}
	}

	return strings.Join(result, "\n")
}
