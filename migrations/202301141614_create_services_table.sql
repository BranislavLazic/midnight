-- +goose Up
create table if not exists services
(
    id                     bigserial primary key,
    name                   varchar(255),
    url                    text    not null,
    check_interval_seconds integer not null default 5
);