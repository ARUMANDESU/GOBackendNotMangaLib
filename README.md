
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
    rating            double precision,
    mangaImg varchar default ''
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
    userid          serial primary key,
    name            varchar(50) not null,
    email           varchar
        constraint unique_email
            unique,
    hashed_password char(60)    not null,
    role            varchar
);

create table chapter
(
    chapterid      serial
        primary key,
    chapter_number numeric(5, 1),
    volume_number  numeric(4),
    title          varchar default ''::character varying,
    images         character varying[]
);

create table manga_chapter(
      chapterid integer
          constraint fk_chapterid
              references chapter(chapterid),
      mangaid integer
          constraint fk_mangaid
              references manga(mangaid)

);
```