# HW15-05…09: UseCase, cookie и composition root

Расширьте Login(ctx,LoginInput)(LoginResult,error), реализуйте
Refresh(ctx,string)(LoginResult,error) и Logout(ctx,string)error.
UseCase уже принимает session repository, generator и SessionConfig; в main
пока стоят nil dependencies. Замените их конкретными реализациями после решения.
Старые Register/password rules и GET /me не переписывайте.

Handler.AuthUseCase и public routes уже объявлены. Реализуйте Login/Refresh/Logout
HTTP responsibilities и helpers cookie в internal/httpapi/handler, явный DTO в dto.
Пример: Login возвращает access JSON и HttpOnly cookie, но ни одного поля raw
refresh в JSON. Refresh не требует действующего access. Cookie-only /me — 401.

Настройте TTL, имя cookie, Secure и env в main. Default Secure=true; local HTTP
может явно использовать false. Не помещайте http.Cookie в user и pgx.Tx в UseCase.
Подсказка: SessionRepository и Handler защищают разные границы системы.
Проверка: make test-acceptance; точные statuses/JSON в README.
Дальше: [tests](05_tests.md).
