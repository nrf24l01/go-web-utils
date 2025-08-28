# Утилиты для бекенд разработки на го(echo)

## Классы утилит
### S3 утилиты
[Основной файл](https://github.com/nrf24l01/go-web-utils/pkg/s3util/s3.go)
#### Документация для pkg/s3util/s3.go

Описание
- Пакет s3util содержит обёртку над minio.Client для простого получения presigned URL'ов на загрузку, управления метками (tags) объектов и получения публичных URL.

Типы
- S3Config
  - Endpoint string — хост S3/MinIO (host:port).
  - AccessKey string — ключ доступа.
  - SecretKey string — секретный ключ.
  - UseSSL bool — использовать TLS (true/false).
  - BaseURL string — базовый URL для замены хоста в presigned URL (может быть пустым).

- Client
  - Обёртка, содержащая внутренний minio.Client и baseURL для формирования публичных ссылок.

Конструктор
- New(cfg S3Config) (*Client, error)
  - Создаёт и возвращает *Client. Валидирует, что Endpoint, AccessKey и SecretKey заданы.
  - При ошибке подключения возвращает ошибку.

Методы
- GeneratePresignedPutURL(ctx context.Context, bucket string, expires time.Duration) (string, string, error)
  - Генерирует уникальное имя объекта (UUID) и presigned PUT URL с указанным временем жизни.
  - Возвращает: имя объекта, URL для загрузки, ошибку.
  - Если в S3Config.BaseURL указан base URL, хост в сгенерированном presigned URL заменяется на baseURL (сохранены схема и путь).

- GetPermanentObjectURL(bucket, object string) string
  - Возвращает публичный URL для уже помеченного как «permanent» объекта.
  - Если baseURL пуст — возвращает относительный путь вида "/api/files/{bucket}/{object}".
  - Корректно обрабатывает наличие/отсутствие завершающего '/' в baseURL.

- MarkObjectPermanent(ctx context.Context, bucket, object string) error
  - Получает текущие теги объекта, устанавливает/заменяет тег "status" на "permanent" и сохраняет теги обратно.

- IsObjectTemporary(ctx context.Context, bucket, object string) (bool, error)
  - Проверяет тег "status" и возвращает true если значение равно "temporary".

Вспомогательные детали
- replaceHostWithBaseURL(originalURL, baseURL string) (string, error)
  - Используется для замены схемы/хоста/пути у сгенерированного presigned URL на значения из baseURL.
  - Если baseURL пуст — возвращает оригинальный URL без изменений.

Примеры использования

1) Создание клиента и генерация presigned PUT URL:
```go
// Пример упрощённого использования
cfg := s3util.S3Config{
    Endpoint:  "play.min.io:9000",
    AccessKey: "YOURACCESSKEY",
    SecretKey: "YOURSECRETKEY",
    UseSSL:    true,
    BaseURL:   "https://cdn.example.com/s3",
}
client, err := s3util.New(cfg)
if err != nil {
    // обработка ошибки
}
objectName, presignedURL, err := client.GeneratePresignedPutURL(ctx, "my-bucket", 15*time.Minute)
if err != nil {
    // обработка ошибки
}
// presignedURL можно вернуть клиенту для загрузки, objectName — сохранить в БД
```

2) Пометить загруженный объект как permanent и получить публичный URL:
```go
err = client.MarkObjectPermanent(ctx, "my-bucket", objectName)
if err != nil {
    // обработка ошибки
}
publicURL := client.GetPermanentObjectURL("my-bucket", objectName)
// publicURL — итоговая доступная ссылка
```

Пояснения по тегам
- В реализации используется тег "status" со значениями "temporary" и "permanent".
- MarkObjectPermanent заменяет или добавляет тег "status":"permanent".
- IsObjectTemporary проверяет текущий статус объекта.

Замечания
- Убедитесь, что бакет существует и у пользователя есть права на PutObjectTagging/GetObjectTagging/PresignedPut.
- Если используете сторонний CDN или прокси, укажите BaseURL для формирования корректных ссылок, чтобы пользователи загружали напрямую через ваш домен.
- Убедитесь, что бакет существует и у пользователя есть права на PutObjectTagging/GetObjectTagging/PresignedPut.
- Если используете сторонний CDN или прокси, укажите BaseURL для формирования корректных ссылок, чтобы пользователи загружали напрямую через ваш домен.

### Echo middleware (pkg/echokit)

