# CLAUDE.md - Руководство для AI-ассистента

## Обзор проекта

**pingerbot** (Mr. Pinger) — Telegram-бот для групповых оповещений с поддержкой тегов. Позволяет упоминать всех пользователей группы командой `/ping` (аналог `@all` в Slack).

**Стек:** Go 1.18, PostgreSQL, Docker, GitHub Actions

**Ключевые особенности:**
- Privacy mode включён (бот видит только команды)
- Username-based система (ID пользователей не хранятся)
- Гибкое тегирование (#backend, #devops, #all и т.д.)

## Архитектура

```
cmd/main.go                    # Точка входа, конфигурация
├── pkg/telegram/              # Telegram Bot SDK (самописный)
│   ├── bot.go                 # Polling, routing событий
│   ├── api.go                 # HTTP клиент для Telegram API
│   ├── types.go               # Структуры данных (Update, Message, etc.)
│   ├── ctx.go                 # Контексты для обработчиков
│   └── constants.go           # Константы ParseMode
├── internal/
│   ├── handlers/              # Обработчики событий
│   │   ├── commands/          # /add, /ping, /ls, etc.
│   │   ├── private-message.handler.go
│   │   ├── bot-joins.handler.go
│   │   ├── bot-leaves.handler.go
│   │   ├── user-joins.handler.go
│   │   └── user-leaves.handler.go
│   ├── state/state.go         # Слой работы с БД
│   └── messages/messages.go   # Шаблоны сообщений
└── migrations/                # SQL миграции
```

**Поток данных:**
```
Telegram API → pkg/telegram/Bot (polling) → Handlers → State → PostgreSQL
```

## Схема БД

```sql
-- Группы (чаты)
CREATE TABLE groups (
  id varchar PRIMARY KEY,      -- Chat ID как строка
  name varchar                 -- Название группы
);

-- Участники с тегами (многие-ко-многим)
CREATE TABLE members (
  username varchar,
  group_id varchar REFERENCES groups(id) ON DELETE CASCADE,
  tag varchar,
  PRIMARY KEY (username, group_id, tag)
);
```

## Команды разработки

Используется [Task](https://taskfile.dev/) (`brew install go-task`).

```bash
# Запуск
task start              # go run cmd/main.go (dev)
task build              # Сборка в build/pingerbot
task start-built        # Сборка + запуск бинарника

# Линтинг
task lint               # golangci-lint run

# Миграции (требуется golang-migrate CLI)
task add-migration -- name   # Создать новую миграцию
task migrate            # Применить миграции (up)
task revert             # Откатить миграции (down)

# Docker Compose
task infra:up           # Поднять PostgreSQL
task infra:down         # Остановить PostgreSQL
task up                 # Поднять всё (infra + app)
task down               # Остановить всё
task rebuild            # Пересобрать и запустить
task logs               # Логи всех сервисов
task clean              # Удалить контейнеры и volumes

# Docker (standalone)
task docker:build       # Собрать образ
task docker:run         # Запустить контейнер
task docker:clean-run   # Пересобрать и запустить
```

## Конфигурация

Переменные окружения (см. `.env.example`):

```bash
BOT_TOKEN=<telegram_bot_token>
LONG_POLLING_TIMEOUT=30s

# PostgreSQL
POSTGRES_USER=pingerbot
POSTGRES_PASSWORD=pingerbot
POSTGRES_DB=pingerbot
POSTGRES_PORT=5447
```

## Функционал бота

**Команды в группе:**
| Команда | Описание |
|---------|----------|
| `/addme [#tags...]` | Добавить себя (по умолчанию в #all) |
| `/add @user1 @user2 [#tags...]` | Добавить других пользователей |
| `/removeme [#tags...]` | Удалить себя |
| `/remove @user1 @user2 [#tags...]` | Удалить пользователей |
| `/ping [#tags...]` | Упомянуть пользователей из тегов |
| `/ls [#tags...]` | Показать список по тегам |

**Автоматические события:**
- При входе бота в группу — регистрация группы, приветствие
- При входе пользователя — добавление в #all (если есть username)
- При выходе пользователя — удаление из всех тегов
- При выходе бота — удаление группы (каскадное удаление)

## Паттерны и конвенции

### Структура обработчиков

```go
// internal/handlers/commands/example.handler.go
func HandleExample(ctx *telegram.CommandCtx, state *state.State) {
    // 1. Извлечь данные из ctx (mentions, hashtags)
    // 2. Выполнить бизнес-логику через state
    // 3. Отправить ответ через ctx.ReplyTxt() или ctx.SendToChat()
}
```

### Работа с контекстами

- `MsgCtx` — базовый контекст сообщения
- `CommandCtx` — расширяет MsgCtx, содержит parsed команду
- `JoinCtx` / `LeaveCtx` — для событий входа/выхода

Методы контекста:
```go
ctx.ReplyTxt(text)           // Ответить на сообщение
ctx.SendToChat(msg)          // Отправить в чат
ctx.Msg.Entities             // Entities для парсинга @mentions, #hashtags
```

### Парсинг entities

```go
// Извлечение упоминаний (@username)
for _, e := range ctx.Msg.Entities {
    if e.Type == "mention" {
        username := ctx.Msg.Text[e.Offset+1 : e.Offset+e.Length] // убираем @
    }
}

// Извлечение хэштегов (#tag)
for _, e := range ctx.Msg.Entities {
    if e.Type == "hashtag" {
        tag := ctx.Msg.Text[e.Offset : e.Offset+e.Length]
    }
}
```

### State слой

```go
state.RememberGroup(chat)                    // Регистрация группы
state.ForgetGroup(chat)                      // Удаление группы
state.RememberMember(groupId, username, tags) // Добавить пользователя
state.ForgetMember(groupId, username, tags)  // Удалить пользователя
state.GetKnownMembers(groupId, tags)         // Получить @usernames для пинга
state.ListGroupMembers(groupId, tags)        // Структурированный список
```

### Преобразование Chat ID

```go
// Chat ID — int64, в БД хранится как string
groupId := strconv.FormatInt(ctx.Msg.Chat.Id, 10)
```

## При внесении изменений

### Добавление новой команды

1. Создать файл `internal/handlers/commands/<name>.handler.go`
2. Реализовать функцию `Handle<Name>(ctx *telegram.CommandCtx, state *state.State)`
3. Зарегистрировать в `pkg/telegram/bot.go` в методе `handleMessage()`:
   ```go
   case "<name>":
       commands.Handle<Name>(cmdCtx, state)
   ```

### Изменение схемы БД

1. `task add-migration -- name` — создаст файлы up/down
2. Написать SQL в `migrations/<timestamp>_<name>.up.sql`
3. Написать откат в `migrations/<timestamp>_<name>.down.sql`
4. `task migrate` — применить
5. Обновить методы в `internal/state/state.go`

### Добавление нового события

1. Проверить тип Update в `pkg/telegram/types.go`
2. Добавить обработку в `pkg/telegram/bot.go` метод `handleUpdate()`
3. Создать обработчик в `internal/handlers/`

## Важные детали

- **Тег по умолчанию:** `#all` — используется если теги не указаны
- **Privacy mode:** Бот НЕ видит обычные сообщения, только команды
- **CASCADE удаление:** При удалении группы все её участники удаляются автоматически
- **Идемпотентность:** Повторное добавление не создаёт дубликаты (PRIMARY KEY)
- **Username обязателен:** Пользователи без username не могут быть добавлены

## Тестирование

> **TODO:** Тесты ещё не написаны

При добавлении тестов использовать стандартную библиотеку `testing`.

## CI/CD

GitHub Actions (`.github/workflows/deploy.yml`):
1. **Build:** компиляция + golangci-lint
2. **Deploy:** SCP на сервер + systemctl restart pingerbot

Секреты: `MCS_HOST`, `MCS_SSH_KEY`

## Известные TODO

- [ ] Unit тесты
- [ ] Мониторинг (Prometheus/Grafana)
- [ ] Webhooks вместо polling
- [ ] Параллельная обработка чатов
