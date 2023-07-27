-- +goose Up
CREATE TABLE IF NOT EXISTS owners (
    id SERIAL PRIMARY KEY NOT NULL,
    telegram_id INT NOT NULL,
    telegram_name varchar (250) NOT NULL
);
CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY NOT NULL,
    token varchar (50) NOT NULL,
    owner_id INT NOT NULL,
    name varchar (250) NOT NULL,
    email varchar (250) NOT NULL,
    CONSTRAINT fk_owner FOREIGN KEY(owner_id) REFERENCES owners(id)
);
CREATE INDEX index_token ON companies (token);
CREATE TABLE IF NOT EXISTS chats (
    id SERIAL PRIMARY KEY NOT NULL,
    company_id INT NOT NULL,
    telegram_id INT NOT NULL,
    CONSTRAINT fk_company FOREIGN KEY(company_id) REFERENCES companies(id)
);
CREATE INDEX index_chat ON chats (company_id);
-- +goose Down
DROP INDEX IF EXISTS index_chat;
DROP TABLE IF EXISTS chats;
DROP INDEX IF EXISTS index_token;
DROP TABLE IF EXISTS companies;
DROP TABLE IF EXISTS owners;