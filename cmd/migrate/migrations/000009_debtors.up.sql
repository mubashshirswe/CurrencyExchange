CREATE TABLE IF NOT EXISTS debtors (
    id bigserial PRIMARY KEY,
    received_amount bigint,
    received_currency  varchar(255),
    debted_amount bigint,
    debted_currency  varchar(255),
    user_id bigint,
    company_id bigint,
    details varchar(255) DEFAULT NULL,
    phone varchar(9) DEFAULT NULL,
    is_balance_effect int DEFAULT 0,
    type bigint,
    status  bigint,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);

ALTER TABLE debtors ADD CONSTRAINT fk_debtors_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE debtors ADD CONSTRAINT fk_debtors_company_id FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL;
