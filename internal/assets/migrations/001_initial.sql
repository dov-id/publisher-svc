-- +migrate Up

create table if not exists feedbacks (
    course text not null,
    content text not null,

    unique(course, content)
);

create index feedbacks_course_idx on feedbacks(course);

-- +migrate Down

drop index if exists feedbacks_course_idx;
drop table if exists feedbacks;