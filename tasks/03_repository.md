# HW15-04: PostgreSQL refresh repository

В internal/repository/postgres/session.go уже есть конструктор и заглушки
Create, GetByTokenHash, Rotate, Revoke. Сигнатуры соответствуют application port.
Реализуйте их с переданным ctx. Ошибка отсутствующей session — ErrInvalidRefreshToken,
повтор hash — ErrSessionAlreadyExists; revoked/expired имеют отдельные application errors.

Контракт Rotate: атомарная замена credential. После успешного вызова old
отозван, next активен; при любой ошибке old остаётся в прежнем состоянии.
Две параллельные попытки не могут успешно заменить один old дважды.
Revoke идемпотентен. Предусловие next: тот же UserID и ещё действующий ExpiresAt.

Подсказка: согласованность должна обеспечиваться и между разными процессами API.
Пример: ошибка unique при сохранении next не меняет active-состояние old.
Проверка: make test-db / make test-db-race. Оценка стоимости учитывает индекс
по hash и сетевой round trip; это не константная операция памяти Go.
Дальше: [auth и HTTP](04_auth_http.md).
