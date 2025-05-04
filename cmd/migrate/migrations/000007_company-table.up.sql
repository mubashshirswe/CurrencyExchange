create table if not exists companies(
    id bigserial PRIMARY KEY,
    name varchar(255) not null,
    details varchar(255) default null,
    password varchar(255) not null,
    created_at timestamp(0) with time zone not null default now()
)

