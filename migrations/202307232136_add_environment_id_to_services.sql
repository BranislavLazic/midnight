-- +goose Up
alter table services
    add column environment_id bigint;

alter table services
    add constraint environment_id_fkey foreign key (environment_id) references environments (id) on update cascade on delete restrict;
