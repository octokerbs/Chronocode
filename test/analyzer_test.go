package test

import (
	"context"
	"errors"
	"testing"

	"github.com/chrono-code-hackathon/chronocode-go/internal"
	"github.com/chrono-code-hackathon/chronocode-go/internal/database"
	"github.com/chrono-code-hackathon/chronocode-go/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AnalyzerTestSuite struct {
	suite.Suite
	agent          *mocks.GenerativeAgentServiceMock
	database       *mocks.DatabaseServiceMock
	sourcecodeHost *mocks.SourcecodeHostServiceMock
}

func TestAnalyzerTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyzerTestSuite))
}

func (a *AnalyzerTestSuite) SetupTest() {
	a.agent = &mocks.GenerativeAgentServiceMock{}
	a.database = &mocks.DatabaseServiceMock{}
	a.sourcecodeHost = &mocks.SourcecodeHostServiceMock{}
}

func (a *AnalyzerTestSuite) TestCannotFetchRepositoryRecordWithInvalidID() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.database.On("GetRepository").Return(&database.RepositoryRecord{}, false, errors.New("repository id not found in database"))

	_, err := internal.NewRepositoryAnalyzer(ctx, a.agent, a.sourcecodeHost, a.database)
	assert.EqualError(a.T(), err, "repository id not found in database")
}

func (a *AnalyzerTestSuite) TestCanFetchRepositoryRecordWithValidID() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	emptyRepositoryRecord := &database.RepositoryRecord{}

	a.database.On("GetRepository").Return(emptyRepositoryRecord, true, nil)

	repoAnalyzer, _ := internal.NewRepositoryAnalyzer(ctx, a.agent, a.sourcecodeHost, a.database)
	assert.Equal(a.T(), repoAnalyzer, &internal.RepositoryAnalyzer{
		GenerativeAgentService: a.agent,
		SourceCodeHostService:  a.sourcecodeHost,
		DatabaseService:        a.database,
		RepositoryRecord:       emptyRepositoryRecord,
		AnalyzedCommits:        []database.CommitRecord{},
	})
}
