Бэк собирается через docker compose,в докер композ можно указать конфиги с помощью env

Фронтенд только через npm run dev

```typescriptreact project="library-management"
...
```

Это базовая структура React приложения "Рабочее место библиотекаря" с использованием Next.js. Вот краткое описание того, что мы создали:

1. Конфигурационный файл (`config.ts`) для хранения адреса бэкенда.
2. Хук `useApi` для работы с API.
3. Основной макет приложения с навигацией.
4. Главная страница (Dashboard) с общей статистикой.
5. Страница управления книгами с возможностью добавления, удаления и поиска книг.
6. Страница управления пользователями с возможностью регистрации, поиска и просмотра информации о выданных книгах.
7. Страница для выдачи и возврата книг.
8. Страница статистики с графиками популярных категорий и активных пользователей.


Для полноценной работы приложения вам потребуется реализовать бэкенд с соответствующими API-эндпоинтами. Также, вы можете дополнить это приложение дополнительными функциями, такими как:

- Аутентификация и авторизация пользователей
- Более детальная информация о книгах и пользователях
- Система уведомлений о просроченных книгах
- Расширенные возможности поиска и фильтрации


Не забудьте настроить переменные окружения для адреса бэкенда перед запуском приложения.




Конечно! Я с удовольствием помогу вам составить документацию для реализации API на бэкенде. Эта документация поможет вам создать совместимый бэкенд для нашего фронтенд приложения "Рабочее место библиотекаря".

# Документация API для "Рабочее место библиотекаря"

## Базовый URL

```plaintext
http://localhost:3001/api
```

## Эндпоинты

### 1. Книги

#### 1.1. Получить список всех книг

- **URL:** `/books`
- **Метод:** GET
- **Ответ:**

```json
[
  {
    "id": 1,
    "title": "Название книги",
    "author": "Автор книги",
    "category": "Категория книги"
  },
  ...
]
```




#### 1.2. Добавить новую книгу

- **URL:** `/books`
- **Метод:** POST
- **Тело запроса:**

```json
{
  "title": "Название книги",
  "author": "Автор книги",
  "category": "Категория книги"
}
```


- **Ответ:**

```json
{
  "id": 1,
  "title": "Название книги",
  "author": "Автор книги",
  "category": "Категория книги"
}
```




#### 1.3. Удалить книгу

- **URL:** `/books/:id`
- **Метод:** DELETE
- **Параметры URL:** `id` - ID книги
- **Ответ:**

```json
{
  "message": "Книга успешно удалена"
}
```




### 2. Пользователи

#### 2.1. Получить список всех пользователей

- **URL:** `/users`
- **Метод:** GET
- **Ответ:**

```json
[
  {
    "id": 1,
    "name": "Имя пользователя",
    "email": "email@example.com",
    "borrowedBooks": ["Название книги 1", "Название книги 2"]
  },
  ...
]
```




#### 2.2. Добавить нового пользователя

- **URL:** `/users`
- **Метод:** POST
- **Тело запроса:**

```json
{
  "name": "Имя пользователя",
  "email": "email@example.com"
}
```


- **Ответ:**

```json
{
  "id": 1,
  "name": "Имя пользователя",
  "email": "email@example.com",
  "borrowedBooks": []
}
```




#### 2.3. Удалить пользователя

- **URL:** `/users/:id`
- **Метод:** DELETE
- **Параметры URL:** `id` - ID пользователя
- **Ответ:**

```json
{
  "message": "Пользователь успешно удален"
}
```




### 3. Выдача и возврат книг

#### 3.1. Получить список всех выдач

- **URL:** `/loans`
- **Метод:** GET
- **Ответ:**

```json
[
  {
    "id": 1,
    "userId": 1,
    "bookId": 1,
    "borrowDate": "2023-06-01T10:00:00Z",
    "returnDate": null
  },
  ...
]
```




#### 3.2. Выдать книгу

- **URL:** `/loans`
- **Метод:** POST
- **Тело запроса:**

```json
{
  "userId": 1,
  "bookId": 1
}
```


- **Ответ:**

```json
{
  "id": 1,
  "userId": 1,
  "bookId": 1,
  "borrowDate": "2023-06-01T10:00:00Z",
  "returnDate": null
}
```




#### 3.3. Вернуть книгу

- **URL:** `/loans/:id/return`
- **Метод:** POST
- **Параметры URL:** `id` - ID выдачи
- **Ответ:**

```json
{
  "id": 1,
  "userId": 1,
  "bookId": 1,
  "borrowDate": "2023-06-01T10:00:00Z",
  "returnDate": "2023-06-15T14:30:00Z"
}
```




### 4. Статистика

#### 4.1. Получить общую статистику

- **URL:** `/statistics`
- **Метод:** GET
- **Ответ:**

```json
{
  "totalBooks": 1000,
  "totalUsers": 500,
  "totalLoans": 250,
  "popularCategories": [
    { "name": "Фантастика", "count": 150 },
    { "name": "Детективы", "count": 100 },
    ...
  ],
  "activeUsers": [
    { "name": "Иван Иванов", "loansCount": 10 },
    { "name": "Мария Петрова", "loansCount": 8 },
    ...
  ]
}
```




## Обработка ошибок

В случае ошибки, API должен возвращать соответствующий HTTP-статус и JSON-объект с описанием ошибки:

```json
{
  "error": "Описание ошибки"
}
```

## Пагинация

Для эндпоинтов, возвращающих списки (книги, пользователи, выдачи), рекомендуется реализовать пагинацию. Пример параметров запроса:

```plaintext
GET /api/books?page=1&limit=20
```

Ответ должен включать метаданные о пагинации:

```json
{
  "data": [...],
  "meta": {
    "currentPage": 1,
    "totalPages": 5,
    "totalItems": 100,
    "itemsPerPage": 20
  }
}
```

## Фильтрация и поиск

Для эндпоинтов, возвращающих списки, рекомендуется реализовать возможность фильтрации и поиска. Пример параметров запроса:

```plaintext
GET /api/books?search=фантастика&category=Научная+фантастика
```

## Аутентификация и авторизация

Рекомендуется реализовать систему аутентификации и авторизации для защиты API. Например, можно использовать JWT (JSON Web Tokens) для аутентификации запросов.

Эта документация предоставляет основу для реализации API вашего бэкенда. Вы можете расширить ее, добавив дополнительные эндпоинты или функциональность по мере необходимости. При реализации бэкенда убедитесь, что он соответствует этой спецификации для обеспечения совместимости с фронтендом.