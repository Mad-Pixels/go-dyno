package schema

import (
	"strings"

	"github.com/Mad-Pixels/go-dyno/internal/generator/attribute"
	"github.com/Mad-Pixels/go-dyno/internal/generator/index"
	"github.com/Mad-Pixels/go-dyno/internal/logger"
)

// Validate performs comprehensive schema validation.
//
// This includes:
//   - Validation of all attributes
//   - Verification that hash/range keys are defined
//   - Validation of index names and definitions
//   - Enforcement of LSI limits
//   - Parsing of composite key definitions
//
// Returns an error if any invalid configuration is found.
func (s *Schema) Validate() error {
	for _, attr := range s.AllAttributes() {
		if err := attr.Validate(); err != nil {
			return err
		}
	}

	if !isAttributeDefined(s.HashKey(), s.AllAttributes()) {
		return logger.NewFailure("hash_key is not defined in attributes", nil).
			With("key", s.HashKey())
	}
	if rk := s.RangeKey(); rk != "" && !isAttributeDefined(rk, s.AllAttributes()) {
		return logger.NewFailure("range_key is not defined in attributes", nil).
			With("key", rk)
	}
	if err := s.ValidateIndexNames(); err != nil {
		return err
	}

	lsiCount := 0
	for i := range s.raw.SecondaryIndexes {
		idx := &s.raw.SecondaryIndexes[i]

		if idx.Type == "" {
			idx.Type = index.GSI
		}
		if idx.IsLSI() {
			idx.HashKey = s.HashKey()
		}
		if err := idx.Validate(s.HashKey(), s.RangeKey()); err != nil {
			return err
		}

		if idx.IsLSI() {
			lsiCount++
			if lsiCount > 10 {
				return logger.NewFailure("too many LSI indexes", nil).
					With("count", lsiCount).
					With("limit", 10)
			}
		}

		if err := validateIndexAttributes(idx, s.AllAttributes()); err != nil {
			return err
		}
		if err := parseIndexCompositeKeys(idx, s.AllAttributes()); err != nil {
			return err
		}
	}
	return nil
}

func isAttributeDefined(name string, attrs []attribute.Attribute) bool {
	for _, a := range attrs {
		if a.Name == name {
			return true
		}
	}
	return false
}

func validateIndexAttributes(idx *index.Index, attrs []attribute.Attribute) error {
	if idx.IsGSI() && idx.HashKey != "" {
		if !isAttributeDefined(idx.HashKey, attrs) && !strings.Contains(idx.HashKey, "#") {
			return logger.NewFailure("GSI hash_key is not defined", nil).
				With("key", idx.HashKey)
		}
	}
	if idx.RangeKey != "" {
		if !isAttributeDefined(idx.RangeKey, attrs) && !strings.Contains(idx.RangeKey, "#") {
			return logger.NewFailure("range_key is not defined", nil).
				With("key", idx.RangeKey)
		}
	}
	for _, nk := range idx.NonKeyAttributes {
		if !isAttributeDefined(nk, attrs) {
			return logger.NewFailure("non_key_attribute is not defined", nil).
				With("key", nk)
		}
	}
	return nil
}

func parseIndexCompositeKeys(idx *index.Index, attrs []attribute.Attribute) error {
	if idx.IsGSI() && strings.Contains(idx.HashKey, "#") {
		parts := strings.Split(idx.HashKey, "#")
		for _, p := range parts {
			if !isAttributeDefined(p, attrs) {
				return logger.NewFailure("invalid composite part in hash_key", nil).
					With("key", p)
			}
		}
		idx.HashKeyParts = make([]index.CompositeKey, len(parts))
		for j, p := range parts {
			idx.HashKeyParts[j] = index.CompositeKey{IsConstant: !isAttributeDefined(p, attrs), Value: p}
		}
	}

	if idx.IsLSI() && strings.Contains(idx.HashKey, "#") {
		parts := strings.Split(idx.HashKey, "#")
		for _, p := range parts {
			if !isAttributeDefined(p, attrs) {
				return logger.NewFailure("invalid composite part in LSI hash_key", nil).
					With("key", p)
			}
		}
		idx.HashKeyParts = make([]index.CompositeKey, len(parts))
		for j, p := range parts {
			idx.HashKeyParts[j] = index.CompositeKey{IsConstant: !isAttributeDefined(p, attrs), Value: p}
		}
	}

	if strings.Contains(idx.RangeKey, "#") {
		parts := strings.Split(idx.RangeKey, "#")
		for _, p := range parts {
			if !isAttributeDefined(p, attrs) {
				indexType := "range_key"

				if idx.IsLSI() {
					indexType = "LSI range_key"
				} else if idx.IsGSI() {
					indexType = "GSI range_key"
				}

				return logger.NewFailure("invalid composite part in "+indexType, nil).
					With("key", p)
			}
		}
		idx.RangeKeyParts = make([]index.CompositeKey, len(parts))
		for j, p := range parts {
			idx.RangeKeyParts[j] = index.CompositeKey{IsConstant: !isAttributeDefined(p, attrs), Value: p}
		}
	}
	return nil
}
