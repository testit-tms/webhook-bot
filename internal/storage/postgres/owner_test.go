package postgres

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/testit-tms/webhook-bot/internal/storage"
)

func TestOwnerStorage_GetOwnerById(t *testing.T) {
	t.Run("with owner", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()

		id := 12
		ownerExp := storage.Owner{
			Id:           id,
			TelegramId:   "123456",
			TelegramName: "Mega Owner",
		}

		rows := sqlmock.NewRows([]string{"id", "telegram_id", "telegram_name"}).
			AddRow("12", "123456", "Mega Owner")

		f.mock.ExpectQuery(regexp.QuoteMeta("SELECT id, telegram_id, telegram_name FROM owners WHERE id=$1")).
			WithArgs(id).
			WillReturnRows(rows)
		repo := New(f.db)

		// Act
		owner, err := repo.GetOwnerById(context.Background(), id)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, ownerExp, owner)
	})

	t.Run("without owner", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()

		id := 12
		f.mock.ExpectQuery("SELECT id, telegram_id, telegram_name FROM owners WHERE id=$1").
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)
		repo := New(f.db)

		// Act
		owner, err := repo.GetOwnerById(context.Background(), id)

		// Assert
		assert.Error(t, storage.ErrNotFound, err)
		assert.Equal(t, storage.Owner{}, owner)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		f := NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")

		id := 12
		f.mock.ExpectQuery("SELECT id, telegram_id, telegram_name FROM owners WHERE id=$1").
			WithArgs(id).
			WillReturnError(expectErr)
		repo := New(f.db)

		// Act
		owner, err := repo.GetOwnerById(context.Background(), id)

		// Assert
		assert.Error(t, expectErr, err)
		assert.Equal(t, storage.Owner{}, owner)
	})

}
