package auditbus

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/realwebdev/garage-sales-system/business/types/domain"
	"github.com/realwebdev/garage-sales-system/business/types/name"
)

// Audit represents information about an individual
// audit record.
type Audit struct {
	ID        uuid.UUID
	objID     uuid.UUID
	objDomain domain.Domain
	objName   name.Name
	ActorID   uuid.UUID
	Action    string
	Data      json.RawMessage
	Message   string
	Timestamp time.Time
}

// NewAudit represents the information needed to create a new audit record.
type NewAudit struct {
	ObjID     uuid.UUID
	objDomain domain.Domain
	objName   name.Name
	ActorID   uuid.UUID
	Action    string
	Data      any
	Message   string
}
