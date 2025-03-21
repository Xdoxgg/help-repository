-- Inserts for the characters table
INSERT INTO characters (name, role, lore, talents_build_emblems)
VALUES ('Лукас', 'fighter', 'Леонин', 'лукас'),
       ('Суё', 'assassin', 'Японский воин', 'суё');

-- Inserts for the teams table
INSERT INTO teams (name, orientation)
VALUES ('TEAM 1', 'Урон и контроль');

-- Inserts for the team_to_character table
INSERT INTO team_to_character (team_id, character_id)
VALUES (1, 1),
       (1, 2);

-- Inserts for the users table
-- INSERT INTO users (username, password, role) VALUES