package dev

import (
	"anthodev/codory/internal/actions"
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

func DecodeUUIDv7Action() *actions.Action {
	return &actions.Action{
		ID:          "decode_uuidv7",
		Name:        "Decode UUIDv7",
		Description: "Decode an UUIDv7",
		Type:        actions.ActionTypeFunction,
		Handler:     decodeUUIDv7,
		Arguments: []actions.ActionArgument{
			{
				Name:        "uuidv7",
				Description: "The UUIDv7 to decode (36 characters maximum)",
				Required:    true,
			},
		},
	}
}

func decodeUUIDv7(ctx context.Context) (string, error) {
	args := ctx.Value(actions.ArgsContextKey).([]string)
	id := args[0]

	uuidDatetime, err := uuid7stringToAtom(id)
	if err != nil {
		return "", fmt.Errorf("invalid UUIDv7")
	}

	return fmt.Sprintf("Datetime decoded for %s: %s", id, uuidDatetime), nil
}

func uuid7stringToAtom(uuid string) (string, error) {
	computedUuid := uuid

	if strings.Contains(uuid, "-") {
		computedUuid = strings.ReplaceAll(uuid, "-", "")
	}

	if len(computedUuid) != 32 && len(computedUuid) != 36 {
		return "", fmt.Errorf("bad length")
	}

	tsBytes, err := hex.DecodeString(computedUuid[:12])
	if err != nil {
		return "", err
	}
	ms := uint64(tsBytes[0])<<40 | uint64(tsBytes[1])<<32 |
		uint64(tsBytes[2])<<24 | uint64(tsBytes[3])<<16 |
		uint64(tsBytes[4])<<8 | uint64(tsBytes[5])

	t := time.UnixMilli(int64(ms)).UTC()
	return t.Format(time.RFC3339), nil
}
