# 📎 Short URL Service

Pet-проект для сокращения ссылок с авторизацией пользователей и статистикой переходов.

## 🚀 Функционал
- 🔐 Авторизация и регистрация пользователя  
- ✂ Создание короткой ссылки (**только для авторизованных пользователей**)  
- 📜 Получение всех своих ссылок (**только для авторизованных пользователей**)  
- 📊 Получение статистики переходов по ссылке  

## 🔧Возможные улучшения
- Добавление пагинации при запросе статистики
- Добавление пагинации при запросе своих ссылок
- Улучшить возрат ошибок

## 🏗️ Стек технологий
- **Go + net/http**
- **PostgreSQL + pgx**
- **Миграции: Goose**
- **Аунтификация: JWT**
- **Docker + Docker Compose**

---

## 📌 API

### 1. Регистрация пользователя
**POST** `http://localhost:8080/auth/register`<br>
**Content-Type**: `application/json`<br>
**Тело запроса**:
```json
{
    "email": "test@example.com",
    "password": "test123"
}
```
**Пример ответа**:
```json
{
    "id": "839839f7-0ea7-4608-8da3-9d4c379a8975",
    "email": "test@example.com"
}
```

### 2. Авторизация
**POST** `http://localhost:8080/auth/login`<br>
**Content-Type**: `application/json`<br>
**Тело запроса**:
```json
{
    "email": "test@example.com",
    "password": "test123"
}
```
**Пример ответа**:
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6..."
}
```

### 3. Создание ссылки
**POST** `http://localhost:8080/url/create` <br>
**Content-Type**: `application/json`<br>
**Авторизация**: `Authorization: Bearer <ВАШ_ТОКЕН>`<br>
**Тело запроса**:
```json
{
    "original_url": "http://url.com"
}
```
**Пример ответа**:
```json
{
    "id": "77b00cfb-1cea-4efd-aefc-bed3df32a0a3",
    "user_id": "839839f7-0ea7-4608-8da3-9d4c379a8975",
    "original_url": "http://url.com",
    "short_url": "http://localhost:8080/url/k0WOZ3",
    "short_code": "k0WOZ3",
    "created_at": "2025-08-13T11:37:53.948841Z"
}
```

### 4. Получение статистики по ссылке
**GET** `http://localhost:8080/url/stats/{shortCode}`<br>
**Авторизация**: `Authorization: Bearer <ВАШ_ТОКЕН>`<br>
**Пример ответа**:
```json
{
    "id": "b672fb1c-9a8e-4b22-8ad3-cdf1c39b3b3d",
    "original_url": "https://i.pinimg.com/originals/6c/ff/e4/6cffe4536d62041d3bd6fc3c1070c4be.jpg",
    "short_url": "http://localhost:8080/url/tTUmEQ",
    "short_code": "tTUmEQ",
    "created_at": "2025-08-01T14:39:24.189504Z",
    "is_active": true,
    "total_clicks": 2,
    "clicks": [
        {
            "id": "87609e71-94db-42ee-967b-dc7dc8fd0b2b",
            "link_id": "b672fb1c-9a8e-4b22-8ad3-cdf1c39b3b3d",
            "timestamp": "2025-08-01T14:39:32.972996Z",
            "ip_address": "...",
            "user_agent": "Mozilla/5.0 ...",
            "referrer": ""
        },
        {
            "id": "3dd4b837-f775-4419-b716-b45e2cf59c21",
            "link_id": "b672fb1c-9a8e-4b22-8ad3-cdf1c39b3b3d",
            "timestamp": "2025-08-01T14:39:36.954587Z",
            "ip_address": "...",
            "user_agent": "Mozilla/5.0 ...",
            "referrer": ""
        }
    ]
}
```

### 5. Получение всех своих ссылок
**GET** `http://localhost:8080/url/links`<br>
**Авторизация**: `Authorization: Bearer <ВАШ_ТОКЕН>`<br>
**Пример ответа**:
```json
[
    {
        "id": "c5bf837c-64b9-4ef7-9a44-aca176d4db01",
        "user_id": "44092648-bca4-48fd-84cc-adf4fc840ef8",
        "original_url": "https://translate.yandex.ru",
        "short_url": "http://localhost:8080/url/8X3Vme",
        "short_code": "8X3Vme",
        "created_at": "2025-08-01T14:18:56.777624Z"
    },
    ...
]
```