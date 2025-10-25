-- Создаем таблицу пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создаем индекс для быстрого поиска по email
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- -- Вставляем тестового пользователя (опционально)
-- INSERT INTO users (email, password) 
-- VALUES ('test@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi') -- password
-- ON CONFLICT (email) DO NOTHING;