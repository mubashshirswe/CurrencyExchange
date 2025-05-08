CREATE TABLE IF NOT EXISTS balance_records (
    id bigserial PRIMARY KEY,
    amount bigint NOT NULL,
    user_id bigint,
    balance_id bigint,
    company_id bigint,
    details varchar(255) NOT NULL,
    currency_id bigint NOT NULL,
    type bigint NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);

ALTER TABLE balance_records ADD CONSTRAINT fk_balance_records_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE balance_records ADD CONSTRAINT fk_balance_records_company_id FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL;
ALTER TABLE balance_records ADD CONSTRAINT fk_balance_records_balance_id FOREIGN KEY (balance_id) REFERENCES balances(id) ON DELETE SET NULL;
ALTER TABLE balance_records ADD CONSTRAINT fk_balance_records_currency_id FOREIGN KEY (currency_id) REFERENCES currencies(id) ON DELETE SET NULL;
