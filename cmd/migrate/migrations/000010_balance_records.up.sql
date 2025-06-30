CREATE TABLE IF NOT EXISTS balance_records (
    id bigserial PRIMARY KEY,
    amount bigint,
    user_id bigint,
    balance_id bigint,
    company_id bigint,
    transaction_id bigint default null,
    debtor_id bigint default null,
    details varchar(255) DEFAULT NULL,
    currency varchar(255),
    type bigint,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);

ALTER TABLE balance_records ADD CONSTRAINT fk_balance_records_transaction_id FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE SET NULL;
ALTER TABLE balance_records ADD CONSTRAINT fk_balance_records_debtor_id FOREIGN KEY (debtor_id) REFERENCES debtors(id) ON DELETE SET NULL;
ALTER TABLE balance_records ADD CONSTRAINT fk_balance_records_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE balance_records ADD CONSTRAINT fk_balance_records_company_id FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL;
ALTER TABLE balance_records ADD CONSTRAINT fk_balance_records_balance_id FOREIGN KEY (balance_id) REFERENCES balances(id) ON DELETE SET NULL;