CREATE TABLE game_config (
	id bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
	game_id varchar(16) NOT NULL DEFAULT '',
	config TEXT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE game_config
    ADD INDEX idx_game_id (game_id), ALGORITHM=INPLACE, LOCK=NONE;

CREATE TABLE user_save (
	id bigint AUTO_INCREMENT NOT NULL PRIMARY KEY,
	username varchar(16) NOT NULL DEFAULT '',
	save TEXT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE user_save
    ADD INDEX idx_username (username), ALGORITHM=INPLACE, LOCK=NONE;