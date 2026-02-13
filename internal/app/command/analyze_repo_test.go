package command

import (
	"context"
	"errors"
	"testing"

	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
	"github.com/octokerbs/chronocode/testing/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AnalyzeRepositoryTestSuite struct {
	suite.Suite
	repoRepository repo.Repository
	codeHost       codehost.CodeHost
	handler        AnalyzeRepoHandler
}

func TestAnalyzeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyzeRepositoryTestSuite))
}

func (s *AnalyzeRepositoryTestSuite) SetupTest() {
	s.repoRepository = mocks.NewRepoRepository()
	s.codeHost = mocks.NewCodeHost()
	s.handler = NewAnalyzeRepoHandler(s.repoRepository, s.codeHost)
}

func (s *AnalyzeRepositoryTestSuite) TestCannotAnalyzeWithInvalidURL() {
	err := s.handler.Handle(context.Background(), AnalyzeRepo{mocks.InvalidRepoURL})
	assert.True(s.T(), errors.Is(err, codehost.ErrInvalidRepoURL))
}

func (s *AnalyzeRepositoryTestSuite) TestAnalyzesValidRepoSuccessfully() {
	err := s.handler.Handle(context.Background(), AnalyzeRepo{mocks.ValidRepoURL})
	assert.Nil(s.T(), err)
}

func (s *AnalyzeRepositoryTestSuite) TestStoresNewRepositoryAfterAnalysis() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{mocks.ValidRepoURL})
	_, err := s.repoRepository.GetRepo(context.Background(), mocks.ValidRepoURL)
	assert.Nil(s.T(), err)
}

//func (s *AnalyzeRepositoryTestSuite) Test04ExistingRepositoryMayHaveOutdatedCommits() {
//
//}
//
//func (s *AnalyzeRepositoryTestSuite) Test05NewRepoHasSubcommitsAfterAnalysis() {
//
//}
//
//func (s *AnalyzeRepositoryTestSuite) Test06ExistingRepoSubcommitsAreAddedToExistingOnes() {
//
//}
