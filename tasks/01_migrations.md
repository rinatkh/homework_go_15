# HW15-01/02: схема и application contract

Реализуйте 002_refresh_sessions.sql в одном goose-файле с Up/Down и
StatementBegin/StatementEnd. 001_users.sql уже готов, повторно users не создавайте.
RefreshSession и RefreshSessionRepository объявлены в internal/user/refresh_session.go.
Сохраните предоставленные сигнатуры.

Хранятся id, user_id, token_hash, expires_at, created_at, revoked_at и связь
replaced_by_session_id. Уникальность hash, FK users и осознанный индекс обязательны.
Добавьте комментарии к таблице и нетривиальным колонкам. Raw token не хранится.
Один Down отменяет только вашу 002. Пример результата: status показывает 1 и 2
applied, после одного Down — 1 applied и 2 pending.

Подсказка: ограничения отвечают за корректность строк даже вне Go-приложения.
Проверка: make up / make status / make down / make up и make test-db.
Дальше: [tokens](02_tokens.md).
