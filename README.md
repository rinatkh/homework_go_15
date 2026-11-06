# Homework Go 15: refresh sessions для Sessions API

Самостоятельный Go 1.22 starter. Сохранены User, PostgreSQL UserRepository,
Register/Login, bcrypt, JWT и защищённый `GET /me`. Задание переносит механизм
classroom Task API на меньший Sessions API. Зависимостей от classroom нет.

## Что готово и что предстоит

Готовая регистрация нормализует email и хеширует пароль; login пока выдаёт
только access. Новые ports, сигнатуры и поля результата предоставлены как
контракт. Generate/Hash, session repository, Refresh/Logout, cookie и расширение
Login содержат TODO. Места реализации явно отмечены HW15-01…10. Реализаций
refresh/rotation в starter нет.

| Задание | Файлы | Условие |
| --- | --- | --- |
| 01 | migrations/002_refresh_sessions.sql | Своя goose Up/Down схема |
| 02 | internal/user/refresh_session.go | Разобрать модель/port, сохранить слои |
| 03 | internal/auth/refresh.go | Криптографический token и SHA-256 |
| 04 | internal/repository/postgres/session.go | Context, SQL и атомарная rotation |
| 05 | internal/user/usecase.go, refresh_session.go | Login сохраняет session; Refresh/Logout |
| 06 | internal/httpapi/handler/auth.go, cookie.go | Cookie, no-store, error mapping |
| 07 | internal/httpapi/dto/auth.go | Явный DTO без raw refresh |
| 08 | internal/httpapi/router.go | Public auth routes уже подключены к TODO handlers |
| 09 | cmd/api/main.go, .env.example | Заменить nil dependencies, config и composition root |
| 10 | *_acceptance_test.go, *_integration_test.go | Добиться зелёных проверок и дополнить сценарии |

Подробные условия, ограничения и подсказки находятся в [tasks/](tasks/01_migrations.md).
Сохраняйте прежние user/profile/middleware tests; не подменяйте ошибки заглушек
успешными ответами и не ослабляйте acceptance assertions.

## Контракт результата

Public: `GET /health`, `POST /register`, `POST /login`, `POST /auth/refresh`,
`POST /auth/logout`. Protected: `GET /me` только с Bearer access JWT.
Login и refresh: 200 `{"access_token":"...","token_type":"Bearer","expires_in":900}`
и refresh cookie. Register: 201, duplicate email 409; invalid credentials 401.
Неизвестный/истёкший/отозванный refresh: одна нейтральная 401. Logout: 204,
отсутствующий credential тоже допустим. Неожиданная DB error: 500 без деталей.

Cookie: настраиваемые Name и Secure, host-only, HttpOnly, SameSite=Lax, Path=/auth,
Expires и MaxAge. Default Secure=true, false только для local HTTP. При logout
совпадают Name/Path/атрибуты, MaxAge=-1 и прошлый Expires. Ответы с credentials —
`Cache-Control: no-store`. HttpOnly не заменяет CSRF-защиту; SameSite и CORS
относятся к разным ограничениям. Raw refresh, passwords и JWT не логируются.

## PostgreSQL и goose

Нужны PostgreSQL 16, psql, Go 1.22, Make и C compiler для race. Используйте
отдельную БД homework, не classroom. Администратор создаёт, например,
`go_homework_sessions` и `go_homework_sessions_test` с собственной ролью-владельцем.

```bash
export DATABASE_URL='postgres://REPLACE_USER:REPLACE_PASSWORD@localhost:5432/go_homework_sessions?sslmode=disable'
export TEST_DATABASE_URL='postgres://REPLACE_USER:REPLACE_PASSWORD@localhost:5432/go_homework_sessions_test?sslmode=disable'
export JWT_SECRET="$(python3 -c 'import secrets; print(secrets.token_hex(32))')"
export ACCESS_TOKEN_TTL=15m
export HTTP_ADDR=:8081
make status
make up
```

`001_users.sql` — готовая база starter. `002_refresh_sessions.sql` содержит
явную защитную TODO-заглушку: **до решения 01 `make up` ожидаемо остановится на
002**, оставив успешно применённую 001. Незавершённая схема не будет записана
как успешная версия 2. Удалите заглушку при написании своей Up/Down, затем
повторите `make up`. Ручные BEGIN/COMMIT не нужны: транзакцией управляет goose.
Файл `.env` автоматически не читается. После 001 уже можно запустить базовый
Login и `/me`; refresh endpoints пока отвечают 501 с названием TODO.

`make status` показывает историю. `make down` откатывает одну применённую
миграцию. После реализации 002 один Down удаляет sessions, ещё один — users.
Откат удаляет данные. Повторный Up не восстанавливает удалённые строки.

## Проверки: starter и готовое решение

```bash
make starter-check
make ci
make run
```

Эти команды проверяют сборку всех тестов, fmt/vet, готовые Register/Login/JWT/
profile tests и race. Они должны проходить на starter. Отдельные проверки
будущего результата **намеренно красные до реализации TODO**:

```bash
make test-acceptance
make test-db
make test-db-race
```

Acceptance tests проверяют генератор/hash, сохранение сессии в Login, refresh,
logout, cookie/DTO/no-store и neutral errors. SQL tests используют настоящую БД
и изолированные схемы, проверяют lifecycle, FK/cascade, rollback и две
конкурирующие rotation. Без TEST_DATABASE_URL они сообщают ошибку, не Skip.
До решения миграции SQL tests останавливаются на понятном HW15-01 TODO.
CI имеет зелёную starter-job и требуемую красную acceptance-job до выполнения
задания. Не отключайте её для получения формально зелёного pipeline.

## Acceptance и самостоятельная проверка

После решения пройдите register → login → /me → access expiration 401 → refresh
→ новый /me → replay старого refresh 401 → logout 204 → refresh 401. Перезапуск
с той же БД и JWT_SECRET сохраняет session. Access после logout может работать
до exp: отзыв refresh не превращает stateless JWT в мгновенно отозванный token.
Два одновременных refresh дают ровно один success; ошибка INSERT не отзывает
старую сессию частично. Unit tests не требуют БД. Добавьте свои expiry/context/
500 и cookie-clearing cases по [условиям](tasks/05_tests.md).

Bonus: logout-all. Redis, OAuth/OIDC, MFA и access denylist не требуются.
