package adapters

import (
	"context"
	"sync"

	"github.com/octokerbs/chronocode/internal/domain/analysis"
)

type InMemoryLocker struct {
	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

func NewInMemoryLocker() *InMemoryLocker {
	return &InMemoryLocker{locks: make(map[string]*sync.Mutex)}
}

func (l *InMemoryLocker) Acquire(_ context.Context, repoURL string) (func(), error) {
	l.mu.Lock()
	repoLock, exists := l.locks[repoURL]
	if !exists {
		repoLock = &sync.Mutex{}
		l.locks[repoURL] = repoLock
	}
	l.mu.Unlock()

	if !repoLock.TryLock() {
		return nil, analysis.ErrAnalysisInProgress
	}

	return func() { repoLock.Unlock() }, nil
}

func (l *InMemoryLocker) IsLocked(_ context.Context, repoURL string) bool {
	l.mu.Lock()
	repoLock, exists := l.locks[repoURL]
	l.mu.Unlock()

	if !exists {
		return false
	}

	if repoLock.TryLock() {
		repoLock.Unlock()
		return false
	}
	return true
}
