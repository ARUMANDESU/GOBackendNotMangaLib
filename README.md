
## Database (PostgreSQL)
```postgresql
create table manga
(
    mangaid           serial
        primary key,
    name              varchar not null,
    description       varchar,
    author            varchar(30),
    type              varchar,
    last_updated_time timestamp,
    status            varchar,
    rating            double precision
);
create table mangateam
(
    id     serial
        primary key,
    teamid integer
        constraint fk_teamid
            references team
);
create table team
(
    teamid serial
        primary key,
    name   varchar(35) not null
);

create table teamuser
(
    id     serial
        primary key,
    teamid integer
        constraint fk_teamid
            references team,
    userid integer
        constraint fk_userid
            references useri
);
create table useri
(
    userid          serial
        primary key,
    name            varchar(50) not null,
    email           varchar
        constraint unique_email
            unique,
    hashed_password char(60)    not null,
    role            varchar
);
```