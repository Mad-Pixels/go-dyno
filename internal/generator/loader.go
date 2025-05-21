package generator

import (
	"encoding/json"
	"io"

	"github.com/Mad-Pixels/go-dyno/internal/logger"
)

func Load(r io.Reader) (*DynamoSchema, error) {
	var s DynamoSchema

	if err := json.NewDecoder(r).Decode(&s); err != nil {
		return nil, logger.NewFailure("failed decode config", err)
	}
	return &s, nil
}
