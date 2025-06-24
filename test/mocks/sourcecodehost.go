package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type SourcecodeHostServiceMock struct {
	mock.Mock
}

func (s *SourcecodeHostServiceMock) ProduceCommits(ctx context.Context, lastAnalyzedCommitSHA string, commits chan<- string, errors chan<- error) {

}

func (s *SourcecodeHostServiceMock) GetCommitDiff(ctx context.Context, commitSHA string) (string, error) {
	return "", nil
}

func (s *SourcecodeHostServiceMock) RepositoryID() int64 {
	return 0
}

func (s *SourcecodeHostServiceMock) GetRepositoryData() map[string]interface{} {
	return nil
}

func (s *SourcecodeHostServiceMock) GetCommitData(ctx context.Context, commitSHA string) map[string]interface{} {
	return nil
}
