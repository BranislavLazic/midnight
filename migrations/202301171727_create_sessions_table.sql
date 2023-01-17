-- +goose Up
create table if not exists sessions
(
    k varchar(100) primary key not null default '',
    v bytea                    not null,
    e bigint                   not null default '0'
);