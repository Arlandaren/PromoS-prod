# Документация к проекту PROMO v2: Promo Code Backend
Проект представляет собой решение задания на хакатон.

Условие задания и публичные тесты доступны [в данном репозитории](https://github.com/Central-University-IT/FAQ-2025/tree/main/backend).

## Обзор решения

Решение представляет собой HTTP-сервер на Go, реализующий API для управления промокодами с разделением на B2B (для компаний) и B2C (для пользователей) функционал. Сервер использует PostgreSQL для хранения данных и Redis для кеширования.

## Архитектура решения

Решение построено по слоистой архитектуре с четким разделением ответственностей:

1. **Transport layer** (HTTP handlers) - обработка запросов и ответов
2. **Service layer** - бизнес-логика приложения
3. **Repository layer** - работа с хранилищами данных
4. **Models/DTO** - структуры данных и валидация

## Основные компоненты

### Конфигурация
- Чтение переменных окружения
- Настройка подключения к БД и Redis
- Инициализация зависимостей через DI контейнер

### B2B функционал (для компаний)
- Регистрация и аутентификация компаний
- Создание и управление промокодами
- Получение статистики по промокодам

### B2C функционал (для пользователей)
- Регистрация и аутентификация пользователей
- Управление профилем пользователя
- Просмотр ленты промокодов
- Лайки и комментарии к промокодам
- Активация промокодов с проверкой через антифрод-сервис

## Запуск приложения

Приложение конфигурируется через переменные окружения:

```bash
SERVER_ADDRESS=0.0.0.0:8080
POSTGRES_CONN=postgres://user:password@localhost:5432/dbname
REDIS_HOST=localhost
REDIS_PORT=6379
ANTIFRAUD_ADDRESS=localhost:9090
RANDOM_SECRET=random128charsstring
```

Для сборки и запуска:
```bash
docker build -t promo-backend .
docker run -e SERVER_ADDRESS=0.0.0.0:8080 -p 8080:8080 promo-backend
```

## API Endpoints

### Общие
- `GET /api/ping` - проверка работоспособности сервера

### B2B Endpoints
- `POST /api/business/auth/sign-up` - регистрация компании
- `POST /api/business/auth/sign-in` - аутентификация компании
- `POST /api/business/promo` - создание промокода
- `GET /api/business/promo` - список промокодов компании
- `GET /api/business/promo/{id}` - получение промокода по ID
- `PATCH /api/business/promo/{id}` - обновление промокода
- `GET /api/business/promo/{id}/stat` - статистика по промокоду

### B2C Endpoints
- `POST /api/user/auth/sign-up` - регистрация пользователя
- `POST /api/user/auth/sign-in` - аутентификация пользователя
- `GET /api/user/profile` - получение профиля пользователя
- `PATCH /api/user/profile` - обновление профиля
- `GET /api/user/feed` - лента промокодов
- `GET /api/user/promo/{id}` - информация о промокоде
- `POST /api/user/promo/{id}/like` - поставить лайк промокоду
- `DELETE /api/user/promo/{id}/like` - убрать лайк
- `POST /api/user/promo/{id}/comments` - добавить комментарий
- `GET /api/user/promo/{id}/comments` - список комментариев
- `GET /api/user/promo/{id}/comments/{comment_id}` - получить комментарий
- `PUT /api/user/promo/{id}/comments/{comment_id}` - изменить комментарий
- `DELETE /api/user/promo/{id}/comments/{comment_id}` - удалить комментарий
- `POST /api/user/promo/{id}/activate` - активировать промокод
- `GET /api/user/promo/history` - история активаций промокодов

## Особенности реализации

1. **Антифрод-интеграция**:
   - При активации промокода проверяется пользователь через антифрод-сервис
   - Реализовано кеширование ответов антифрода согласно `cache_until`
   - Повторные запросы при ошибках (retry logic)

2. **Типы промокодов**:
   - COMMON - фиксированное значение, ограниченное количество активаций
   - UNIQUE - уникальные значения из списка, выдается по одному

3. **Безопасность**:
   - Хеширование паролей (bcrypt)
   - JWT токены для аутентификации
   - Проверка прав доступа к ресурсам

4. **Производительность**:
   - Кеширование в Redis
   - Пагинация и фильтрация на уровне БД
   - Оптимизированные запросы

5. **Надежность**:
   - Обработка ошибок
   - Валидация входных данных
   - Транзакции для критичных операций

## Технологический стек

- Язык: Go 1.21+
- Фреймворк: Gin
- ORM: GORM
- Базы данных: PostgreSQL, Redis
- DI контейнер: собственный простой DI
- Документация API: OpenAPI 3.0 (Swagger)

## Локальное тестирование

Для локального тестирования можно использовать docker-compose:

```bash
docker-compose up -d
```

Запуск тестов:
```bash
cd tests
python3 -m pip install -r requirements.txt
export BASE_URL="http://localhost:8080/api"
export ANTIFRAUD_URL="http://localhost:9090/internal"
py.test test_01_ping.tavern.yml
```

made with ❤️❤️❤️ by @Arlandaren