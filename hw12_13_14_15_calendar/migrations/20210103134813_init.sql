-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE event (
   id serial PRIMARY KEY,
   title VARCHAR (100) not null,
   time timestamp without time zone NOT NULL,
   timezone smallint NOT NULL,
   duration smallint NOT NULL,
   description VARCHAR ( 255 ) NULL,
   user_id int NOT NULL,
   notification_time timestamp,
   is_notification_send boolean DEFAULT false
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
drop table event;