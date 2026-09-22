  # Order API

  REST API для управления заказами, пользователями и товарами. Реализовал на Go с двухфакторной авторизацией по SMS-коду, JWT-токенами и связью many-to-many между заказами и товарами.

  ## Стек технологий 

  - **Язык:** Go 1.26
  - **Роутинг:** gorilla/mux
  - **База данных:** PostgreSQL 16.4, GORM
  - **Авторизация:** JWT (HS256)
  - **Валидация:** go-playground/validator
  - **Конфигурация:** godotenv
  - **Контейнеризация:** Docker, Docker Compose
  - **Тестирование:** sqlmock (unit), httptest (E2E)

  ## Архитектура

  Проект построен на слоистой архитектуре.

  - **Handler** - HTTP слой: принимает запросы, возвращает  ответы
  - **Service** - бизнес-логика
  - **Repository** - работа с базой данных
    
  ```
6-order-api-cart/
├── cmd/
│   ├── main.go             — точка входа
│   └── order_test.go       — E2E-тест
├── configs/
│   └── conf.go             — конфигурация
├── internal/
│   ├── auth/               — авторизация
│   ├── order/              — заказы
│   └── product/            — товары
├── migrations/
│   └── auto.go             — миграции
├── pkg/
│   ├── contextKeys/        — типизированные ключи контекста
│   ├── db/                 — подключение к БД
│   ├── jwt/                — работа с JWT
│   ├── middleware/         — middleware
│   ├── req/                — парсинг и валидация запросов
│   └── resp/               — ответы
├── .dockerignore           — исключения для Docker
├── .env                    — пример переменных окружения
├── .gitignore              — исключения для Git
├── docker-compose.yml      — оркестрация контейнеров
├── Dockerfile              — сборка образа
├── go.mod                  — зависимости Go
├── go.sum                  — контрольные суммы зависимостей
├── Makefile                — команды управления
└── README.md               — документация
  ```

  ## Запуск

  ### Требования 

  - Docker
  - Docker Compose
  - .env file
  ### Запуск через Docker 

  ```bash
  docker compose up -d
  ```

  API будет доступен через http://localhost:8081
  
  Миграции (создание таблиц) выполняются автоматически при старте.

  ### Остановка

  ```bash
  docker compose down
  ```

  ### Полная очистка (с удалением данных)

  ```bash
  docker compose down -v
  ```

  ## Эндпоинты

  ### Авторизация
  | Метод | Путь          | Описание                 |
  | ----- | ------------- | ------------------------ |
  | POST  | /authbyphone  | Отправить номер телефона |
  | POST  | /verifybycode | Получить токен           |

  ### Товары 
  | Метод | Путь         | Описание            |
  | ----- | ------------ | ------------------- |
  | POST  | /product     | Создать товар       |
  | GET   | /product/{id}| Получить товар по ID|
  | PATCH | /product/{id}| Обновить товар      |
  | DELETE| /product/{id}| Удалить товар       |

  ### Заказы
  | Метод | Путь      | Описание                             |
  | ----- | --------- | ------------------------------------ |      
  |POST   |/order     | Создать заказ (требует токен)        |
  |GET    |/order/{id}| Получить заказ по ID (требует токен) |
  |GET    |/my-orders | Получить свои заказы (требует токен) |

  ## Примеры запросов
  ### 1. Авторизация по телефону

  ```bash
  curl -X POST http://localhost:8081/authbyphone \
  -H "Content-Type: application/json" \
  -d '{"phone": "9281112233"}'
  ```

  Ответ:

  В dev-режиме код возвращается в ответе для удобства тестирования
    
  ```json
  {
    "sessionId": "abc123",
    "code": 1234
  }
  ```

  ### 2. Подтверждение кода

  ```bash
  curl -X POST http://localhost:8081/verifybycode \
  -H "Content-Type: application/json" \
  -d '{"sessionId": "abc123", "code": 1234}'
  ```

  Ответ:

  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
  ```

  ### 3. Создание товара

  ```bash
  curl -X POST http://localhost:8081/product \
  -H "Content-Type: application/json" \
  -d '{
        "name": "Apple",
        "description": "Fruit",
        "images": ["http://example.com/apple.jpg"]
  }'
  ```

### 3.1 Получение товара по ID

```bash
curl -X GET http://localhost:8081/product/1
```

### 3.2 Обновление товара

```bash
curl -X PATCH http://localhost:8081/product/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Green Apple",
    "description": "Fresh fruit",
    "images": [
    	"http://example.com/image1.jpg",
    	"http://example.com/image2.jpg"
    ]
  }'
```

### 3.3 Удаление товара

```bash
curl -X DELETE http://localhost:8081/product/1
```

### 4. Создание заказа

  ```bash
  curl -X POST http://localhost:8081/order \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"products_id": [1,2]}'
  ```

  Ответ:

  ```json
  {
    "id": 1,
    "userid": 1,
    "product": [
      {
        "id": 1,
        "CreatedAt": "2026-09-12T10:28:45.552653Z",
        "UpdatedAt": "2026-09-12T10:28:45.552653Z",
        "DeletedAt": null,
        "Name": "mango",
        "Description": "Fruit",
        "Images": [
          "http://example.com/image1.jpg",
          "http://example.com/image2.jpg"
        ]
      }
    ],
    "status": "pending"
  }
  ```

  ### 5. Получение заказа по ID

  ```bash
  curl -X GET http://localhost:8081/order/{id} \
  -H "Authorization: Bearer <token>"
  ```

  ### 6. Получение своих заказов

  ```bash
  curl -X GET http://localhost:8081/my-orders \
  -H "Authorization: Bearer <token>"
  ```

  ## Тестирование

  ### Unit-тесты

  ```bash
  go test -v ./internal/order/...
  ```

  ### E2E-тесты

  Требуют запущенной базы данных

  ```bash
  docker compose up -d postgres
  go test ./cmd/...
  ```

  ### Все тесты

  ```bash
  go test ./cmd/... ./internal/...
  ```

  ## Makefile

  | Команда        | Описание                  |
  | -------------- | ------------------------- |
  | make up        | Поднять все через Docker  |
  | make down      | Остановить контейнеры     |
  | make logs      | Все логи                  |
  | make logs-app  | Логи приложения           |
  | make build     | Собрать образ             |
  | make E2E-test  | Запустить E2E-тест        |
  | make Unit-test | Запустить UNIT-тесты      |
  | make clean     | Полная очистка с данными  |
  | make help      | Справка                   |

  ### Переменные окружения

  Создайте файл .env в корне проекта:

  ```text
  DSN=host=postgres user=postgres password=<your_pass> dbname=orders port=5432 sslmode=disable
  SECRET=<your-secret-key>
  ```

  ### Структура базы данных

- auth_by_phones - пользователи (телефон, сессия, код)
- products - товары (название, описание, изображения)
- orders — заказы (пользователь, статус)
- order_items - промежуточная таблица (заказ ↔ товар, many-to-many)
