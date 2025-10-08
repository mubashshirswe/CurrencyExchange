CREATE TABLE IF NOT EXISTS exchanges (
    id bigserial PRIMARY KEY,
    received_money bigint,
    received_currency varchar(255),
    selled_money bigint,
    selled_currency varchar(255),
    user_id bigint REFERENCES users(id),
    status bigint,
    company_id bigint REFERENCES companies(id),
    details varchar(255),
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);