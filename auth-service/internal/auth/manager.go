package auth

import "time"

// TokenData — данные, извлечённые из токена.
// Их можно будет класть в контекст запроса.
type TokenData struct {
	UserID int64
	Exp    time.Time
}

// AuthManager — интерфейс для управления токенами.
//
// Любая реализация (простое JWT, JWT+Redis, cookie-сессии)
// должна реализовывать эти методы.
//
// Такой подход позволяет быстро подменить механизм авторизации,
// не трогая хендлеры или бизнес-логику.
type AuthManager interface {
	// GenerateToken — создаёт токен для пользователя.
	// Возвращает строку токена (обычно JWT) и ошибку.
	GenerateToken(userID int64) (string, error)

	// ParseToken — проверяет токен, валидирует подпись и срок действия,
	// возвращает информацию о пользователе (UserID, Exp).
	ParseToken(token string) (*TokenData, error)

	// Invalidate — инвалидирует токен пользователя.
	// В JWT без сессий может быть no-op,
	// а в JWT+Redis — очищает версию в Redis.
	Invalidate(userID int64) error
}
