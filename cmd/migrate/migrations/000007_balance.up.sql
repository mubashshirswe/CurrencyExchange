CREATE TABLE IF NOT EXISTS balances (
    id bigserial PRIMARY KEY,
    balance bigint DEFAULT NULL,
    user_id bigint,
    in_out_lay bigint DEFAULT NULL,
    out_in_lay bigint DEFAULT NULL,
    currency_id bigint,
    company_id bigint,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT now()
);

ALTER TABLE balances ADD CONSTRAINT fk_balances_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE balances ADD CONSTRAINT fk_balances_currency_id FOREIGN KEY (currency_id) REFERENCES currencies(id) ON DELETE SET NULL;
ALTER TABLE balances ADD CONSTRAINT fk_balances_company_id FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL;
