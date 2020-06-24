CREATE EXTENSION IF NOT EXISTS citext;

ALTER SYSTEM SET max_connections = '200';
ALTER SYSTEM SET shared_buffers = '512B';
ALTER SYSTEM SET effective_cache_size = '768MB';
ALTER SYSTEM SET checkpoint_completion_target = '0.9';
ALTER SYSTEM SET random_page_cost = '1.1';
ALTER SYSTEM SET effective_io_concurrency = '200';
ALTER SYSTEM SET wal_buffers = '6912kB';
ALTER SYSTEM SET default_statistics_target = '100';
ALTER SYSTEM SET seq_page_cost = '1.1';
/*CREATE DATABASE forum
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'Russian_Russia.1251'
    LC_CTYPE = 'Russian_Russia.1251'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;
 */

CREATE UNLOGGED TABLE users(
                      id serial PRIMARY KEY,
                      email citext UNIQUE NOT NULL,
                      fullname citext NOT NULL,
                      nickname citext UNIQUE NOT NULL ,
                      about text
);

-- Покрывающие индексы
--Get User
CREATE INDEX check_lower_name ON users USING hash(lower(nickname));
--User Conflict
CREATE INDEX index_name_get_user ON users (lower(nickname), lower(email));
CLUSTER users USING check_lower_name;


CREATE UNLOGGED TABLE forums (
                        title citext  NOT NULL,
                        author citext references users(nickname),
                        slug citext PRIMARY KEY,
                        posts int,
                        threads int
);

CREATE UNIQUE INDEX lower_slug ON forums(lower(slug));
CREATE UNIQUE INDEX lower_slug_title ON forums(lower(slug), title);

CREATE UNLOGGED TABLE forum_users(
                      nickname citext references users(nickname),
                      forum citext references forums(slug) --,
                     -- CONSTRAINT fk UNIQUE(nickname, forum)

);

--CREATE INDEX lower_forum_users ON forum_users(nickname, lower(forum));
CREATE INDEX lower_forum ON forum_users USING hash(forum);
CREATE INDEX lower_both ON forum_users(lower(forum), lower(nickname));
--CLUSTER forum_users USING lower_forum;

CREATE UNLOGGED TABLE threads (
                         id serial PRIMARY KEY,
                         author citext references users(nickname),
                         message citext NOT NULL,
                         title citext NOT NULL,

                         created_at timestamp with time zone,
                         forum citext references forums(slug),
                         slug citext,
                         votes int
);


CREATE OR REPLACE FUNCTION update_user_forum_thread() RETURNS TRIGGER AS
    $update_user_forum_thread$
BEGIN
    INSERT INTO forum_users(
        nickname,
        forum) VALUES (new.author, new.forum);
    UPDATE forums SET threads = threads + 1 WHERE lower(slug) = lower(new.forum);
    RETURN new;
END
    $update_user_forum_thread$ LANGUAGE plpgsql;

CREATE TRIGGER table_update_threads
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_user_forum_thread();


CREATE INDEX forum_date ON threads(lower(forum), created_at);
--CREATE INDEX lower_thread_name_id ON threads(lower(slug));
CREATE INDEX id_check ON threads(id);
--CLUSTER threads USING lower_thread_name_id;


CREATE UNLOGGED TABLE posts (
                       id serial PRIMARY KEY ,
                       author citext references users(nickname),
                       post citext NOT NULL,

                       created_at timestamp with time zone,
                       forum citext references forums(slug),
                       isEdited bool,
                       parent int,
                       thread int references threads(id),
                       path  INTEGER[]
);


--CREATE INDEX posts_hash_forum on posts USING hash(forum);
--CREATE INDEX posts_thread ON posts(thread);

--CREATE INDEX posts_id_thread ON posts(parent, thread, id);
--CREATE INDEX posts_id_thread ON posts(thread) WHERE parent = 0;;
--CREATE INDEX thread_path_null ON posts((path[1]));
--CREATE INDEX posts_id_thread_a ON posts(thread, id, path) WHERE id > 5000;
--CREATE INDEX posts_id_thread_b ON posts(thread, id, path) WHERE id <= 5000;
--CREATE INDEX posts_thread_path ON posts(thread, id);

--CREATE INDEX parent_thread_check ON posts (thread, parent) WHERE parent = 0;
--CREATE INDEX posts_id_thread_parent ON posts (id, thread, parent);
CREATE INDEX thread_path_null ON posts((path[1]) DESC, path,  id);

CREATE INDEX thread_path_tree ON posts( path,  id);


CREATE OR REPLACE FUNCTION update_path() RETURNS TRIGGER AS
$update_path$
DECLARE
    parent_path  INTEGER[];
    parent_thread int;
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.path := array_append(new.path, new.id);
    ELSE
        SELECT thread FROM posts WHERE id = new.parent INTO parent_thread;
        IF NOT FOUND OR parent_thread != NEW.thread THEN
            RAISE EXCEPTION 'this is an exception' USING ERRCODE = '22000';
        end if;

        SELECT path FROM posts WHERE id = new.parent INTO parent_path;
        NEW.path := parent_path || new.id;
    END IF;
    RETURN new;
END
$update_path$ LANGUAGE plpgsql;


CREATE TRIGGER path_update_trigger
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_path();



CREATE OR REPLACE FUNCTION update_user_forum() RETURNS TRIGGER AS
$update_user_forum$
BEGIN
    UPDATE forums SET posts = posts + 1 WHERE lower(slug) = lower(new.forum);
    INSERT INTO forum_users(
        nickname,
        forum) VALUES (new.author, new.forum);
    RETURN new;
END
$update_user_forum$ LANGUAGE plpgsql;

CREATE TRIGGER table_update
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_user_forum();


CREATE UNLOGGED TABLE votes (

                         author citext references users(nickname),
                         vote int,
                         thread int references threads(id), --slug thread
                         CONSTRAINT checks UNIQUE(author, thread)
);
CLUSTER votes USING checks;

CREATE OR REPLACE FUNCTION add_votes() RETURNS TRIGGER AS
$add_votes$
BEGIN
    UPDATE threads SET votes=(votes+NEW.vote) WHERE id = NEW.thread;
    return NEW;
end
$add_votes$ LANGUAGE plpgsql;

CREATE TRIGGER add_vote
    BEFORE INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE add_votes();


CREATE OR REPLACE FUNCTION update_votes() RETURNS TRIGGER AS
$update_votes$
BEGIN
    UPDATE threads SET votes=votes - old.vote + new.vote WHERE id = new.thread;
    return NEW;
end
$update_votes$ LANGUAGE plpgsql;


CREATE TRIGGER update_vote
    BEFORE UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_votes();


