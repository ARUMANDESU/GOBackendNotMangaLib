
## Database (PostgreSQL)
```postgresql
create table useri(
    userid serial primary key,
    name varchar(50) not null
);

create table manga(
    mangaid serial primary key,
    name varchar not null,
    description varchar,
    author varchar(30),
    type varchar(20)
);

alter table manga add last_updated_time timestamp;

create table team(
    teamid serial primary key,
    name varchar(35) not null
);

create table teamuser(
    id serial primary key,
    teamid int,
    userid int,
    constraint fk_teamid
        foreign key (teamid)
        references team(teamid),
    constraint fk_userid
        foreign key (userid)
        references useri(userid)
);

create table mangateam(
    id serial primary key,
    teamid int,
    constraint fk_teamid
        foreign key (teamid)
        references team(teamid)
);
```