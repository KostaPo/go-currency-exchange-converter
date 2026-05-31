PRAGMA foreign_keys=on;

-- =========================
-- CURRENCIES
-- =========================

CREATE TABLE IF NOT EXISTS Currencies (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Code TEXT NOT NULL,
    FullName TEXT NOT NULL,
    Sign TEXT,

    UNIQUE(Code)
);

CREATE INDEX IF NOT EXISTS currencies_id_index
ON Currencies (ID);

-- =========================
-- EXCHANGE RATES
-- (use business keys instead of IDs)
-- =========================

CREATE TABLE IF NOT EXISTS ExchangeRates (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,

    BaseCurrencyCode TEXT NOT NULL,
    TargetCurrencyCode TEXT NOT NULL,

    Rate REAL NOT NULL,

    CONSTRAINT fk_base_currency
        FOREIGN KEY (BaseCurrencyCode)
        REFERENCES Currencies (Code)
        ON DELETE CASCADE,

    CONSTRAINT fk_target_currency
        FOREIGN KEY (TargetCurrencyCode)
        REFERENCES Currencies (Code)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS curr_pair_unique_index
ON ExchangeRates (BaseCurrencyCode, TargetCurrencyCode);

CREATE INDEX IF NOT EXISTS exchange_rates_id_index
ON ExchangeRates (ID);

-- =========================
-- SEED: CURRENCIES
-- =========================

INSERT OR IGNORE INTO Currencies (Code, FullName, Sign)
VALUES ('RUB', 'Рубль', '₽');

INSERT OR IGNORE INTO Currencies (Code, FullName, Sign)
VALUES ('EUR', 'Евро', '€');

INSERT OR IGNORE INTO Currencies (Code, FullName, Sign)
VALUES ('USD', 'Доллар', '$');

INSERT OR IGNORE INTO Currencies (Code, FullName, Sign)
VALUES ('TRY', 'Лира', '₺');

INSERT OR IGNORE INTO Currencies (Code, FullName, Sign)
VALUES ('JPY', 'Иена', '¥');

INSERT OR IGNORE INTO Currencies (Code, FullName, Sign)
VALUES ('KZT', 'Тенге', '₸');

INSERT OR IGNORE INTO Currencies (Code, FullName, Sign)
VALUES ('GBP', 'Фунт', '£');

INSERT OR IGNORE INTO Currencies (Code, FullName, Sign)
VALUES ('ILS', 'Шекель', '₪');

INSERT OR IGNORE INTO Currencies (Code, FullName, Sign)
VALUES ('CNY', 'Юань', '¥');

INSERT OR IGNORE INTO Currencies (Code, FullName, Sign)
VALUES ('AZN', 'Манат', '₼');

-- =========================
-- SEED: EXCHANGE RATES
-- (business-key based)
-- =========================

INSERT OR IGNORE INTO ExchangeRates (BaseCurrencyCode, TargetCurrencyCode, Rate)
VALUES ('USD', 'RUB', 90.80);

INSERT OR IGNORE INTO ExchangeRates (BaseCurrencyCode, TargetCurrencyCode, Rate)
VALUES ('RUB', 'EUR', 0.0099);

INSERT OR IGNORE INTO ExchangeRates (BaseCurrencyCode, TargetCurrencyCode, Rate)
VALUES ('TRY', 'KZT', 16.53);

INSERT OR IGNORE INTO ExchangeRates (BaseCurrencyCode, TargetCurrencyCode, Rate)
VALUES ('JPY', 'GBP', 0.0055);

INSERT OR IGNORE INTO ExchangeRates (BaseCurrencyCode, TargetCurrencyCode, Rate)
VALUES ('USD', 'GBP', 0.78);