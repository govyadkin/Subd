CREATE UNLOGGED TABLE "users" (
  "about" varchar NOT NULL,
  "email" varchar NOT NULL,
  "fullname" varchar NOT NULL,
  "nickname" varchar PRIMARY KEY
);

CREATE INDEX index_users_all ON users (lower(nickname));

CREATE UNLOGGED TABLE "forums" (
  "username" varchar NOT null,
  "posts" BIGINT DEFAULT 0,
  "threads" int DEFAULT 0,
  "slug" varchar PRIMARY KEY,
  "title" varchar NOT NULL,
  FOREIGN KEY ("username") REFERENCES "users" (nickname)
);

CREATE INDEX index_forums_slug ON forums (lower(slug));
CREATE INDEX index_users_fk ON forums (lower(author));
-- CREATE INDEX index_forum_all ON forums (slug, title, author, posts, threads);

CREATE UNLOGGED TABLE "threads" (
  "id" SERIAL PRIMARY KEY,
  "author" varchar NOT NULL,
  "created" timestamptz DEFAULT now(),
  "forum" varchar NOT NULL,
  "message" varchar NOT NULL,
  "slug" varchar,
  "title" varchar NOT NULL,
  "votes" int DEFAULT 0,
  FOREIGN KEY (author) REFERENCES "users" (nickname)
--   FOREIGN KEY (forum) REFERENCES "forums" (slug)
);

CREATE INDEX index_threads_slug ON threads (lower(slug));
CREATE INDEX index_thread_slug_hash ON threads USING hash (lower(slug));
CREATE INDEX index_thread_users_fk ON threads (lower(author));
-- CREATE INDEX index_thread_forum_fk ON threads (forum);
CREATE INDEX index_thread_forum_created ON threads (lower(forum), created);
-- CREATE INDEX index_thread_all ON threads (title, message, created, slug, author, forum, votes);

CREATE UNLOGGED TABLE "posts" (
  "author" varchar NOT NULL,
  "created" timestamp DEFAULT now(),
  "forum" varchar NOT NULL,
  "id" BIGSERIAL PRIMARY KEY,
  "is_edited" BOOL DEFAULT false,
  "message" varchar NOT NULL,
  "parent" BIGINT DEFAULT 0,
  "thread" int,
  "path" BIGINT[] DEFAULT ARRAY []::INTEGER[],
  
  FOREIGN KEY (author) REFERENCES "users" (nickname)
--   FOREIGN KEY (forum) REFERENCES "forums" (slug)
--   FOREIGN KEY (thread) REFERENCES "threads" (id)
--   FOREIGN KEY (parent) REFERENCES "posts" (id)
);

CREATE INDEX index_posts_thread ON posts (thread);
CREATE INDEX index_posts_author ON posts (thread, id, path);
CREATE INDEX index_posts_author ON posts (thread, path);
CREATE INDEX index_posts_author ON posts (thread, (path[1]));
CREATE INDEX index_posts_author ON posts (thread, parent);
-- CREATE INDEX index_post_thread_parent_path ON posts (thread, parent, path);
CREATE INDEX index_post_path1_path ON posts ((path[1]), path);
-- CREATE INDEX index_post_forum_fk ON posts (forum);
-- CREATE INDEX index_post_thread_created_id ON posts (thread, created, id);

CREATE UNLOGGED TABLE "votes" (
  "nickname" varchar NOT NULL,
  "voice" int,
  "thread" int,
  
   FOREIGN KEY (nickname) REFERENCES "users" (nickname),
--    FOREIGN KEY (thread) REFERENCES "threads" (id),
   UNIQUE (nickname, thread)
);

CREATE INDEX index_votes_thread_nick ON votes (thread, lower(nickname));

CREATE UNLOGGED TABLE forum_users
(
    author varchar REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
--     slug   varchar REFERENCES forums (slug) ON DELETE CASCADE NOT NULL,
    slug   varchar NOT NULL,
    UNIQUE (author, slug)
);
CREATE INDEX on forum_users (lower(slug));

CREATE OR REPLACE FUNCTION update_threads_count() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    UPDATE forums SET Threads=(Threads+1) WHERE LOWER(slug)=LOWER(NEW.forum);
    return NEW;
END
$update_users_forum$ LANGUAGE plpgsql;

CREATE TRIGGER add_thread
    BEFORE INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_threads_count();

CREATE OR REPLACE FUNCTION update_forum_users_by_insert_th_or_post()
RETURNS TRIGGER AS
$BODY$
BEGIN
    INSERT INTO forum_users values (NEW.author, NEW.forum)
    ON CONFLICT DO NOTHING;
    RETURN NULL;
END;
$BODY$ LANGUAGE plpgsql;

CREATE TRIGGER thread_insert_forum
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE update_forum_users_by_insert_th_or_post();

CREATE OR REPLACE FUNCTION update_path() RETURNS TRIGGER AS
$update_path$
DECLARE
    parent_path         BIGINT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.path := array_append(NEW.path, NEW.id);
    ELSE
        SELECT thread, path FROM posts
        WHERE thread = NEW.thread AND id = NEW.parent
        INTO first_parent_thread, parent_path;

        IF NOT FOUND THEN
            RAISE EXCEPTION 'bad parent thread';
        END IF;

        NEW.path := parent_path || NEW.id;
    END IF;
    UPDATE forums SET Posts=Posts + 1 WHERE lower(forums.slug) = lower(new.forum);
    RETURN NEW;
END
$update_path$ LANGUAGE plpgsql;

CREATE TRIGGER path_update
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_path();


CREATE OR REPLACE FUNCTION insert_votes() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    IF NEW.voice > 0 THEN
        UPDATE threads SET votes = (votes + 1) WHERE id = NEW.thread;
    ELSE
        UPDATE threads SET votes = (votes - 1) WHERE id = NEW.thread;
    END IF;
    return NEW;
END
$update_users_forum$ LANGUAGE plpgsql;

CREATE TRIGGER add_vote
    BEFORE INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE insert_votes();


CREATE OR REPLACE FUNCTION update_votes() RETURNS TRIGGER AS
$update_users_forum$
BEGIN
    IF NEW.voice = OLD.voice THEN
        RETURN NEW;
    END IF;
    IF NEW.voice > 0 THEN
        UPDATE threads SET votes = (votes + 2) WHERE id = NEW.thread;
    ELSE
        UPDATE threads SET votes = (votes - 2) WHERE id = NEW.thread;
    END IF;
    return NEW;
END
$update_users_forum$ LANGUAGE plpgsql;

CREATE TRIGGER edit_vote
    BEFORE UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE update_votes();
