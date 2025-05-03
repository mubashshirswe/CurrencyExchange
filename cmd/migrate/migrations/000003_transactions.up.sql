create table if not exists transactions(
    id bigserial PRIMARY KEY,
    amount bigint not null,
    service_fee bigint default null,
    from_currency_type_id bigint not null,
    to_currency_type_id bigint not null,
    sender_id bigint not null,
    from_city_id bigint not null,
    to_city_id bigint not null,
    receiver_name varchar(255) not null,
    receiver_phone varchar(9) not null,
    details varchar(255) not null,
    type bigint not null,
    received_time timestamp(0) with  default null,
    delivered_time timestamp(0) with  default null
)

