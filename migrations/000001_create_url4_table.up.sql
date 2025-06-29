CREATE TABLE url5 (
                     id bigserial primary key,
                     url varchar not null unique,
                     alias varchar not null unique,
                     dateAdd time not null default current_time
)