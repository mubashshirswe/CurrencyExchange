CREATE TABLE IF NOT EXISTS user_sessions (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    device_id varchar(255) NOT NULL,
    fcm_token text NOT NULL,
    refresh_token text,
    platform varchar(64),
    app_version varchar(64),
    user_agent varchar(512),
    last_seen_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    CONSTRAINT uq_user_sessions_user_device UNIQUE (user_id, device_id)
);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions (user_id);
