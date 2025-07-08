CREATE TABLE IF NOT EXISTS balance_records (
    id bigserial PRIMARY KEY,
    amount bigint,
    user_id bigint REFERENCES users(id),
    balance_id bigint REFERENCES balances(id),
    company_id bigint REFERENCES companies(id),
    exchange_id bigint default null REFERENCES exchanges(id),
    transaction_id bigint default null REFERENCES transactions(id),
    debtor_id bigint default null REFERENCES debtors(id),
    details varchar(255) DEFAULT NULL,
    currency varchar(255),
    type bigint,
    status bigint,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);