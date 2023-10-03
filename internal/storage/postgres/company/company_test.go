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
			ID:      12,
			Token:   "bguFFFTF&ffdR9*9u",
			OwnerID: 21,
			Name:    "MyCompany",
			Email:   "info@google.com",
		}
		rows := sqlmock.NewRows([]string{"id", "token", "owner_id", "name", "email"}).
			AddRow(12, "bguFFFTF&ffdR9*9u", 21, "MyCompany", "info@google.com")

		f.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO companies (token, owner_id, name, email) VALUES ($1, $2, $3, $4) RETURNING id, token, owner_id, name, email")).
			WithArgs(expectedCompany.Token, expectedCompany.OwnerID, expectedCompany.Name, expectedCompany.Email).
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
			ID:      12,
			Token:   "bguFFFTF&ffdR9*9u",
			OwnerID: 21,
			Name:    "MyCompany",
			Email:   "info@google.com",
		}

		f.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO companies (token, owner_id, name, email) VALUES ($1, $2, $3, $4) RETURNING id, token, owner_id, name, email")).
			WithArgs(expectedCompany.Token, expectedCompany.OwnerID, expectedCompany.Name, expectedCompany.Email).
			WillReturnError(expectErr)

		repo := New(f.DB)

		// Act
		chat, err := repo.AddCompany(context.Background(), expectedCompany)

		// Assert
		assert.ErrorIs(t, err, expectErr)
		assert.Equal(t, entities.Company{}, chat)
	})
}

func TestCompanyStorage_GetCompanyByOwnerTelegramId(t *testing.T) {
	t.Run("with company", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		var id int64 = 21
		companyExp := entities.Company{
			ID:      12,
			OwnerID: 13,
			Token:   "bguFFFTF&ffdR9*9u",
			Name:    "MyCompany",
			Email:   "info@ya.ru",
		}

		rows := sqlmock.NewRows([]string{"id", "token", "owner_id", "name", "email"}).
			AddRow(12, "bguFFFTF&ffdR9*9u", 13, "MyCompany", "info@ya.ru")

		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT c.id, c.token, c.owner_id, c.name, c.email FROM companies AS c INNER JOIN owners As o ON o.id = c.owner_id WHERE o.telegram_id=$1")).
			WithArgs(id).
			WillReturnRows(rows)
		repo := New(f.DB)

		// Act
		company, err := repo.GetCompanyByOwnerTelegramId(context.Background(), id)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, companyExp, company)
	})

	t.Run("without company", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		var id int64 = 21
		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT c.id, c.token, c.owner_id, c.name, c.email FROM companies AS c INNER JOIN owners As o ON o.id = c.owner_id WHERE o.telegram_id=$1")).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)
		repo := New(f.DB)

		// Act
		company, err := repo.GetCompanyByOwnerTelegramId(context.Background(), id)

		// Assert
		assert.ErrorIs(t, err, storage.ErrNotFound)
		assert.Equal(t, entities.Company{}, company)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")

		var id int64 = 21
		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT c.id, c.token, c.owner_id, c.name, c.email FROM companies AS c INNER JOIN owners As o ON o.id = c.owner_id WHERE o.telegram_id=$1")).
			WithArgs(id).
			WillReturnError(expectErr)
		repo := New(f.DB)

		// Act
		company, err := repo.GetCompanyByOwnerTelegramId(context.Background(), id)

		// Assert
		assert.Error(t, expectErr, err)
		assert.Equal(t, entities.Company{}, company)
	})
}
