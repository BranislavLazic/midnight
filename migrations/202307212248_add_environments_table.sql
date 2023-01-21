-- +goose Up
create table if not exists environments
(
    id   serial primary key,
    name varchar(255) not null
);

insert into environments (name)
values ('PROD'),
       ('DEV'),
       ('UAT');