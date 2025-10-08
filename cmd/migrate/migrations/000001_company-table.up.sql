create table if not exists companies(
    id bigserial PRIMARY KEY,
    name varchar(255),
    details varchar(255),
    password varchar(255),
    created_at timestamp(0) with time zone not null default now()
)

