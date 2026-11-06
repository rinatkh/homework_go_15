# HW15-03: два разных вида hash

Реализуйте CryptoRefreshTokenGenerator.Generate() (string,error) и Hash(string) string.
Generate возвращает opaque base64url credential из минимум 32 случайных байт.
Hash возвращает SHA-256 в hex. PasswordHasher остаётся bcrypt.

Пример свойства: Hash одного credential совпадает при повторном вызове, а два
Generate возвращают разные значения; raw credential и digest различаются.
Подсказка: выбирайте источник случайности для секретов, не генератор тестовых данных.
Время Hash — O(m) по длине credential, дополнительная память digest — O(1).
Проверка: go test -tags=acceptance ./internal/auth. Статистика не доказывает
криптографическую стойкость; требуется ревью выбора API.
Дальше: [repository](03_repository.md).
