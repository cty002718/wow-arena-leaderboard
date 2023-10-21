CREATE TABLE IF NOT EXISTS token_price_logs (
    id uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    last_updated_time TIMESTAMP with time zone NOT NULL,
    price INT NOT NULL
);