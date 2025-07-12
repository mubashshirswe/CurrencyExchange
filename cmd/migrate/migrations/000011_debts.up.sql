CREATE TABLE IF NOT EXISTS debts (
    id bigserial PRIMARY KEY,
    received_amount bigint,
    received_currency  varchar(255),
    debted_amount bigint,
    debted_currency  varchar(255),
    user_id bigint,
    company_id bigint,
    debtor_id bigint,
    details varchar(255) DEFAULT NULL,
    phone varchar(9) DEFAULT NULL,
    is_balance_effect int DEFAULT 0,
    type bigint,
    status  bigint,
    state bigint default 0,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);

ALTER TABLE debts ADD CONSTRAINT fk_debts_debtor_id FOREIGN KEY (debtor_id) REFERENCES debtors(id) ON DELETE SET NULL;
ALTER TABLE debts ADD CONSTRAINT fk_debts_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE debts ADD CONSTRAINT fk_debts_company_id FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL;
