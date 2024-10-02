-- Table: gobloks.game

DROP TABLE IF EXISTS gobloks.game CASCADE;

CREATE TABLE IF NOT EXISTS gobloks.game
(
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    created timestamp without time zone NOT NULL,
    last_active timestamp without time zone NOT NULL,
    game_status int NOT NULL,
    player_count int NOT NULL,
    public boolean NOT NULL
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS gobloks.game OWNER to gobloks_user;

-- Table: gobloks.game_config

DROP TABLE IF EXISTS gobloks.game_config;

CREATE TABLE IF NOT EXISTS gobloks.game_config
(
    id bigint PRIMARY KEY,
    CONSTRAINT fk_game_id FOREIGN KEY (id)
        REFERENCES gobloks.game (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    players integer NOT NULL,
    block_degree integer NOT NULL,
    density float,
    turns BOOLEAN,
    time_seconds integer,
    time_bonus  integer,
    hints integer
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS gobloks.game_config OWNER to gobloks_user;