DROP TABLE IF EXISTS users;
CREATE TABLE users (
                       id bigserial primary key,
                       username varchar(40) unique,
                       password varchar(40),
                       admin bool,
                       created_at timestamp default now() not null,
                       updated_at timestamp default now() not null
);

insert into users (id, username, password, admin) VALUES (1, 'james', 'password', true);
insert into users (id, username, password, admin) VALUES (2, 'david0122', 'password', false);
insert into users (id, username, password, admin) VALUES (3, 'nickC121', 'password', false);
insert into users (id, username, password, admin) VALUES (4, 'Sarah0123', 'password', false);
insert into users (id, username, password, admin) VALUES (5, 'Nathan39024', 'password', false);
insert into users (id, username, password, admin) VALUES (6, 'Mary43243', 'password', false);

select * from users;