package command

import (
	"context"
	"errors"
	"testing"

	"github.com/octokerbs/chronocode/internal/adapters"
	"github.com/octokerbs/chronocode/internal/domain/agent"
	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
	"github.com/octokerbs/chronocode/internal/domain/subcommit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AnalyzeRepositoryTestSuite struct {
	suite.Suite
	repoRepository      repo.Repository
	subcommitRepository subcommit.Repository
	agent               agent.Agent
	codeHost            codehost.CodeHost
	handler             AnalyzeRepoHandler
}

func TestAnalyzeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyzeRepositoryTestSuite))
}

func (s *AnalyzeRepositoryTestSuite) SetupTest() {
	s.repoRepository = adapters.NewRepoRepository()
	s.subcommitRepository = adapters.NewSubcommitRepository()
	s.agent = adapters.NewAgent()
	s.codeHost = adapters.NewCodeHost()
	s.handler = NewAnalyzeRepoHandler(s.repoRepository, s.subcommitRepository, s.agent, s.codeHost)
}

func (s *AnalyzeRepositoryTestSuite) TestCannotAnalyzeWithInvalidURL() {
	err := s.handler.Handle(context.Background(), AnalyzeRepo{adapters.InvalidRepoURL})
	assert.True(s.T(), errors.Is(err, codehost.ErrInvalidRepoURL))
}

func (s *AnalyzeRepositoryTestSuite) TestAnalyzesValidRepoSuccessfully() {
	err := s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL})
	assert.Nil(s.T(), err)
}

func (s *AnalyzeRepositoryTestSuite) TestStoresNewRepositoryAfterAnalysis() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL})
	_, err := s.repoRepository.GetRepo(context.Background(), adapters.ValidRepoURL)
	assert.Nil(s.T(), err)
}

func (s *AnalyzeRepositoryTestSuite) TestNewRepoHasSubcommitsAfterAnalysis() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL})
	subcommits, err := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidRepoID)

	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), subcommits)
}

func (s *AnalyzeRepositoryTestSuite) TestNewRepoWithoutCommitsHasNoSubcommits() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidEmptyRepoURL})
	subcommits, err := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidEmptyRepoID)

	assert.Nil(s.T(), err)
	assert.Empty(s.T(), subcommits)
}

func (s *AnalyzeRepositoryTestSuite) TestExistingRepositoryMayHaveOutdatedSubcommits() {
	_ = s.repoRepository.StoreRepo(context.Background(), repo.NewRepo(adapters.ValidRepoID, "chronocode", adapters.ValidRepoURL, "old-sha"))
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL})
	subcommits, err := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidRepoID)

	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), subcommits)
}

func (s *AnalyzeRepositoryTestSuite) TestExistingRepoSubcommitsAreAddedToExistingOnes() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL})
	subcommitsBefore, _ := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidRepoID)

	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL})
	subcommitsAfter, _ := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidRepoID)

	assert.Greater(s.T(), len(subcommitsAfter), len(subcommitsBefore))
}
