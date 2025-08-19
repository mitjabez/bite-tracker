package repository

import (
	"context"
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/mitjabez/bite-tracker/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepoTestSuite struct {
	suite.Suite
	pgContext *testhelpers.PostgresContext
	repo      *UserRepo
	ctx       context.Context
}

func (suite *UserRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContext, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}
	suite.pgContext = pgContext
	suite.repo = NewUserRepo(&suite.pgContext.DBContext)
}

func (suite *UserRepoTestSuite) TeardownSuite() {
	if err := suite.pgContext.PostgresContainer.Terminate(suite.ctx); err != nil {
		log.Fatal("Error terminating container: ", err)
	}
}

func (suite *UserRepoTestSuite) TestUserExists() {
	t := suite.T()

	isUser, err := suite.repo.UserExists(suite.ctx, "sj@dot.com")
	assert.NoError(t, err)
	assert.True(t, isUser)
}

func (suite *UserRepoTestSuite) TestUserNotExists() {
	t := suite.T()

	isUser, err := suite.repo.UserExists(suite.ctx, "noname@dot.com")
	assert.NoError(t, err)
	assert.False(t, isUser)
}

func (suite *UserRepoTestSuite) TestCreateUser() {
	t := suite.T()

	fullName := "New User"
	email := "new@email.com"
	passwordHash := "123"
	user, err := suite.repo.CreateUser(suite.ctx, fullName, email, passwordHash)
	assert.NoError(t, err)
	assert.Equal(t, fullName, user.FullName)
	assert.Equal(t, email, user.Email)
	assert.Equal(t, passwordHash, user.PasswordHash)

	dbUser, err := suite.repo.GetUser(suite.ctx, user.Id)
	assert.NoError(t, err)
	assert.Equal(t, user, dbUser)
}

func (suite *UserRepoTestSuite) TestGetUserByEmail() {
	t := suite.T()

	user, err := suite.repo.GetUserByEmail(suite.ctx, "sj@dot.com")
	assert.NoError(t, err)
	assert.Equal(t, uuid.MustParse("f41ad27a-881d-4f7f-a908-f16a26ce7b78"), user.Id)
	assert.Equal(t, "sj@dot.com", user.Email)
	assert.Equal(t, "Salsa Jimmy", user.FullName)
	assert.Equal(t, "$2a$12$F22j/9fE8wI2nfjFADc/reQgm/TpKAxUWIyPhzZybV3GuvZP49rtu", user.PasswordHash)
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}
