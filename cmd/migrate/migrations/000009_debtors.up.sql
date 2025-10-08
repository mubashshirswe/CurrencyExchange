CREATE TABLE IF NOT EXISTS debtors (
    id bigserial PRIMARY KEY,
    balance bigint,
    currency  varchar(255),
    user_id bigint,
    company_id bigint,
    phone varchar(255),
    full_name varchar(255),
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);

ALTER TABLE debtors ADD CONSTRAINT fk_debtors_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE debtors ADD CONSTRAINT fk_debtors_company_id FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL;
