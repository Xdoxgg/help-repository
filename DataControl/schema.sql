-- Удаление таблиц, если они существуют
DROP TABLE IF EXISTS characters;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS team_to_character;

-- Создание таблицы characters
CREATE TABLE characters
(
    id                    SERIAL PRIMARY KEY,
    name                  VARCHAR(255) NOT NULL,
    role                  VARCHAR(255) NOT NULL,
    lore                  VARCHAR(800) NOT NULL, -- если будут траблы с мелким размером просто измени 800 на больше
    talents_build_emblems VARCHAR(255) NOT NULL
);

-- Создание таблицы teams
CREATE TABLE teams
(
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    orientation VARCHAR(255) NOT NULL
);

-- Промежуток между персонажами и командами
CREATE TABLE team_to_character
(
    id           SERIAL PRIMARY KEY,
    team_id      INT,
    character_id INT,
    FOREIGN KEY (team_id) REFERENCES teams (id),          -- Исправлено на правильную таблицу
    FOREIGN KEY (character_id) REFERENCES characters (id) -- Исправлено на правильную таблицу
);

-- Создание таблицы users
CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role     BOOLEAN -- 1-админ, 0-юзер
);