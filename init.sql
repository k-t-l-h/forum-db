CREATE EXTENSION IF NOT EXISTS citext;

ALTER SYSTEM SET max_connections = '200';
ALTER SYSTEM SET shared_buffers = '256B';
ALTER SYSTEM SET effective_cache_size = '768MB';
ALTER SYSTEM SET checkpoint_completion_target = '0.9';
ALTER SYSTEM SET random_page_cost = '1.1';
ALTER SYSTEM SET effective_io_concurrency = '400';
ALTER SYSTEM SET seq_page_cost = '1.1';
ALTER SYSTEM SET random_page_cost = '1.1';

ALTER SYSTEM SET max_worker_processes = '4';
ALTER SYSTEM SET max_parallel_workers_per_gather = '2';
ALTER SYSTEM SET max_parallel_workers = '4';
ALTER SYSTEM SET max_parallel_maintenance_workers = '2';


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
CREATE INDEX check_name ON users(nickname);
CREATE INDEX index_name ON users(email, fullname, nickname, about, lower(nickname), lower(email));
CREATE INDEX index_name_get_user ON users(email, fullname, nickname, about, lower(nickname));


CREATE UNLOGGED TABLE forums (
                        title citext  NOT NULL,
                        author citext references users(nickname),
                        slug citext PRIMARY KEY,
                        posts int,
                        threads int
);


CREATE INDEX  forum_slug ON forums(slug);
CREATE INDEX forum_exist ON forums(slug, lower(slug));
CREATE INDEX lower_slug ON forums(lower(slug));
CLUSTER forums USING lower_slug;

CREATE UNLOGGED TABLE forum_users(
                      nickname citext references users(nickname),
                      forum citext references forums(slug)
);

CREATE INDEX lower_forum_users ON forum_users(nickname, lower(forum));
CREATE INDEX lower_forum ON forum_users(lower(forum));

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
    UPDATE forums SET threads = threads + 1 WHERE slug = new.forum;
    RETURN new;
END
    $update_user_forum_thread$ LANGUAGE plpgsql;

CREATE TRIGGER table_update_threads
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_user_forum_thread();


CREATE INDEX threads_all ON threads(id, author, message, title, created_at, forum, slug, votes);
CREATE INDEX lower_thread_name ON threads(lower(slug));
CLUSTER threads USING lower_thread_name;

---CREATE POST INDEX
CREATE INDEX threads_id_forum ON threads(id, forum);
CREATE INDEX threads_id_forum_lower ON threads(id, forum, lower(slug));

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


CREATE INDEX posts_path_first on posts((path[1]));
CREATE INDEX posts_thread_id on posts(id, thread);
CREATE INDEX posts_trigger_update on posts(id, path);


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
    UPDATE forums SET posts = posts + 1 WHERE slug = new.forum;
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
                         UNIQUE(author, thread)
);

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


