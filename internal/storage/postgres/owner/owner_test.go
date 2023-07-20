package owner

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/testit-tms/webhook-bot/internal/storage"
	"github.com/testit-tms/webhook-bot/pkg/database"
)

func TestOwnerStorage_GetOwnerById(t *testing.T) {
	t.Run("with owner", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		id := 12
		ownerExp := storage.Owner{
			Id:           id,
			TelegramId:   "123456",
			TelegramName: "Mega Owner",
		}

		rows := sqlmock.NewRows([]string{"id", "telegram_id", "telegram_name"}).
			AddRow("12", "123456", "Mega Owner")

		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT id, telegram_id, telegram_name FROM owners WHERE id=$1")).
			WithArgs(id).
			WillReturnRows(rows)
		repo := New(f.DB)

		// Act
		owner, err := repo.GetOwnerById(context.Background(), id)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, ownerExp, owner)
	})

	t.Run("without owner", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		id := 12
		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT id, telegram_id, telegram_name FROM owners WHERE id=$1")).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)
		repo := New(f.DB)

		// Act
		owner, err := repo.GetOwnerById(context.Background(), id)

		// Assert
		assert.ErrorIs(t, err, storage.ErrNotFound)
		assert.Equal(t, storage.Owner{}, owner)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")

		id := 12
		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT id, telegram_id, telegram_name FROM owners WHERE id=$1")).
			WithArgs(id).
			WillReturnError(expectErr)
		repo := New(f.DB)

		// Act
		owner, err := repo.GetOwnerById(context.Background(), id)

		// Assert
		assert.ErrorIs(t, err, expectErr)
		assert.Equal(t, storage.Owner{}, owner)
	})
}

func TestOwnerStorage_AddOwner(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()
		expectedOwner := storage.Owner{
			Id:           12,
			TelegramId:   "123456",
			TelegramName: "MyName",
		}
		rows := sqlmock.NewRows([]string{"id", "telegram_id", "telegram_name"}).
			AddRow(12, "123456", "MyName")

		f.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO owners (telegram_id, telegram_name) VALUES ($1, $2) RETURNING id, telegram_id, telegram_name")).
			WithArgs(expectedOwner.TelegramId, expectedOwner.TelegramName).
			WillReturnRows(rows)

		repo := New(f.DB)

		// Act
		owner, err := repo.AddOwner(context.Background(), expectedOwner)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedOwner, owner)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")
		expectedOwner := storage.Owner{
			Id:           12,
			TelegramId:   "123456",
			TelegramName: "MyName",
		}

		f.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO owners (telegram_id, telegram_name) VALUES ($1, $2) RETURNING id, telegram_id, telegram_name")).
			WithArgs(expectedOwner.TelegramId, expectedOwner.TelegramName).
			WillReturnError(expectErr)

		repo := New(f.DB)

		// Act
		owner, err := repo.AddOwner(context.Background(), expectedOwner)

		// Assert
		assert.ErrorIs(t, err, expectErr)
		assert.Equal(t, storage.Owner{}, owner)
	})
}

func TestOwnerStorage_DeleteOwnerById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()
		id := 12

		f.Mock.ExpectExec(regexp.QuoteMeta("DELETE FROM owners WHERE id=$1")).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := New(f.DB)

		// Act
		err := repo.DeleteOwnerById(context.Background(), id)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")
		id := 12

		f.Mock.ExpectExec(regexp.QuoteMeta("DELETE FROM owners WHERE id=$1")).
			WithArgs(id).
			WillReturnError(expectErr)

		repo := New(f.DB)

		// Act
		err := repo.DeleteOwnerById(context.Background(), id)

		// Assert
		assert.ErrorIs(t, err, expectErr)
	})
}
