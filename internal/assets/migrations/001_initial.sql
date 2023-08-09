-- +migrate Up

create type requests_status_enum as enum ('success', 'in progress', 'failed', 'pending');

create table if not exists requests (
    id uuid primary key,
    status requests_status_enum not null,
    error text not null
);

create table if not exists feedbacks (
    course text not null,
    content text not null,

    unique(course, content)
);

create index feedbacks_course_idx on feedbacks(course);

-- +migrate Down

drop index if exists feedbacks_course_idx;
drop table if exists feedbacks;
drop table if exists requests;
drop type if exists requests_status_enum;