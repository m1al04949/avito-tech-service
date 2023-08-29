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
    USERS (user_id, created_at) - таблица для ведения пользователей
    SEGMENTS (segment_name, created_at) - таблица для ведения сегментов
    USER_SEGMENTS (user_id, segment_name) - таблица принадлежности пользователя к конкретному сегменту (имеет внешние ключи с таблицами выше)

Изначально рекомендуется завести пользователей в таблицу USERS путём отправки POST-запросов на адресс:"adress/users" JSON-документв в формате:
{
    "user_id": XXX
}
, где XXX - любое ЧИСЛОВОЕ выражение. ВНИМАНИЕ! Сервис хранит имена пользователей в чиловом формате, последовательность - произовальная.

Затем необходимо завести сегменты в таблице SEGMENTS, аналогично, на адресс "adress/segments", формат JSON:
{
    "slug": "SEGMENT_NAME"
}
.

После этих настроек можно делать соответствие между пользователями и сегментами, путем отправки POST/DELETE запросов на адресс "adress/users/id=XXX", где XXX - заведенный ранее в таблице USERS id пользователя, JSON-документов формата:
{
    "segments": [
        { "slug": "segment1" },
        { "slug": "segment2" },
        { "slug": "segment3" }
  ]
}
