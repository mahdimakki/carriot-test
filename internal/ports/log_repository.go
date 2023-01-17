package ports

import (
	"context"

	"git.carriot.ir/warning-detector/internal/models"
)

type LogRepository interface {
	Store(context.Context, models.TempLog)
}
