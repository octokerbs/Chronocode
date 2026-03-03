package query

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/octokerbs/chronocode/internal/adapters/memory"
	"github.com/octokerbs/chronocode/internal/domain/codehost"
	"github.com/octokerbs/chronocode/internal/domain/repo"
	"github.com/octokerbs/chronocode/internal/domain/subcommit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GetSubcommitsTestSuite struct {
	suite.Suite
	repoRepository      repo.Repository
	subcommitRepository subcommit.Repository
	codeHostFactory     codehost.CodeHostFactory
	handler             GetSubcommitsHandler
}

func TestGetSubcommitsTestSuite(t *testing.T) {
	suite.Run(t, new(GetSubcommitsTestSuite))
}

func (s *GetSubcommitsTestSuite) SetupTest() {
	s.repoRepository = memory.NewRepoRepository()
	s.subcommitRepository = memory.NewSubcommitRepository()
	s.codeHostFactory = memory.NewCodeHostFactory()
	s.handler = NewGetSubcommitsHandler(s.repoRepository, s.subcommitRepository, s.codeHostFactory)
}

func (s *GetSubcommitsTestSuite) TestCannotGetSubcommitsWithoutAccessToken() {
	_ = s.repoRepository.StoreRepo(context.Background(), repo.NewRepo(memory.ValidRepoID, "chronocode", memory.ValidRepoURL, "", time.Time{}))
	_, err := s.handler.Handle(context.Background(), GetSubcommits{memory.ValidRepoID, ""})
	assert.NotNil(s.T(), err)
}

func (s *GetSubcommitsTestSuite) TestCannotGetSubcommitsForInaccessibleRepo() {
	_ = s.repoRepository.StoreRepo(context.Background(), repo.NewRepo(memory.ForbiddenRepoID, "forbidden", memory.ForbiddenRepoURL, "", time.Time{}))
	_, err := s.handler.Handle(context.Background(), GetSubcommits{memory.ForbiddenRepoID, memory.ValidAccessToken})
	assert.True(s.T(), errors.Is(err, codehost.ErrAccessDenied))
}

func (s *GetSubcommitsTestSuite) TestCannotGetSubcommitsForNonExistentRepo() {
	_, err := s.handler.Handle(context.Background(), GetSubcommits{memory.ValidRepoID, memory.ValidAccessToken})
	assert.True(s.T(), errors.Is(err, repo.ErrRepositoryNotFound))
}

func (s *GetSubcommitsTestSuite) TestReturnsSubcommitsForExistingRepo() {
	_ = s.repoRepository.StoreRepo(context.Background(), repo.NewRepo(memory.ValidRepoID, "chronocode", memory.ValidRepoURL, "FFFFFF", time.Time{}))
	result, err := s.handler.Handle(context.Background(), GetSubcommits{memory.ValidRepoID, memory.ValidAccessToken})

	assert.Nil(s.T(), err)
	assert.Empty(s.T(), result.Subcommits)
	assert.Equal(s.T(), memory.ValidRepoURL, result.RepoURL)
}
