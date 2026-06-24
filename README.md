# LogLint

Кастомный линтер для Go, который проверяет стиль логов и бьет по рукам, если в них летят секреты.

Поддерживает стандартный `log`, а также `log/slog` и `go.uber.org/zap`.

Работает автономно как CLI-утилита, а также встраивается в `golangci-lint` (и как `.so` плагин, и через новую Module Plugin System).

## Что умеет

Проверяет 4 основных правила:

1. **Lowercase:** Сообщение должно начинаться с маленькой буквы (`"starting server"`, а не `"Starting server"`).
2. **English Only:** Сообщения только на английском.
3. **No Emoji/Special Chars:** Запрещены эмодзи, всякие префиксы типа `error:` и мусор вроде `!`, `?` (сделано через быстрый O(1) allow-list).
4. **No Sensitive Data:** Ищет аргументы (переменные/мапы) с именами вроде `password`, `token`, `secret`, чтобы они не улетели в stdout.

*Для первых трех правил есть авто-фиксы (`-fix`)*.

---

## Как запустить

### 1. Как standalone CLI

Самый простой вариант, чтобы просто прогнать по проекту или починить логи авто-фиксом:

```bash
# Обычный запуск
go run ./cmd/loglint/main.go ./...

# С авто-фиксом
go run ./cmd/loglint/main.go -fix ./...
```

### 2. Через golangci-lint (Module Plugin System)

Рекомендуемый способ. Собирает бинарь линтера с вшитым плагином, чтобы не возиться с `.so` файлами.

Создаем `.custom-gcl.yml` в корне вашего проекта:

```yaml
version: v1.57.2
plugins:
  - module: github.com/selectel-tasks/loglint/pkg/plugin
    import: github.com/selectel-tasks/loglint/pkg/plugin
    version: latest
```

Собираем:

```bash
golangci-lint custom
```

Появится бинарник `./custom-gcl`. Используйте его вместо обычного `golangci-lint`.

### 3. Как легаси плагин (.so)

Сборка `.so` плагина. Убедитесь, что версия Go совпадает с той, которой собран ваш `golangci-lint`.

```bash
go build -buildmode=plugin -o loglint.so ./plugin
```

Подключаем в `.golangci.yml`:

```yaml
linters:
  enable:
    - loglint

linters-settings:
  custom:
    loglint:
      path: ./loglint.so
      description: "Custom log validator"
      original-url: github.com/selectel-tasks/loglint
```

---

## Конфиг

Если подключаете через `.golangci.yml`, можно подтюнить правила и добавить свои стоп-слова для секретов.

```yaml
linters-settings:
  custom:
    loglint:
      path: ./loglint.so
      settings:
        rules:
          lowercase: true         
          english_only: true      
          no_special_chars: true   
          no_sensitive: true       
        
        # Кастомные стоп-слова для поиска в именах переменных
        sensitive_keywords:        
          - password
          - token
          - api_key
          - secret
        
        # Можно даже регулярки, если очень нужно
        sensitive_patterns:        
          - "^secret_.*$"
          - "^key_[0-9]+$"
          
        # Запрещенные спецсимволы
        forbidden_chars: "!?"      
```

---

## Тесты

Тесты гоняются через дефолтный `analysistest`. Чтобы не лезть в сеть при парсинге, для `zap` накинул мок в `testdata`.

```bash
go test -v ./pkg/analyzer/...
```
