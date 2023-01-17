package ports

import "git.carriot.ir/warning-detector/internal/models"

type WarningDetector interface {
	Detect(models.TempLog) (*models.HasWarnings, error)
}
