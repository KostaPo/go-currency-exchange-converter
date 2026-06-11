# Currency Exchange API

REST API сервис для управления валютами, обменными курсами и конвертации денежных сумм.

## Запуск

### Клонирование репозитория

```bash
git clone https://github.com/KostaPo/go-currency-exchange-converter.git
cd go-currency-exchange-converter
```

### Локально

```bash
go run ./cmd/api -config=configs/config.local.yml
```

### На сервере

```bash
docker compose up -d
```

## API

### Валюты

| Метод | Endpoint           | Описание                |
| ----- | ------------------ | ----------------------- |
| GET   | `/currencies`      | Получить список валют   |
| GET   | `/currency/{code}` | Получить валюту по коду |
| POST  | `/currencies`      | Создать новую валюту    |

### Обменные курсы

| Метод | Endpoint               | Описание                        |
| ----- | ---------------------- | ------------------------------- |
| GET   | `/exchangeRates`       | Получить список обменных курсов |
| GET   | `/exchangeRate/{pair}` | Получить курс по валютной паре  |
| POST  | `/exchangeRates`       | Создать новый обменный курс     |
| PATCH | `/exchangeRate/{pair}` | Обновить существующий курс      |

### Конвертация

| Метод | Endpoint                                        | Описание                                      |
| ----- | ----------------------------------------------- | --------------------------------------------- |
| GET   | `/exchange?from={FROM}&to={TO}&amount={AMOUNT}` | Конвертировать сумму из одной валюты в другую |
