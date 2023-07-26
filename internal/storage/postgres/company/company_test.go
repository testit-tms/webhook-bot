package company

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/testit-tms/webhook-bot/internal/entities"
	"github.com/testit-tms/webhook-bot/internal/storage"
	"github.com/testit-tms/webhook-bot/pkg/database"
)

func TestCompanyStorage_AddCompany(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()
		expectedCompany := entities.Company{
			Id:      12,
			Token:   "bguFFFTF&ffdR9*9u",
			OwnerId: 21,
			Name:    "MyCompany",
			Email:   "info@google.com",
		}
		rows := sqlmock.NewRows([]string{"id", "token", "owner_id", "name", "email"}).
			AddRow(12, "bguFFFTF&ffdR9*9u", 21, "MyCompany", "info@google.com")

		f.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO companies (token, owner_id, name, email) VALUES ($1, $2, $3, $4) RETURNING id, token, owner_id, name, email")).
			WithArgs(expectedCompany.Token, expectedCompany.OwnerId, expectedCompany.Name, expectedCompany.Email).
			WillReturnRows(rows)

		repo := New(f.DB)

		// Act
		chat, err := repo.AddCompany(context.Background(), expectedCompany)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCompany, chat)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")
		expectedCompany := entities.Company{
			Id:      12,
			Token:   "bguFFFTF&ffdR9*9u",
			OwnerId: 21,
			Name:    "MyCompany",
			Email:   "info@google.com",
		}

		f.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO companies (token, owner_id, name, email) VALUES ($1, $2, $3, $4) RETURNING id, token, owner_id, name, email")).
			WithArgs(expectedCompany.Token, expectedCompany.OwnerId, expectedCompany.Name, expectedCompany.Email).
			WillReturnError(expectErr)

		repo := New(f.DB)

		// Act
		chat, err := repo.AddCompany(context.Background(), expectedCompany)

		// Assert
		assert.ErrorIs(t, err, expectErr)
		assert.Equal(t, entities.Company{}, chat)
	})
}

func TestCompanyStorage_GetCompaniesByOwnerId(t *testing.T) {
	t.Run("with companies", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		var id int64 = 21
		companiesExp := []entities.Company{
			{
				Id:      12,
				OwnerId: id,
				Name:    "MyCompany",
				Email:   "info@ya.ru",
				Token:   "23r32r23",
			},
			{
				Id:      13,
				OwnerId: id,
				Name:    "AnyCompany",
				Email:   "info@ya.ru",
				Token:   "rwe23t23t",
			},
		}

		rows := sqlmock.NewRows([]string{"id", "token", "owner_id", "name", "email"}).
			AddRow("12", "23r32r23", "21", "MyCompany", "info@ya.ru").
			AddRow("13", "rwe23t23t", "21", "AnyCompany", "info@ya.ru")

		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT id, token, owner_id, name, email FROM companies WHERE owner_id = $1")).
			WithArgs(id).
			WillReturnRows(rows)
		repo := New(f.DB)

		// Act
		companies, err := repo.GetCompaniesByOwnerId(context.Background(), id)

		// Assert
		assert.NoError(t, err)
		assert.ElementsMatch(t, companiesExp, companies)
	})

	t.Run("without companies", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		var id int64 = 21
		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT id, token, owner_id, name, email FROM companies WHERE owner_id = $1")).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)
		repo := New(f.DB)

		// Act
		companies, err := repo.GetCompaniesByOwnerId(context.Background(), id)

		// Assert
		assert.ErrorIs(t, err, storage.ErrNotFound)
		assert.Equal(t, []entities.Company{}, companies)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")

		var id int64 = 21
		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT id, token, owner_id, name, email FROM companies WHERE owner_id = $1")).
			WithArgs(id).
			WillReturnError(expectErr)
		repo := New(f.DB)

		// Act
		companies, err := repo.GetCompaniesByOwnerId(context.Background(), id)

		// Assert
		assert.Error(t, expectErr, err)
		assert.Equal(t, []entities.Company{}, companies)
	})
}
