CREATE TABLE IF NOT EXISTS debtors (
    id bigserial PRIMARY KEY,
    amount bigint NOT NULL,
    serial_no varchar(255) unique,
    user_id bigint,
    balance_id bigint,
    company_id bigint,
    details varchar(255) NOT NULL,
    debtors_name varchar(255) NOT NULL,
    debtors_phone varchar(9) NOT NULL,
    currency_id bigint NOT NULL,
    is_balance_effect int NOT NULL,
    currency_type varchar(255) NOT NULL,
    type bigint NOT NULL,
    status  bigint DEFAULT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);

ALTER TABLE debtors ADD CONSTRAINT fk_debtors_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE debtors ADD CONSTRAINT fk_debtors_company_id FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL;
ALTER TABLE debtors ADD CONSTRAINT fk_debtors_balance_id FOREIGN KEY (balance_id) REFERENCES balances(id) ON DELETE SET NULL;
ALTER TABLE debtors ADD CONSTRAINT fk_debtors_currency_id FOREIGN KEY (currency_id) REFERENCES currencies(id) ON DELETE SET NULL;
