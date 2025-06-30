package entities

import (
	"database/sql/driver"
	"fmt"

	"github.com/Gabriel-Schiestl/authgate/internal/src/domain/models"
)

type IdentifierType models.IdentifierType

func (it IdentifierType) Value() (driver.Value, error) {
    return int32(it), nil
}

func (it *IdentifierType) Scan(value interface{}) error {
    if value == nil {
        *it = IdentifierType(models.IdentifierUnspecified)
        return nil
    }
    
    switch v := value.(type) {
    case int64:
        *it = IdentifierType(v)
    case int32:
        *it = IdentifierType(v)
    case int:
        *it = IdentifierType(v)
    default:
        return fmt.Errorf("cannot scan %T into IdentifierType", value)
    }
    
    return nil
}

