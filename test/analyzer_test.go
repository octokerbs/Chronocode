package test

import (
	"context"
	"errors"
	"testing"

	"github.com/octokerbs/chronocode-go/internal/repository"
	"github.com/octokerbs/chronocode-go/internal/service"
	"github.com/octokerbs/chronocode-go/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AnalyzerTestSuite struct {
	suite.Suite
	agent          *mocks.GenerativeAgentServiceMock
	repository     *mocks.DatabaseServiceMock
	sourcecodeHost *mocks.SourcecodeHostServiceMock
}

func TestAnalyzerTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyzerTestSuite))
}

func (a *AnalyzerTestSuite) SetupTest() {
	a.agent = &mocks.GenerativeAgentServiceMock{}
	a.repository = &mocks.DatabaseServiceMock{}
	a.sourcecodeHost = &mocks.SourcecodeHostServiceMock{}
}

func (a *AnalyzerTestSuite) TestCannotFetchRepositoryRecordWithInvalidID() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.repository.On("GetRepository").Return(&repository.RepositoryRecord{}, false, errors.New("repository id not found in repository"))

	_, err := service.NewRepositoryAnalyzer(ctx, a.agent, a.sourcecodeHost, a.repository)
	assert.EqualError(a.T(), err, "repository id not found in repository")
}

func (a *AnalyzerTestSuite) TestCanFetchRepositoryRecordWithValidID() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	emptyRepositoryRecord := &repository.RepositoryRecord{}

	a.repository.On("GetRepository").Return(emptyRepositoryRecord, true, nil)

	repoAnalyzer, _ := internal.NewRepositoryAnalyzer(ctx, a.agent, a.sourcecodeHost, a.repository)
	assert.Equal(a.T(), repoAnalyzer, &internal.RepositoryAnalyzer{
		GenerativeAgentService: a.agent,
		SourceCodeHostService:  a.sourcecodeHost,
		DatabaseService:        a.repository,
		RepositoryRecord:       emptyRepositoryRecord,
		AnalyzedCommits:        []repository.CommitRecord{},
	})
}
