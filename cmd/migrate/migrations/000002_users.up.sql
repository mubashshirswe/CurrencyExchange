create table if not exists users(
    id bigserial PRIMARY KEY,
    phone varchar(9) unique not null,
    role bigint not null,
    avatar varchar(255) default null,
    username varchar(255) not null,
    password bytea not null,
    company_id bigint default null REFERENCES companies(id),
    created_at timestamp(0) with time zone not null default now()
)