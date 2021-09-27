package reporter

import "github.com/10Pines/tracker/v2/internal/report"

type Reporter func(report report.Report) error
