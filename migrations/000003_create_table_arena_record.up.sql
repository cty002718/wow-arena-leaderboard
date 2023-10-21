CREATE TABLE IF NOT EXISTS servers (
  id INT PRIMARY KEY NOT NULL,
  name TEXT NOT NULL,
  slug TEXT NOT NULL,
  type TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS characters (
  id BIGINT PRIMARY KEY NOT NULL,
  server_id INT NOT NULL,
  name TEXT NOT NULL,
  faction TEXT NOT NULL,
  CONSTRAINT character_server_id_fk FOREIGN KEY (server_id) REFERENCES servers(id)
);

CREATE TABLE IF NOT EXISTS seasons (
  id INT PRIMARY KEY NOT NULL,
  started_at TIMESTAMP with time zone NOT NULL,
  ended_at TIMESTAMP with time zone NOT NULL
);

CREATE TABLE IF NOT EXISTS leaderboards (
  id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  season_id INT NOT NULL,
  created_at TIMESTAMP with time zone NOT NULL,
  bracket TEXT NOT NULL,
  CONSTRAINT leaderboard_season_id_fk FOREIGN KEY (season_id) REFERENCES seasons(id)
);

CREATE TABLE IF NOT EXISTS arena_records (
  id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
  leaderboard_id uuid NOT NULL,
  character_id BIGINT NOT NULL,
  rank INT NOT NULL,
  rating INT NOT NULL,
  won INT NOT NULL,
  lost INT NOT NULL,
  CONSTRAINT arena_record_leaderboard_id_fk FOREIGN KEY (leaderboard_id) REFERENCES leaderboards(id),
  CONSTRAINT arena_record_character_id_fk FOREIGN KEY (character_id) REFERENCES characters(id)
);
