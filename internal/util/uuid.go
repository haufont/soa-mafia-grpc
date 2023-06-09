package util

import "github.com/google/uuid"

func UUIDToBytes(uuid uuid.UUID) []byte {
	return []byte(uuid[:16])
}

func BytesToUUID(b []byte) (uuid uuid.UUID) {
	copy(uuid[:], b)
	return
}

func NewUUID() (id uuid.UUID, bid []byte) {
	id = uuid.New()
	bid = UUIDToBytes(id)
	return
}
