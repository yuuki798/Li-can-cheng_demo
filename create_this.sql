GRANT ALL PRIVILEGES ON todoDB.* TO 'user'@'%' IDENTIFIED BY 'password';


create table todos
(
    id      int auto_increment
        primary key,
    content text       not null,
    done    tinyint(1) null
);

create table users
(
    id       int auto_increment
        primary key,
    username varchar(255) not null,
    password varchar(255) not null,
    constraint username
        unique (username)
);


