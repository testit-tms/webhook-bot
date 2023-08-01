package chat

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

func TestChatStorage_GetChatsByCompanyId(t *testing.T) {
	t.Run("with chats", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		var id int64 = 21
		chatsExp := []entities.Chat{
			{
				Id:         12,
				CompanyId:  id,
				TelegramId: 123456,
			},
			{
				Id:         13,
				CompanyId:  id,
				TelegramId: 654321,
			},
		}

		rows := sqlmock.NewRows([]string{"id", "company_id", "telegram_id"}).
			AddRow("12", "21", "123456").
			AddRow("13", "21", "654321")

		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT id, company_id, telegram_id FROM chats WHERE company_id=$1")).
			WithArgs(id).
			WillReturnRows(rows)
		repo := New(f.DB)

		// Act
		chats, err := repo.GetChatsByCompanyId(context.Background(), id)

		// Assert
		assert.NoError(t, err)
		assert.ElementsMatch(t, chatsExp, chats)
	})

	t.Run("without chats", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		var id int64 = 21
		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT id, company_id, telegram_id FROM chats WHERE company_id=$1")).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)
		repo := New(f.DB)

		// Act
		chats, err := repo.GetChatsByCompanyId(context.Background(), id)

		// Assert
		assert.ErrorIs(t, err, storage.ErrNotFound)
		assert.Equal(t, []entities.Chat{}, chats)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")

		var id int64 = 21
		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT id, company_id, telegram_id FROM chats WHERE company_id=$1")).
			WithArgs(id).
			WillReturnError(expectErr)
		repo := New(f.DB)

		// Act
		chats, err := repo.GetChatsByCompanyId(context.Background(), id)

		// Assert
		assert.Error(t, expectErr, err)
		assert.Equal(t, []entities.Chat{}, chats)
	})
}

func TestChatStorage_GetChatsByCompanyToken(t *testing.T) {
	t.Run("with chats", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		token := "123"
		chatsExp := []entities.Chat{
			{
				Id:         12,
				CompanyId:  21,
				TelegramId: 123456,
			},
			{
				Id:         13,
				CompanyId:  21,
				TelegramId: 654321,
			},
		}

		rows := sqlmock.NewRows([]string{"id", "company_id", "telegram_id"}).
			AddRow("12", "21", "123456").
			AddRow("13", "21", "654321")

		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT id, company_id, telegram_id FROM chats WHERE company_id=(SELECT id FROM companies WHERE token=$1)")).
			WithArgs(token).
			WillReturnRows(rows)
		repo := New(f.DB)

		// Act
		chats, err := repo.GetChatsByCompanyToken(context.Background(), token)

		// Assert
		assert.NoError(t, err)
		assert.ElementsMatch(t, chatsExp, chats)
	})

	t.Run("without chats", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		token := "123"
		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT id, company_id, telegram_id FROM chats WHERE company_id=(SELECT id FROM companies WHERE token=$1)")).
			WithArgs(token).
			WillReturnError(sql.ErrNoRows)
		repo := New(f.DB)

		// Act
		chats, err := repo.GetChatsByCompanyToken(context.Background(), token)

		// Assert
		assert.ErrorIs(t, err, storage.ErrNotFound)
		assert.Equal(t, []entities.Chat{}, chats)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")

		token := "123"
		f.Mock.ExpectQuery(regexp.QuoteMeta("SELECT id, company_id, telegram_id FROM chats WHERE company_id=(SELECT id FROM companies WHERE token=$1)")).
			WithArgs(token).
			WillReturnError(expectErr)
		repo := New(f.DB)

		// Act
		chats, err := repo.GetChatsByCompanyToken(context.Background(), token)

		// Assert
		assert.Error(t, expectErr, err)
		assert.Equal(t, []entities.Chat{}, chats)
	})
}

func TestChatStorage_AddChat(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()
		expectedChat := entities.Chat{
			Id:         12,
			CompanyId:  21,
			TelegramId: 123456,
		}
		rows := sqlmock.NewRows([]string{"id", "company_id", "telegram_id"}).
			AddRow(12, 21, "123456")

		f.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO chats (company_id, telegram_id) VALUES ($1, $2) RETURNING id, company_id, telegram_id")).
			WithArgs(expectedChat.CompanyId, expectedChat.TelegramId).
			WillReturnRows(rows)

		repo := New(f.DB)

		// Act
		chat, err := repo.AddChat(context.Background(), expectedChat)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedChat, chat)
	})

	t.Run("with error", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()

		expectErr := errors.New("test error")
		expectedChat := entities.Chat{
			Id:         12,
			CompanyId:  21,
			TelegramId: 123456,
		}

		f.Mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO chats (company_id, telegram_id) VALUES ($1, $2) RETURNING id, company_id, telegram_id")).
			WithArgs(expectedChat.CompanyId, expectedChat.TelegramId).
			WillReturnError(expectErr)

		repo := New(f.DB)

		// Act
		chat, err := repo.AddChat(context.Background(), expectedChat)

		// Assert
		assert.ErrorIs(t, err, expectErr)
		assert.Equal(t, entities.Chat{}, chat)
	})
}

func TestChatStorage_DeleteChatById(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()
		id := 12

		f.Mock.ExpectExec(regexp.QuoteMeta("DELETE FROM chats WHERE id=$1")).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := New(f.DB)

		// Act
		err := repo.DeleteChatById(context.Background(), id)

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

		f.Mock.ExpectExec(regexp.QuoteMeta("DELETE FROM chats WHERE id=$1")).
			WithArgs(id).
			WillReturnError(expectErr)

		repo := New(f.DB)

		// Act
		err := repo.DeleteChatById(context.Background(), id)

		// Assert
		assert.ErrorIs(t, err, expectErr)
	})
}

func TestChatStorage_DeleteChatByCompanyId(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		t.Parallel()
		f := database.NewFixture(t)
		defer f.Teardown()
		id := 12

		f.Mock.ExpectExec(regexp.QuoteMeta("DELETE FROM chats WHERE company_id=$1")).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := New(f.DB)

		// Act
		err := repo.DeleteChatByCompanyId(context.Background(), id)

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

		f.Mock.ExpectExec(regexp.QuoteMeta("DELETE FROM chats WHERE company_id=$1")).
			WithArgs(id).
			WillReturnError(expectErr)

		repo := New(f.DB)

		// Act
		err := repo.DeleteChatByCompanyId(context.Background(), id)

		// Assert
		assert.ErrorIs(t, err, expectErr)
	})
}
