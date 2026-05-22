CREATE TABLE IF NOT EXISTS laptops (
    id            SERIAL PRIMARY KEY,
    brand         VARCHAR(50)    NOT NULL,
    model         VARCHAR(50)    NOT NULL,
    cpu           TEXT           NOT NULL,
    ram           INTEGER        NOT NULL,
    storage       INTEGER        NOT NULL,
    price         NUMERIC(12,2)  NOT NULL,
    year          INTEGER        NOT NULL,
    serial_number TEXT           NOT NULL
);
