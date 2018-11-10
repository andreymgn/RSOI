CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    uid UUID PRIMARY KEY,
    username VARCHAR(30) NOT NULL,
    password_hash CHAR(60) NOT NULL
);
