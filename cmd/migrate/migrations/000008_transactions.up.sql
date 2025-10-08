CREATE TABLE IF NOT EXISTS transactions (
    id bigserial PRIMARY KEY,
    service_fee varchar,
    received_company_id bigint NULL,
    received_user_id bigint NULL,
    received_incomes jsonb DEFAULT '[]'::jsonb,
    delivered_outcomes jsonb DEFAULT '[]'::jsonb,
    delivered_company_id bigint NULL,
    delivered_user_id bigint NULL,
    phone varchar(255),
    details varchar(255),
    status bigint,
    type bigint,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);

ALTER TABLE transactions
    ADD CONSTRAINT fk_transactions_delivered_user_id FOREIGN KEY (delivered_user_id) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE transactions
    ADD CONSTRAINT fk_transactions_received_user_id FOREIGN KEY (received_user_id) REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE transactions
    ADD CONSTRAINT fk_transactions_received_company_id FOREIGN KEY (received_company_id) REFERENCES companies(id) ON DELETE SET NULL;

ALTER TABLE transactions
    ADD CONSTRAINT fk_transactions_delivered_company_id FOREIGN KEY (delivered_company_id) REFERENCES companies(id) ON DELETE SET NULL;
