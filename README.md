# gomax

Go-клиент для мессенджера Max. Порт библиотеки [PyMax](https://github.com/MaxApiTeam/PyMax) с Python на Go с помощью LLM.

Используйте на свой страх и риск,

## Возможности

- Авторизация по телефону и SMS-коду через `Client`
- QR-авторизация web-клиента через `WebClient`
- Роутеры, фильтры, `on_start`, raw-события и typed events
- Сообщения: отправка, forward, reply, реакции, pin, read, delete и история
- Чаты, группы, участники, invite-ссылки и настройки групп
- Пользователи, контакты, профиль, папки, активные сессии и 2FA
- Вложения: `Photo`, `File`, `Video`
- SQLite-сессии, sync-state, reconnect
- Debug-логирование через `-debug`

## Установка

```bash
go get github.com/akovardin/gomax
```

## Быстрый старт

### TCP-клиент с SMS-авторизацией

```go
package main

import (
    "log"

    "github.com/akovardin/gomax"
    "github.com/akovardin/gomax/types"
)

func main() {
    client := gomax.NewClient(
        "+79990000000",          // phone
        "session-mobile.db",     // sessionName
        ".",                     // workDir
        nil,                     // extraConfig
        nil,                     // smsCodeProvider
        nil,                     // passwordProvider
    )

    client.OnStart()(func(c interface{}) error {
        cl := c.(*gomax.Client)
        if me := cl.Me(); me != nil && me.Contact != nil {
            log.Printf("Запущен: %v", me.Contact.ID)
        }
        return nil
    })

    client.OnMessage()(func(event interface{}, c interface{}) error {
        msg := event.(*types.Message)
        log.Printf("chat=%v sender=%v text=%q id=%v", msg.ChatID, msg.Sender, msg.Text, msg.ID)

        if msg.ChatID != nil && msg.Text != "" {
            _, err := client.API().Messages.SendMessage(
                msg.ChatID.Int(), "Привет от gomax", nil, nil, false,
            )
            return err
        }
        return nil
    })

    if err := client.Start(); err != nil {
        log.Fatal(err)
    }
}
```

Запуск:
```bash
go run cmd/mobile/main.go
```

### WebSocket-клиент с QR-авторизацией

```go
package main

import (
    "log"

    "github.com/akovardin/gomax"
)

func main() {
    client := gomax.NewWebClient(
        "session-web.db",    // sessionName
        ".",                 // workDir
        nil,                 // extraConfig
        nil,                 // qrProvider
        nil,                 // passwordProvider
    )

    client.OnStart()(func(c interface{}) error {
        log.Println("Web-клиент запущен")
        return nil
    })

    if err := client.Start(); err != nil {
        log.Fatal(err)
    }
}
```

## Роутеры и фильтры

```go
package main

import (
    "fmt"

    "github.com/akovardin/gomax"
    "github.com/akovardin/gomax/dispatch"
    "github.com/akovardin/gomax/types"
)

func main() {
    client := gomax.NewClient("+79990000000", "session-mobile.db", ".", nil, nil, nil)

    router := dispatch.NewRouter()

    command := func(name string) dispatch.FilterCallback {
        return func(event interface{}) bool {
            msg, ok := event.(*types.Message)
            return ok && msg.Text == "/"+name
        }
    }

    router.OnMessage(command("start"))(func(event interface{}, c interface{}) error {
        msg := event.(*types.Message)
        _, err := client.API().Messages.SendMessage(msg.ChatID.Int(), "Привет!", nil, nil, false)
        return err
    })

    router.OnMessage(command("me"))(func(event interface{}, c interface{}) error {
        msg := event.(*types.Message)
        if me := client.Me(); me != nil && me.Contact != nil {
            _, err := client.API().Messages.SendMessage(msg.ChatID.Int(),
                fmt.Sprintf("Ваш ID: %v", me.Contact.ID), nil, nil, false)
            return err
        }
        return nil
    })

    client.IncludeRouter(router)
    client.Start()
}
```

### Множественные фильтры

```go
onlyText := func(event interface{}) bool {
    msg, ok := event.(*types.Message)
    return ok && msg.Text != ""
}

router.OnMessage(onlyText)(func(event interface{}, c interface{}) error {
    msg := event.(*types.Message)
    _, err := client.API().Messages.SendMessage(msg.ChatID.Int(), "Текст получил", nil, nil, false)
    return err
})
```

## Работа с сообщениями

```go
// Отправка
msg, err := client.API().Messages.SendMessage(chatID, "текст", nil, nil, false)

// Ответ на сообщение (reply)
replyTo := messageID
msg, err = client.API().Messages.SendMessage(chatID, "ответ", &replyTo, nil, false)

// Forward сообщения
msg, err = client.API().Messages.ForwardMessage(targetChatID, messageID, sourceChatID, true)

// Редактирование
msg, err = client.API().Messages.EditMessage(chatID, messageID, "новый текст", nil)

// Удаление
ok, err := client.API().Messages.DeleteMessage(chatID, []int{messageID}, false)

// Pin сообщения
ok, err = client.API().Messages.PinMessage(chatID, messageID, true)

// Реакции
info, err := client.API().Messages.AddReaction(chatID, messageID, "👍")
info, err = client.API().Messages.RemoveReaction(chatID, messageID)
reactions, err := client.API().Messages.GetReactions(chatID, []int{messageID})

// Отметка прочтения
state, err := client.API().Messages.ReadMessage(messageID, chatID)

// Получение сообщений
msg, err = client.API().Messages.GetMessage(chatID, messageID)
msgs, err := client.API().Messages.GetMessages(chatID, []int{id1, id2})

// История чата
history, err := client.API().Messages.FetchHistory(chatID, 0, 50)
```

## Отправка файлов

```go
import (
    "bytes"

    "github.com/akovardin/gomax/files"
)

// Фото
photo, err := files.NewPhotoFromPath("image.jpg")
photoData, _ := photo.Read()
photoPayload, err := client.API().Uploads.UploadPhoto(
    bytes.NewReader(photoData), photo.Name(), false,
)
msg, err := client.API().Messages.SendMessage(chatID, "Фото",
    nil, []interface{}{photoPayload}, false,
)

// Файл
file, _ := files.NewFileFromPath("report.pdf")
fileData, _ := file.Read()
fileSize, _ := file.Size()
filePayload, err := client.API().Uploads.UploadFile(
    bytes.NewReader(fileData), file.Name(), fileSize,
)
msg, err = client.API().Messages.SendMessage(chatID, "Документ",
    nil, []interface{}{filePayload}, false,
)

// Видео
video, _ := files.NewVideoFromPath("clip.mp4")
videoData, _ := video.Read()
videoSize, _ := video.Size()
videoPayload, err := client.API().Uploads.UploadVideo(
    bytes.NewReader(videoData), video.Name(), videoSize,
)
msg, err = client.API().Messages.SendMessage(chatID, "Видео",
    nil, []interface{}{videoPayload}, false,
)
```

## Получение вложений

```go
import (
    "encoding/json"
    "github.com/akovardin/gomax/types/attachments"
)

client.OnMessage()(func(event interface{}, c interface{}) error {
    msg := event.(*types.Message)
    if msg.ChatID == nil {
        return nil
    }

    for _, attach := range msg.Attaches {
        switch attach.Type {
        case "PHOTO":
            var photo attachments.PhotoAttachment
            json.Unmarshal(attach.Raw, &photo)
            log.Printf("photo: %v %s", photo.PhotoID, photo.BaseURL)

        case "FILE":
            var file attachments.FileAttachment
            json.Unmarshal(attach.Raw, &file)
            fileInfo, _ := client.API().Messages.GetFileByID(
                msg.ChatID.Int(), msg.ID.Int(), file.FileID,
            )
            if fileInfo != nil {
                log.Printf("file: %s", fileInfo.URL)
            }
        }
    }
    return nil
})
```

## Сервисные события

```go
import "github.com/akovardin/gomax/types"

client.OnTyping()(func(event interface{}, c interface{}) error {
    e := event.(*types.TypingEvent)
    log.Printf("typing: chat=%v user=%v", e.ChatID, e.UserID)
    return nil
})

client.OnPresence()(func(event interface{}, c interface{}) error {
    e := event.(*types.PresenceEvent)
    log.Printf("presence: user=%v status=%v", e.UserID, e.Presence.Status)
    return nil
})

client.OnMessageRead()(func(event interface{}, c interface{}) error {
    e := event.(*types.MessageReadEvent)
    log.Printf("read: chat=%v mark=%v", e.ChatID, e.Mark)
    return nil
})

client.OnReactionUpdate()(func(event interface{}, c interface{}) error {
    e := event.(*types.ReactionUpdateEvent)
    log.Printf("reactions: msg=%s count=%v", e.MessageID, e.TotalCount)
    return nil
})

client.OnMessageDelete()(func(event interface{}, c interface{}) error {
    e := event.(*types.MessageDeleteEvent)
    log.Printf("deleted in chat: %v", e.ChatID)
    return nil
})

client.OnMessageEdit()(func(event interface{}, c interface{}) error {
    msg := event.(*types.Message)
    log.Printf("edited: %s", msg.Text)
    return nil
})
```

## Raw-события

```go
import "github.com/akovardin/gomax/protocol"

client.OnRaw()(func(event interface{}, c interface{}) error {
    frame := event.(*protocol.InboundFrame)
    log.Printf("opcode: %d, payload: %v", frame.Opcode, frame.Payload)
    return nil
})
```

## Обработчики ошибок и отключений

```go
import (
    "github.com/akovardin/gomax"
    "github.com/akovardin/gomax/api/core"
    "github.com/akovardin/gomax/dispatch"
)

client.OnError(gomax.ErrorScopeGlobal)(func(err error, ctx *dispatch.ErrorContext) error {
    if apiErr, ok := err.(*core.ApiError); ok && apiErr.ErrorCode() == "FAIL_LOGIN_TOKEN" {
        return client.Relogin(false, true)
    }
    log.Printf("error: %v", err)
    return nil
})

client.OnDisconnect()(func(err error, reconnect bool, delay float64) error {
    log.Printf("connection lost: %v, reconnect=%v, delay=%.1f", err, reconnect, delay)
    return nil
})
```

## Конфигурация и debug

```go
client := gomax.NewClient(
    "+79990000000",
    "account.db",
    "cache",
    &gomax.ExtraConfig{
        LogLevel:       "debug",
        Reconnect:      true,
        ReconnectDelay: 3.0,
        Telemetry:      false,
    },
    nil,
    nil,
)
```


```bash
go run cmd/mobile/main.go
go run cmd/web/main.go
```

## Токен-авторизация

```go
client := gomax.NewClient(
    "+79990000000",
    "session-mobile.db",
    "cache",
    &gomax.ExtraConfig{Token: "YOUR_TOKEN"},
    nil,
    nil,
)
```

## Relogin

```go
// Переавторизация со сбросом токена
client.Relogin(true, true)

// Переавторизация с сохранением токена из ExtraConfig
client.Relogin(false, true)
```

## Регистрация нового аккаунта

```go
client := gomax.NewClient(
    "+79990000000",
    "session-mobile.db",
    "cache",
    &gomax.ExtraConfig{
        RegistrationConfig: &types.RegistrationConfig{
            FirstName: "Max",
            LastName:  "User",
        },
    },
    nil,
    nil,
)
```

## Кастомные провайдеры

### SMS-код из очереди

```go
type QueueSmsProvider struct {
    queue chan string
}

func (p *QueueSmsProvider) GetCode(phone string) (string, error) {
    fmt.Printf("Waiting for SMS code for %s...\n", phone)
    return <-p.queue, nil
}

func (p *QueueSmsProvider) SetCode(code string) {
    p.queue <- code
}

provider := &QueueSmsProvider{queue: make(chan string, 1)}
client := gomax.NewClient("+79990000000", "s.db", ".", nil, provider, nil)
```

### Пароль 2FA из переменной окружения

```go
type EnvPasswordProvider struct{}

func (p *EnvPasswordProvider) GetPassword(hint string) (string, error) {
    return os.Getenv("MAX_2FA_PASSWORD"), nil
}

client := gomax.NewClient("+79990000000", "s.db", ".",
    nil, nil, &EnvPasswordProvider{},
)
```

### Кастомный QR-обработчик

```go
type PrintQRHandler struct{}

func (h *PrintQRHandler) ShowQR(qrURL string) error {
    fmt.Printf("Откройте ссылку: %s\n", qrURL)
    return nil
}

client := gomax.NewWebClient("web.db", ".",
    nil, &PrintQRHandler{}, nil,
)
```

## Работа с чатами

```go
// Создание группы
chat, msg, err := client.API().Chats.CreateGroup("Моя группа", []int{111, 222}, true)

// Пригласить участников
chat, err = client.API().Chats.InviteUsersToGroup(chatID, []int{333, 444}, true)

// Удалить участников
ok, err := client.API().Chats.RemoveUsersFromGroup(chatID, []int{333}, 0)

// Настройки группы
err = client.API().Chats.ChangeGroupSettings(chatID, strPtr("Новое название"), strPtr("Описание"))

// Получить чат
chat, err = client.API().Chats.GetChat(chatID)

// Получить несколько чатов
chats, err := client.API().Chats.GetChats([]int{id1, id2})

// Все чаты (с сервера)
chats, err = client.API().Chats.FetchChats(nil)

// Покинуть / удалить
err = client.API().Chats.LeaveGroup(chatID)
err = client.API().Chats.DeleteChat(chatID)

// Invite-ссылки
chat, err = client.API().Chats.JoinGroup("https://max.ru/join/XXXX")
chat, err = client.API().Chats.ResolveGroupByLink("https://max.ru/join/XXXX")
chat, err = client.API().Chats.ReworkInviteLink(chatID)

// Запросы на вступление
members, err := client.API().Chats.GetJoinRequests(chatID, 10)
ok, err = client.API().Chats.ConfirmJoinRequests(chatID, []int{userID})
ok, err = client.API().Chats.DeclineJoinRequests(chatID, []int{userID})

func strPtr(s string) *string { return &s }
```

## Пользователи и контакты

```go
// Получить пользователя
user, err := client.API().Users.GetUser(123)

// Получить нескольких пользователей
users, err := client.API().Users.GetUsers([]int{123, 456})

// Поиск по номеру телефона
user, err = client.API().Users.SearchByPhone("+79990000000")

// Добавить / удалить контакт
err = client.API().Users.AddContact(123)
err = client.API().Users.RemoveContact(123)

// Импорт контактов
contacts, err := client.API().Users.ImportContacts([]types.ContactInfo{
    {Phone: "+79990000000", FirstName: "Ada", LastName: strPtr("Lovelace")},
})

// ID личного чата (XOR)
chatID := client.API().Users.GetChatID(123, 456)
```

## Профиль и аккаунт

```go
// Изменить профиль
ok, err := client.API().Self.ChangeProfile("Alex", "PyMax", "Testing", nil, "")

// С фото
photo, _ := files.NewPhotoFromPath("avatar.jpg")
photoData, _ := photo.Read()
photoPayload, _ := client.API().Uploads.UploadPhoto(
    bytes.NewReader(photoData), photo.Name(), true,
)
ok, err = client.API().Self.ChangeProfile("Alex", "", "",
    photoPayload, photoPayload.PhotoToken,
)

// Папки
update, err := client.API().Self.CreateFolder("Работа", []int{chatID1, chatID2}, nil)
folders, err := client.API().Self.GetFolders(0)
update, err = client.API().Self.UpdateFolder(folderID, "Новое имя", []int{chatID1}, nil)
update, err = client.API().Self.DeleteFolder(folderID)

// Сессии
sessions, err := client.API().Users.GetSessions()
for _, s := range sessions {
    log.Printf("session: %v %s %s current=%v", s.ID, s.DeviceName, s.DeviceType, s.Current)
}
ok, err = client.API().Self.CloseAllSessions()

// Выход
ok, err = client.API().Self.Logout()
```

## 2FA

```go
// Установить 2FA
ok, err := client.API().Auth.Set2FA("strong-password", "", "мой пароль")

// С email
ok, err = client.API().Auth.Set2FA("strong-password", "user@example.com", "мой пароль")

// Удалить 2FA
ok, err = client.API().Auth.Remove2FA("strong-password")

// Проверить 2FA
ok, err = client.API().Auth.Check2FA()
```

## Markdown-форматирование

```go
import "github.com/akovardin/gomax/formatting"

formatter := formatting.NewFormatter()
cleanText, elements := formatter.FormatMarkdown("**Важное:** проверьте [документацию](https://example.com)")

// Отправка с элементами форматирования
msg, err := client.API().Messages.SendMessage(chatID, cleanText, nil,
    []interface{}{map[string]interface{}{"elements": elements}}, false,
)
```

## Структура проекта

```
gomax/
├── config.go            # ExtraConfig, ClientConfig, UA-генерация
├── enums.go             # EventType, ChatType, MessageStatus и др.
├── errors.go            # Базовые типы ошибок
├── app.go               # App — runtime (invoke, ping, dispatch)
├── base.go              # BaseClient — общая логика, reconnect
├── client.go            # Client — TCP + SMS-авторизация
├── client_web.go        # WebClient — WebSocket + QR-авторизация
├── logging.go           # SetLogLevel (публичная обёртка)
├── logging/             # Уровневое логирование (LogDebug/LogInfo/LogWarn/LogError)
├── protocol/            # Сетевой протокол (msgpack+LZ4/Zstd + JSON)
├── transport/           # TCP (TLS/proxy) и WebSocket транспорт
├── connection/          # ConnectionManager, PendingRequests
├── api/                 # API-сервисы (auth, messages, chats, users, ...)
├── types/               # Domain-модели и события (FlexInt, Message, Chat, ...)
├── session/             # SQLite-сессии
├── auth/                # Auth-флоу (SMS, QR) и провайдеры
├── dispatch/            # Router, Dispatcher, EventMapper
├── files/               # Photo, Video, File (для загрузки)
├── formatting/          # Markdown → Max Elements
└── telemetry/           # Телеметрия (NAV/PERF)
```

## Типы данных

Сервер Max может возвращать числовые ID как строки или числа. Для обработки используется `FlexInt` — тип, который при десериализации принимает и `int`, и `string`:

```go
type FlexInt int

// Поддерживает JSON (строки и числа) и msgpack (int/string)
func (f *FlexInt) UnmarshalJSON(data []byte) error { ... }
func (f *FlexInt) DecodeMsgpack(dec *msgpack.Decoder) error { ... }

// Для получения int: id.Int()
// Для форматирования: %v или id.Int()
```

При отправке значений в API используйте `int` напрямую. При чтении из ответов — `FlexInt.Int()`.

## Запуск примеров

```bash
# TCP-клиент (SMS-авторизация)
go run cmd/mobile/main.go

# WebSocket-клиент (QR-авторизация)
go run cmd/web/main.go
```
