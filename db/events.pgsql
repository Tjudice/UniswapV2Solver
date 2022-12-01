CREATE TABLE IF NOT EXISTS sync_event (
        primary_key VARCHAR(256) NOT NULL,

        block BIGINT NOT NULL,
        transaction BYTEA NOT NULL,
        transaction_index BIGINT NOT NULL,
        log_index BIGINT NOT NULL,
        address BYTEA NOT NULL,

        reserve0 TEXT NOT NULL,
        reserve1 TEXT NOT NULL,

        PRIMARY KEY (primary_key)
);

CREATE INDEX IF NOT EXISTS sync_event_block on sync_event (block);
CREATE INDEX IF NOT EXISTS sync_event_transaction on sync_event (transaction);
CREATE INDEX IF NOT EXISTS sync_event_address on sync_event (address);


CREATE TABLE IF NOT EXISTS pair_created_event (
        primary_key VARCHAR(256) NOT NULL,

        block BIGINT NOT NULL,
        transaction BYTEA NOT NULL,
        transaction_index BIGINT NOT NULL,
        log_index BIGINT NOT NULL,
        address BYTEA NOT NULL,

        token0 BYTEA NOT NULL,
        token1 BYTEA NOT NULL,
        pair BYTEA NOT NULL,
        pair_id TEXT NOT NULL,
        PRIMARY KEY (primary_key)
);

CREATE INDEX IF NOT EXISTS pair_created_event_block on pair_created_event (block);
CREATE INDEX IF NOT EXISTS pair_created_event_transaction on pair_created_event (transaction);
CREATE INDEX IF NOT EXISTS pair_created_event_token0 on pair_created_event (token0);
CREATE INDEX IF NOT EXISTS pair_created_event_token1 on pair_created_event (token1);
CREATE INDEX IF NOT EXISTS pair_created_event_pair on pair_created_event (pair);
CREATE INDEX IF NOT EXISTS pair_created_event_pair_id on pair_created_event (pair_id);


