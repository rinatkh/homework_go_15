-- +goose Up
-- +goose StatementBegin

CREATE TABLE users (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email text NOT NULL CONSTRAINT users_email_key UNIQUE,
    password_hash text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);
COMMENT ON TABLE users IS 'Пользователи самостоятельного Sessions API';
COMMENT ON COLUMN users.password_hash IS 'Только bcrypt hash, без исходного пароля';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS users;

-- +goose StatementEnd
