CREATE EXTENSION IF NOT EXISTS citext;

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
                     -- id serial PRIMARY KEY,
                      email citext UNIQUE NOT NULL,
                      fullname citext NOT NULL,
                      nickname citext PRIMARY KEY,
                      about text NOT NULL
);

-- Покрывающие индексы
--Get User
CREATE INDEX users_full ON users (nickname, email, fullname, about);
CREATE INDEX index_name_get_user ON users (nickname, email);
CREATE INDEX check_user ON users (nickname DESC );



CREATE UNLOGGED TABLE forums (
                        title citext  NOT NULL,
                        author citext references users(nickname),
                        slug citext PRIMARY KEY,
                        posts int,
                        threads int
);

--CREATE UNIQUE INDEX lower_slug_title ON forums(slug, title);

CREATE UNLOGGED TABLE forum_users(
                      nickname citext references users(nickname),
                      forum citext references forums(slug),
                     CONSTRAINT fk UNIQUE(nickname, forum)

);
CREATE INDEX fu_nick ON forum_users(nickname);
CREATE INDEX fu_for ON forum_users(forum);
--CREATE INDEX lower_forum_users ON forum_users(nickname, lower(forum));
--CREATE INDEX lower_forum ON forum_users USING hash(forum);
--CREATE INDEX lower_both ON forum_users(forum, nickname);
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

CREATE INDEX ON threads(id, forum);

CREATE OR REPLACE FUNCTION update_user_forum_thread() RETURNS TRIGGER AS
    $update_user_forum_thread$
BEGIN
    INSERT INTO forum_users(
        nickname,
        forum) VALUES (new.author, new.forum) ON CONFLICT DO NOTHING;
    UPDATE forums SET threads = threads + 1 WHERE lower(slug) = lower(new.forum);
    RETURN new;
END
    $update_user_forum_thread$ LANGUAGE plpgsql;

CREATE TRIGGER table_update_threads
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_user_forum_thread();



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


CREATE INDEX parent_tree_index
    ON posts ((path[1]), path DESC, id);

CREATE INDEX parent_tree_index2
    ON posts (id, (path[1]));


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
        forum) VALUES (new.author, new.forum) ON CONFLICT DO NOTHING;
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


