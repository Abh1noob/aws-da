Create table users(
    email varchar(255) primary key,
    password varchar(255) not null,
    username varchar(255) not null,
)

create table posts(
    post_id varchar(255) primary key,
    email varchar(255),
    image_url varchar(255),
    is_visible boolean,
    created_at timestamp,
)