# Spy Cat Agency Management System

## Опис

Ця система дозволяє керувати агентством шпигунських котів, включаючи додавання, видалення, оновлення та відображення котів, місій та цілей.

## Вимоги

- Go 1.16+
- PostgreSQL

## Інструкції по запуску

1. Клонуйте репозиторій.
2. Створіть файл `.env` і додайте ваші дані для підключення до бази даних PostgreSQL: - ? або замінити в 1 місці 2 змінні
3. Виконайте `go mod tidy` для встановлення залежностей.
4. Запустіть додаток за допомогою команди `go run main.go`.

## Ендпоінти

- `POST /cats` - створення кота.
- `GET /cats` - отримання списку котів.
- `GET /cats/:id` - отримання інформації про кота.
- `PUT /cats/:id` - оновлення зарплати кота.
- `DELETE /cats/:id` - видалення кота.
- `POST /missions` - створення місії.
- `GET /missions` - отримання списку місій.
- `GET /missions/:id` - отримання інформації про місію.
- `PUT /missions/:id` - оновлення статусу місії.
- `DELETE /missions/:id` - видалення місії.
- `POST /targets/:mission_id` - додавання цілі до місії.
- `PUT /targets/:id` - оновлення нотаток та статусу цілі.
- `DELETE /targets/:id` - видалення цілі.
- `PUT /complete_mission/:id` - завершення місії та її цілей.

## env.exemaple - file in project
