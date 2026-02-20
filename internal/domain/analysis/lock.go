package analysis

import (
	"context"
	"errors"
)

var ErrAnalysisInProgress = errors.New("analysis already in progress for this repository")

type Locker interface {
	Acquire(ctx context.Context, repoURL string) (release func(), err error)
	IsLocked(ctx context.Context, repoURL string) bool
}
