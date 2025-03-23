create table if not exists consumer_1
(
    id serial primary key,
    message varchar(100)
);

create table if not exists consumer_2
(
    id serial primary key,
    message varchar(100)
);
