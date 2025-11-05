Сервис позволяет:
создавать события (/events)
получать данные о событиях (/events/{id})
бронировать место (/events/{id}/book?user={userId})
подтверждать бронь (/events/{id}/confirm?book={bookId})
создавать пользователей (/users)


Создание пользователя
POST /users
Пример запроса:
{
  "email": "john@example.com",
  "name": "John Doe",
  "role": "user",
  "telegram_id": "@john_doe"
}

Пример ответа:
{
  "id": "c1a6f4c1-b7f0-4e77-8c6e-123456789abc"
}


Создание события
POST /events
Пример запроса:
{
  "title": "Go Meetup",
  "description": "Встреча разработчиков на Go",
  "date": "2025-11-15T18:00:00Z",
  "total_seats": 100,
  "available_seats": 100,
  "requires_payment": false,
  "booking_ttl": 3600
}

Пример ответа:
{
  "status": "ok",
  "eventId": "f8d5a2c4-88e7-4d7f-b3e8-5678abcd1234"
}


Получение события
GET /events/{id}
Пример запроса:
GET /events/f8d5a2c4-88e7-4d7f-b3e8-5678abcd1234
Пример ответа:
{
  "id": "f8d5a2c4-88e7-4d7f-b3e8-5678abcd1234",
  "title": "Go Meetup",
  "description": "Встреча разработчиков на Go",
  "date": "2025-11-15T18:00:00Z",
  "total_seats": 100,
  "available_seats": 95,
  "requires_payment": false,
  "booking_ttl": 3600,
  "created_at": "2025-11-01T12:00:00Z"
}


Создание брони
POST /events/{eventId}/book?user={userId}
Пример запроса:
POST /events/f8d5a2c4-88e7-4d7f-b3e8-5678abcd1234/book?user=c1a6f4c1-b7f0-4e77-8c6e-123456789abc
Пример ответа:
{
  "status": "ok",
  "book_id": "bb3e4b1e-55aa-47c3-a612-8d2d8892c1f0"
}


Подтверждение брони
POST /events/{eventId}/confirm?book={bookId}
Пример запроса:
POST /events/f8d5a2c4-88e7-4d7f-b3e8-5678abcd1234/confirm?book=bb3e4b1e-55aa-47c3-a612-8d2d8892c1f0
Пример ответа:
{
  "status": "ok"
}
Если бронь не имеет статус pending, сервер вернёт ошибку:
{
  "error": "cannot confirm booking with status=\"confirmed\""
}