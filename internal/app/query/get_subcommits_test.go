package query

import (
	"context"
	"errors"
	"testing"

	"github.com/octokerbs/chronocode/internal/adapters"
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
	s.repoRepository = adapters.NewRepoRepository()
	s.subcommitRepository = adapters.NewSubcommitRepository()
	s.codeHostFactory = adapters.NewCodeHostFactory()
	s.handler = NewGetSubcommitsHandler(s.repoRepository, s.subcommitRepository, s.codeHostFactory)
}

func (s *GetSubcommitsTestSuite) TestCannotGetSubcommitsWithoutAccessToken() {
	_, err := s.handler.Handle(context.Background(), GetSubcommits{adapters.ValidRepoURL, ""})
	assert.NotNil(s.T(), err)
}

func (s *GetSubcommitsTestSuite) TestCannotGetSubcommitsForInaccessibleRepo() {
	_, err := s.handler.Handle(context.Background(), GetSubcommits{adapters.ForbiddenRepoURL, adapters.ValidAccessToken})
	assert.True(s.T(), errors.Is(err, codehost.ErrAccessDenied))
}

func (s *GetSubcommitsTestSuite) TestCannotGetSubcommitsForNonExistentRepo() {
	_, err := s.handler.Handle(context.Background(), GetSubcommits{adapters.ValidRepoURL, adapters.ValidAccessToken})
	assert.True(s.T(), errors.Is(err, repo.ErrRepositoryNotFound))
}

func (s *GetSubcommitsTestSuite) TestReturnsSubcommitsForExistingRepo() {
	_ = s.repoRepository.StoreRepo(context.Background(), repo.NewRepo(adapters.ValidRepoID, "chronocode", adapters.ValidRepoURL, "FFFFFF"))
	subcommits, err := s.handler.Handle(context.Background(), GetSubcommits{adapters.ValidRepoURL, adapters.ValidAccessToken})

	assert.Nil(s.T(), err)
	assert.Empty(s.T(), subcommits)
}
