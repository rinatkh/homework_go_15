-- +goose Up
-- +goose StatementBegin

-- TODO HW15-01: замените защитную заглушку своей схемой refresh sessions.
-- Незавершённая миграция не должна регистрироваться как успешно применённая.
DO $migration$
BEGIN
    RAISE EXCEPTION 'TODO HW15-01: реализуйте migrations/002_refresh_sessions.sql';
END
$migration$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- TODO HW15-01: реализуйте обратное изменение после написания Up.

-- +goose StatementEnd
