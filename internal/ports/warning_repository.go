package ports

import (
	"context"

	"git.carriot.ir/warning-detector/internal/models"
)

type WarningRepository interface {
	Store(context.Context, models.HasWarnings) error
}
