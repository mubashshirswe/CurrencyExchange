CREATE TABLE IF NOT EXISTS employees (
    id bigserial PRIMARY KEY,
    phone varchar(9) UNIQUE NOT NULL,
    role bigint NOT NULL,
    avatar varchar(255) DEFAULT NULL,
    username varchar(255) UNIQUE NOT NULL,
    password bytea NOT NULL,
    company_id bigint,
    created_at timestamp(0) with time zone NOT NULL DEFAULT now()
);

ALTER TABLE employees
    ADD CONSTRAINT fk_employees_company_id FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE SET NULL;
