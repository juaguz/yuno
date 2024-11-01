CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     user_id VARCHAR(255) UNIQUE NOT NULL,
     username VARCHAR(255) NOT NULL,
     email VARCHAR(255) UNIQUE NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cards (
                                     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                     created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     deleted_at TIMESTAMP, -- Puede ser NULL para soportar sql.NullTime
                                     card_holder VARCHAR(255) NOT NULL,
                                     user_id UUID NOT NULL,
                                     last_digits CHAR(4) NOT NULL, -- Últimos 4 dígitos de la tarjeta
                                     CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_cards_deleted_at ON cards (deleted_at);
