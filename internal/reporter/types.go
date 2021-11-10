package reporter

import (
	"errors"
	"log"

	"github.com/10Pines/tracker/v2/internal/shared"
)

// Reporter defines a common set of methods for reporters
type Reporter interface {
	// Name returns the reporter Name
	Name() string
	// Process communicates the report using the underlying transport
	Process(report shared.Report) error
}

type multiple struct {
	reporters []Reporter
}

// Multiple creates a compound Reporter that invokes given reporters sequentially
func Multiple(reporters ...Reporter) Reporter {
	return multiple{reporters: reporters}
}

func (m multiple) Name() string {
	return "Multiple"
}

func (m multiple) Process(report shared.Report) error {
	var err error
	for _, reporter := range m.reporters {
		log.Printf("reporter[%s]", reporter.Name())
		err = reporter.Process(report)
		if err != nil {
			log.Printf("reporter[%s] failed: %v", reporter.Name(), err)
		}
		log.Printf("reporter[%s] done!", reporter.Name())
	}
	if err != nil {
		err = errors.New("reporter failed")
		return err
	}
	return nil
}
