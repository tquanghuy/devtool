package checker

import (
	"context"
	"devtool/internal/config"
	"sync"
)

type Checker interface {
	Check(ctx context.Context, cfg *config.AppConfig) StatusResult
}

type Orchestrator struct {
	checkers []Checker
	cfg      *config.AppConfig
}

func NewOrchestrator(cfg *config.AppConfig) *Orchestrator {
	return &Orchestrator{
		checkers: make([]Checker, 0),
		cfg:      cfg,
	}
}

func (o *Orchestrator) Register(c Checker) {
	o.checkers = append(o.checkers, c)
}

func (o *Orchestrator) RunAll(ctx context.Context) []StatusResult {
	results := make([]StatusResult, len(o.checkers))
	var wg sync.WaitGroup

	for i, c := range o.checkers {
		wg.Add(1)
		go func(index int, ch Checker) {
			defer wg.Done()
			results[index] = ch.Check(ctx, o.cfg)
		}(i, c)
	}

	wg.Wait()
	return results
}
