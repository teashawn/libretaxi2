-- drop index idx_created_at_utc_posts;
-- drop index users_geog_idx;
-- drop table posts;
-- drop table users;

create table if not exists users
(
    "userId" bigint not null,
    "menuId" int,
    "username" text,
    "firstName" text,
    "lastName" text,
	"lon" double precision,
    "lat" double precision,
	"geog" geography(POINT, 4326),
	"createdAtUtc" timestamp without time zone not null default (now() at time zone 'utc'),

    primary key ("userId")
);

create index if not exists users_geog_idx on users using gist(geog);

create table if not exists posts
(
    "postId" bigint generated by default as identity primary key,
    "userId" bigint not null references "users" ("userId"),
    "text" text,
	"lon" double precision,
    "lat" double precision,
	"geog" geography(POINT, 4326),
	"createdAtUtc" timestamp without time zone not null default (now() at time zone 'utc')
);

create index if not exists idx_created_at_utc_posts on "posts"("createdAtUtc");
