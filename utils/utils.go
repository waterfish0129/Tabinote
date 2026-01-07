package utils

import "github.com/jackc/pgx/v5/pgtype"

func StringPtr(s string) *string {
	return &s
}
func ToPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}