Файл: pkg/echokit/validatemw.go
Кратко
- Middleware для валидации тела запроса с использованием go-playground/validator и механизма Bind/Validate Echo.
- Сохраняет валидированную структуру в контексте под ключом "validatedBody".

Типы и функции
- CustomValidator
  - Поле Validator *validator.Validate — адаптер для Echo.
  - Метод Validate(i interface{}) error — вызывает Validator.Struct.

- ValidationMiddleware(schemaFactory func() interface{}) echo.MiddlewareFunc
  - Принимает фабрику схемы (функция возвращает новую пустую структуру).
  - Выполняет c.Bind(schema) и c.Validate(schema).
  - При ошибке возвращает 422 (bind) или 400 (validate) с кратким сообщением.
  - При успехе кладёт schema в контекст: c.Set("validatedBody", schema) и пропускает дальше.

Пример использования
```go
// инициализация
e.Validator = &echokit.CustomValidator{Validator: validator.New()}

// регистрация middleware для маршрута
e.POST("/items", func(c echo.Context) error {
    v := c.Get("validatedBody").(*YourRequestType)
    // обработка v
    return c.JSON(200, v)
}, echokit.ValidationMiddleware(func() interface{} { return new(YourRequestType) }))
```

Файл: pkg/echokit/jwtmw.go
Кратко
- Middleware для извлечения и валидации JWT из заголовка Authorization: "Bearer <token>".
- Использует jwtutils.ValidateToken (в проекте — pkg/jwtutil/jwt.go) и кладёт userID в контекст под ключом "userID".

Поведение
- Если header отсутствует или не в формате Bearer — возвращает 401.
- Если секрет пуст — возвращает 500.
- Если валидация не прошла — возвращает 401.
- Ожидает, что claims содержат "user_id" строкой.

Пример использования
```go
e.Use(echokit.JWTMiddleware([]byte(os.Getenv("JWT_SECRET"))))

// в хендлере
userID := c.Get("userID").(string)
```

### JWT утилиты (pkg/jwtutil)

Файл: pkg/jwtutil/jwt.go
Кратко
- Функции для генерации access/refresh токенов и их валидации на основе HMAC (HS256).

Функции
- GenerateAccessToken(userID string, username string, accessSecret []byte) (string, error)
  - Генерирует access token с claims: user_id, username, exp (15 минут), iat.

- GenerateRefreshToken(userID string, refreshSecret []byte) (string, error)
  - Генерирует refresh token с claims: user_id, exp (7 дней), iat.

- ValidateToken(tokenString string, secret []byte) (jwt.MapClaims, error)
  - Парсит и проверяет подпись/алгоритм, возвращает claims при валидном токене.

Пример использования
```go
access, err := jwtutil.GenerateAccessToken(userID, username, []byte(cfg.AccessSecret))
refresh, err := jwtutil.GenerateRefreshToken(userID, []byte(cfg.RefreshSecret))

claims, err := jwtutil.ValidateToken(access, []byte(cfg.AccessSecret))
if err != nil {
    // обработка
}
uid := claims["user_id"].(string)
```

Замечания
- Храните секреты в безопасном месте.
- Следите за временем жизни токенов и реализацией механизма обновления (refresh).

### Passhash (pkg/passhash)

Файл: pkg/passhash/argon2id.go
Кратко
- Утилиты для хеширования паролей с argon2id и сравнения хешей.

Типы и переменные
- Params — параметры argon2 (Memory, Time, Parallelism, SaltLength, KeyLength).
- DefaultParams — рекомендованные параметры по умолчанию.

Функции
- HashPassword(password string, p *Params) (string, error)
  - Генерирует соль, вычисляет argon2id хеш и возвращает строку в формате:
    $argon2id$v=19$m=<memory>,t=<time>,p=<parallelism>$<salt_b64>$<hash_b64>

- CheckPassword(password string, encodedHash string) (bool, error)
  - Разбирает закодированную строку, декодирует параметры/соль/хеш и сравнивает безопасно.

- subtleCompare(a, b []byte) bool
  - Постоянное по времени сравнение байтовых слайсов.

Пример использования
```go
hash, err := passhash.HashPassword("secret123", passhash.DefaultParams)
ok, err := passhash.CheckPassword("secret123", hash)
if ok {
    // пароль верный
}
```

Замечания по безопасности
- Используйте сильные параметры (DefaultParams подходят для многих случаев, но подстройте под вашу инфраструктуру).
- Никогда не логируйте сырые пароли.
- Сохраняйте строку-хеш в БД и проверяйте её при логине.
