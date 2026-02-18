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
	codeHostFactory     codehost.CodeHostFactory
	handler             AnalyzeRepoHandler
}

func TestAnalyzeRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AnalyzeRepositoryTestSuite))
}

func (s *AnalyzeRepositoryTestSuite) SetupTest() {
	s.repoRepository = adapters.NewRepoRepository()
	s.subcommitRepository = adapters.NewSubcommitRepository()
	s.agent = adapters.NewAgent()
	s.codeHostFactory = adapters.NewCodeHostFactory()
	s.handler = NewAnalyzeRepoHandler(s.repoRepository, s.subcommitRepository, s.agent, s.codeHostFactory)
}

func (s *AnalyzeRepositoryTestSuite) TestCannotAnalyzeWithoutAccessToken() {
	err := s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL, ""})
	assert.NotNil(s.T(), err)
}

func (s *AnalyzeRepositoryTestSuite) TestCannotAnalyzeInaccessibleRepo() {
	err := s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ForbiddenRepoURL, adapters.ValidAccessToken})
	assert.True(s.T(), errors.Is(err, codehost.ErrAccessDenied))
}

func (s *AnalyzeRepositoryTestSuite) TestCannotAnalyzeWithInvalidURL() {
	err := s.handler.Handle(context.Background(), AnalyzeRepo{adapters.InvalidRepoURL, adapters.ValidAccessToken})
	assert.True(s.T(), errors.Is(err, codehost.ErrInvalidRepoURL))
}

func (s *AnalyzeRepositoryTestSuite) TestAnalyzesValidRepoSuccessfully() {
	err := s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL, adapters.ValidAccessToken})
	assert.Nil(s.T(), err)
}

func (s *AnalyzeRepositoryTestSuite) TestStoresNewRepositoryAfterAnalysis() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL, adapters.ValidAccessToken})
	_, err := s.repoRepository.GetRepo(context.Background(), adapters.ValidRepoURL)
	assert.Nil(s.T(), err)
}

func (s *AnalyzeRepositoryTestSuite) TestNewRepoHasSubcommitsAfterAnalysis() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL, adapters.ValidAccessToken})
	subcommits, err := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidRepoID)

	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), subcommits)
}

func (s *AnalyzeRepositoryTestSuite) TestNewRepoWithoutCommitsHasNoSubcommits() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidEmptyRepoURL, adapters.ValidAccessToken})
	subcommits, err := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidEmptyRepoID)

	assert.Nil(s.T(), err)
	assert.Empty(s.T(), subcommits)
}

func (s *AnalyzeRepositoryTestSuite) TestExistingRepositoryMayHaveOutdatedSubcommits() {
	_ = s.repoRepository.StoreRepo(context.Background(), repo.NewRepo(adapters.ValidRepoID, "chronocode", adapters.ValidRepoURL, "old-sha"))
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL, adapters.ValidAccessToken})
	subcommits, err := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidRepoID)

	assert.Nil(s.T(), err)
	assert.NotEmpty(s.T(), subcommits)
}

func (s *AnalyzeRepositoryTestSuite) TestExistingRepoSubcommitsAreAddedToExistingOnes() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL, adapters.ValidAccessToken})
	subcommitsBefore, _ := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidRepoID)

	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL, adapters.ValidAccessToken})
	subcommitsAfter, _ := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidRepoID)

	assert.Greater(s.T(), len(subcommitsAfter), len(subcommitsBefore))
}

func (s *AnalyzeRepositoryTestSuite) TestInvalidURLDoesNotStoreRepo() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.InvalidRepoURL, adapters.ValidAccessToken})
	_, err := s.repoRepository.GetRepo(context.Background(), adapters.InvalidRepoURL)

	assert.True(s.T(), errors.Is(err, repo.ErrRepositoryNotFound))
}

func (s *AnalyzeRepositoryTestSuite) TestInvalidURLDoesNotStoreSubcommits() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.InvalidRepoURL, adapters.ValidAccessToken})
	subcommits, err := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidRepoID)

	assert.Nil(s.T(), err)
	assert.Empty(s.T(), subcommits)
}

func (s *AnalyzeRepositoryTestSuite) TestAnalyzingTwoReposDoesNotMixSubcommits() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL, adapters.ValidAccessToken})
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidEmptyRepoURL, adapters.ValidAccessToken})
	subcommits, _ := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidEmptyRepoID)

	assert.Empty(s.T(), subcommits)
}

func (s *AnalyzeRepositoryTestSuite) TestSubcommitsBelongToAnalyzedRepo() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL, adapters.ValidAccessToken})
	subcommits, _ := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidRepoID)

	for _, sc := range subcommits {
		assert.Equal(s.T(), adapters.ValidRepoID, sc.RepoID())
	}
}

func (s *AnalyzeRepositoryTestSuite) TestExistingRepoIsNotDuplicatedAfterReAnalysis() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL, adapters.ValidAccessToken})
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL, adapters.ValidAccessToken})
	_, err := s.repoRepository.GetRepo(context.Background(), adapters.ValidRepoURL)

	assert.Nil(s.T(), err)
}

func (s *AnalyzeRepositoryTestSuite) TestEmptyRepoIsStillStored() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidEmptyRepoURL, adapters.ValidAccessToken})
	_, err := s.repoRepository.GetRepo(context.Background(), adapters.ValidEmptyRepoURL)

	assert.Nil(s.T(), err)
}

func (s *AnalyzeRepositoryTestSuite) TestEachCommitProducesAtLeastOneSubcommit() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.ValidRepoURL, adapters.ValidAccessToken})
	subcommits, _ := s.subcommitRepository.GetSubcommits(context.Background(), adapters.ValidRepoID)

	assert.GreaterOrEqual(s.T(), len(subcommits), 1)
}

func (s *AnalyzeRepositoryTestSuite) TestAgentFailureReturnsError() {
	err := s.handler.Handle(context.Background(), AnalyzeRepo{adapters.FailingAgentRepoURL, adapters.ValidAccessToken})
	assert.True(s.T(), errors.Is(err, agent.ErrAnalysisFailed))
}

func (s *AnalyzeRepositoryTestSuite) TestAgentFailureDoesNotStoreRepo() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.FailingAgentRepoURL, adapters.ValidAccessToken})
	_, err := s.repoRepository.GetRepo(context.Background(), adapters.FailingAgentRepoURL)

	assert.True(s.T(), errors.Is(err, repo.ErrRepositoryNotFound))
}

func (s *AnalyzeRepositoryTestSuite) TestAgentFailureDoesNotStoreSubcommits() {
	_ = s.handler.Handle(context.Background(), AnalyzeRepo{adapters.FailingAgentRepoURL, adapters.ValidAccessToken})
	subcommits, err := s.subcommitRepository.GetSubcommits(context.Background(), adapters.FailingAgentRepoID)

	assert.Nil(s.T(), err)
	assert.Empty(s.T(), subcommits)
}
