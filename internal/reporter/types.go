package reporter

import "github.com/10Pines/tracker/internal/report"

type Reporter func(report report.Report) error
