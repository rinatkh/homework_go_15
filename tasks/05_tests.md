# HW15-10: доказательство поведения

Сохраните базовые tests; добейтесь прохождения acceptance и PostgreSQL tests.
Файлы с build tag acceptance проверяют будущий результат и не являются готовыми
реализациями. SQL tests создают свою схему, применяют goose и удаляют только её.

Добавьте собственные cases: истёкшая session, отменённый ctx, ошибка storage
до выдачи cookie, корректное удаление cookie, invalid JWT algorithm, reconnect
после restart. Для concurrency важен ровно один success, а не отсутствие Go race.
Не используйте time.Sleep в unit tests expiry: контролируйте время или claims.

Составьте curl jar flow без вывода credentials. Проверьте 401/409/500 и два
параллельных refresh; докажите rollback при failed insert. Укажите HTTP statuses,
что сохраняется в БД и почему access после logout ещё может работать.

Подсказка: зелёная компиляция TODO подтверждает только корректность starter.
Команды: make ci / make test-acceptance / make test-db / make test-db-race.
