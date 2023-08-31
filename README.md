# avito-tech-service
Сервис динамического сегментирования пользователей

Для запуска сервиса необходимо разместить config файл в корне проекта. Пример заполнения yaml-файла ниже:
    // Уровень запуска
        env: "local" # local, dev, prod
    // Настройка БД
        storage_path: "/project_name/storagepath"
        database_url: "host=localhost user=username password=userpass dbname=dbname sslmode=disable"
    // Параметры HTTP сервера:
        http_server:
        address : "adress:port"
        timeout: 4s
        idle_timeout: 60s
        user: "username"     // параметры авторизации
        password: "password"

При первичном запуске сервиса инициализируются три таблицы: 
    USERS (user_id, created_at) - таблица для ведения пользователей с датой создания;
    SEGMENTS (segment_name, created_at) - таблица для ведения сегментов с датой создания;
    USER_SEGMENTS (user_id, segment_name) - таблица принадлежности пользователя к конкретному сегменту (имеет внешние ключи с таблицами выше).

Изначально рекомендуется завести пользователей в таблицу USERS путём отправки POST-запросов на адресс "service_adress/users" JSON в формате:
{
    "user_id": XXX
}
, где XXX - любое ЧИСЛОВОЕ выражение. ВНИМАНИЕ! Сервис хранит имена пользователей в чиловом формате, последовательность - произовальная, задаётся пользователем.

Затем необходимо завести сегменты в таблице SEGMENTS, аналогично, на адресс "service_adress/segments", формат JSON:
{
    "slug": "SEGMENT_NAME"
}
.

После этих настроек можно делать соответствие между пользователями и сегментами, путем отправки POST/DELETE запросов на адресс "service_adress/users/id=XXX", где XXX - заведенный ранее в таблице USERS id пользователя, JSON формат:
{
    "segments": [
        { "slug": "segment1" },
        { "slug": "segment2" },
        { "slug": "segment3" }
  ]
}
.

Получить информацию о сегментах, в которых состоит тот или иной пользователь, можно путем отправки GET запроса на "service_adress/users/id=XXX".
Пример запроса и ответа.
    GET: http://127.0.0.1:8080/users/id=1000
    JSON ответ:
    {
    "status": "OK",
    "user_id": 1000,
    "segments": [
        "AVITO_VOICE_MESSAGES",
        "AVITO_PERFORMANCE_VAS",
        "AVITO_DISCOUNT_50"
    ],
    "Method": "GET"
}

В случае, если пользователь заведен, но не принадлежит ни одному сегменту, ответ будет:
{
    "status": "OK",
    "user_id": 100,
    "segments": null,
    "Method": "GET"
}

Реализован просто функциональный тест, который создаёт случайного пользователя, создаёт случайный сегмент, добавляет этот сегмент к пользователю и запрашивает сегменты, которые относятся к данному пользователю.
