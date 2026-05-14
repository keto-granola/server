package utils

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDFrom(pgUUID pgtype.UUID) uuid.UUID {
	return uuid.UUID(pgUUID.Bytes)
}

func PGTextFrom(text string) pgtype.Text {
	return pgtype.Text{
		String: text,
		Valid:  true,
	}
}
