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
-- =========================

CREATE TABLE IF NOT EXISTS ExchangeRates (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,

    BaseCurrencyId INTEGER NOT NULL,
    TargetCurrencyId INTEGER NOT NULL,

    Rate REAL NOT NULL,

    CONSTRAINT fk_base_currency
        FOREIGN KEY (BaseCurrencyId)
        REFERENCES Currencies (ID)
        ON DELETE CASCADE,

    CONSTRAINT fk_target_currency
        FOREIGN KEY (TargetCurrencyId)
        REFERENCES Currencies (ID)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS curr_pair_unique_index
ON ExchangeRates (BaseCurrencyId, TargetCurrencyId);

CREATE INDEX IF NOT EXISTS exchange_rates_id_index
ON ExchangeRates (ID);

-- =========================
-- SEED: CURRENCIES
-- =========================

INSERT OR IGNORE INTO Currencies (Code, FullName, Sign) VALUES ('RUB', 'Рубль',   '₽');
INSERT OR IGNORE INTO Currencies (Code, FullName, Sign) VALUES ('EUR', 'Евро',    '€');
INSERT OR IGNORE INTO Currencies (Code, FullName, Sign) VALUES ('USD', 'Доллар',  '$');
INSERT OR IGNORE INTO Currencies (Code, FullName, Sign) VALUES ('TRY', 'Лира',    '₺');
INSERT OR IGNORE INTO Currencies (Code, FullName, Sign) VALUES ('JPY', 'Иена',    '¥');
INSERT OR IGNORE INTO Currencies (Code, FullName, Sign) VALUES ('KZT', 'Тенге',   '₸');
INSERT OR IGNORE INTO Currencies (Code, FullName, Sign) VALUES ('GBP', 'Фунт',    '£');
INSERT OR IGNORE INTO Currencies (Code, FullName, Sign) VALUES ('ILS', 'Шекель',  '₪');
INSERT OR IGNORE INTO Currencies (Code, FullName, Sign) VALUES ('CNY', 'Юань',    '¥');
INSERT OR IGNORE INTO Currencies (Code, FullName, Sign) VALUES ('AZN', 'Манат',   '₼');

-- =========================
-- SEED: EXCHANGE RATES
-- =========================

INSERT OR IGNORE INTO ExchangeRates (BaseCurrencyId, TargetCurrencyId, Rate)
VALUES (
    (SELECT ID FROM Currencies WHERE Code = 'USD'),
    (SELECT ID FROM Currencies WHERE Code = 'RUB'),
    90.80
);

INSERT OR IGNORE INTO ExchangeRates (BaseCurrencyId, TargetCurrencyId, Rate)
VALUES (
    (SELECT ID FROM Currencies WHERE Code = 'RUB'),
    (SELECT ID FROM Currencies WHERE Code = 'EUR'),
    0.0099
);

INSERT OR IGNORE INTO ExchangeRates (BaseCurrencyId, TargetCurrencyId, Rate)
VALUES (
    (SELECT ID FROM Currencies WHERE Code = 'TRY'),
    (SELECT ID FROM Currencies WHERE Code = 'KZT'),
    16.53
);

INSERT OR IGNORE INTO ExchangeRates (BaseCurrencyId, TargetCurrencyId, Rate)
VALUES (
    (SELECT ID FROM Currencies WHERE Code = 'JPY'),
    (SELECT ID FROM Currencies WHERE Code = 'GBP'),
    0.0055
);

INSERT OR IGNORE INTO ExchangeRates (BaseCurrencyId, TargetCurrencyId, Rate)
VALUES (
    (SELECT ID FROM Currencies WHERE Code = 'USD'),
    (SELECT ID FROM Currencies WHERE Code = 'GBP'),
    0.78
);