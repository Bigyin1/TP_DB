CREATE EXTENSION
IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users
CASCADE;
DROP TABLE IF EXISTS forums
CASCADE;
DROP TABLE IF EXISTS threads
CASCADE;
DROP TABLE IF EXISTS posts
CASCADE;
DROP TABLE IF EXISTS votes
CASCADE;

DROP TABLE IF EXISTS UsersForum cascade;

DROP INDEX IF EXISTS idx_user_email;
DROP INDEX IF EXISTS idx_thread_slug;
DROP INDEX IF EXISTS idx_thread_forum;
DROP INDEX IF EXISTS idx_thread_author;
DROP INDEX IF EXISTS idx_post_author;
DROP INDEX IF EXISTS idx_post_thread;
DROP INDEX IF EXISTS idx_post_forum;
DROP INDEX IF EXISTS idx_user_nickname_email;
DROP INDEX IF EXISTS idx_post_branch;
DROP INDEX IF EXISTS idx_post_thread_parent;
DROP INDEX IF EXISTS idx_thread_created;
DROP INDEX IF EXISTS idx_usersForum;
DROP INDEX IF EXISTS idx_usersForum_forum;
DROP INDEX IF EXISTS idx_usersForum_name;




----------------USERS----------------
CREATE TABLE
IF NOT EXISTS users
(
  nickname CITEXT NOT NULL PRIMARY KEY COLLATE "C",
  fullname CITEXT                     NOT NULL,
  email    CITEXT                     NOT NULL UNIQUE,
  about    TEXT
);
CREATE INDEX IF NOT EXISTS idx_user_nickname_email ON Users USING btree  (nickname, email);
CREATE INDEX IF NOT EXISTS idx_user_email ON Users USING btree  (email);


----------------FORUMS----------------
CREATE TABLE
IF NOT EXISTS forums
(
  id      BIGSERIAL  PRIMARY KEY,
  posts   INT  NOT NULL  DEFAULT 0,
  slug    CITEXT  NOT NULL  UNIQUE,
  threads INT   NOT NULL  DEFAULT 0,
  title   varchar (100)  NOT NULL,
  owner CITEXT  NOT NULL  REFERENCES users
(nickname)
);

----------------THREADS----------------
CREATE TABLE
IF NOT EXISTS threads
(
  id        BIGSERIAL  PRIMARY KEY,
  author    CITEXT  NOT NULL  REFERENCES users
(nickname),
  created   TIMESTAMPTZ  NOT NULL DEFAULT (NOW () AT TIME ZONE 'UTC'),
  forum     CITEXT  NOT NULL  REFERENCES forums (slug),
  message   TEXT  NOT NULL,
  slug      CITEXT  DEFAULT NULL UNIQUE,
  title     CITEXT  NOT NULL,
  votes     INT  NOT NULL  DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_thread_created ON threads  USING btree (created);
CREATE INDEX IF NOT EXISTS idx_thread_author ON threads  USING btree (author);
CREATE INDEX IF NOT EXISTS idx_thread_slug ON threads  USING btree (slug);
CREATE INDEX IF NOT EXISTS idx_thread_forum ON threads  USING btree (forum);


----------------POSTS----------------
CREATE TABLE
IF NOT EXISTS posts
(
  id        BIGSERIAL                             PRIMARY KEY,
  author    CITEXT                      NOT NULL  REFERENCES users
(nickname),
  created   TIMESTAMPTZ                 NOT NULL  DEFAULT (NOW() AT TIME ZONE 'UTC'),
  forum     CITEXT                      NOT NULL  REFERENCES forums (slug),
  is_edited BOOLEAN                     NOT NULL  DEFAULT FALSE,
  message   CITEXT                      NOT NULL,
  parent    BIGINT DEFAULT 0            NOT NULL ,
  thread    BIGINT                      NOT NULL REFERENCES threads
(id),
  branch    BIGINT                      NOT NULL DEFAULT 0,
  path      BIGINT[]
);

CREATE INDEX IF NOT EXISTS idx_post_author ON posts USING btree (author COLLATE "C");
CREATE INDEX IF NOT EXISTS idx_post_thread ON posts USING btree (thread);
CREATE INDEX IF NOT EXISTS idx_post_forum ON posts USING btree (forum);
CREATE INDEX IF NOT EXISTS idx_post_branch ON posts USING btree (branch);
CREATE INDEX IF NOT EXISTS idx_post_thread_parent ON posts USING btree (thread, parent);
drop index if exists idx_post_path;
CREATE INDEX IF NOT EXISTS idx_post_path ON posts USING btree (path);

----------------VOTES----------------
CREATE TABLE
IF NOT EXISTS votes
(
  nickname  CITEXT                          NOT NULL          REFERENCES users (nickname),
  thread    BIGINT                          NOT NULL          REFERENCES threads (id),
  voice     SMALLINT                        NOT NULL,
  PRIMARY KEY (nickname, thread)
);

CREATE INDEX IF NOT EXISTS idx_vote_nickname_threadId ON votes USING btree (nickname, thread);

CREATE TABLE IF NOT EXISTS UsersForum (
  forum  CITEXT NOT NULL,
  userNickname  CITEXT NOT NULL,
  PRIMARY KEY (forum, userNickname)
);

CREATE INDEX IF NOT EXISTS idx_usersForum ON UsersForum USING btree (userNickname COLLATE "C", forum);
CREATE INDEX IF NOT EXISTS idx_usersForum_forum ON UsersForum USING btree (forum);
CREATE INDEX IF NOT EXISTS idx_usersForum_name ON UsersForum USING btree (userNickname COLLATE "C");

CREATE OR REPLACE FUNCTION insertPost
()
  RETURNS TRIGGER AS $$
DECLARE 
    parent_branch INT;
    parent_path BIGINT[];
BEGIN

	update forums f
	SET posts
	= posts + 1
  WHERE f.slug = NEW.forum;

	SELECT path, branch
	into parent_path
	, parent_branch
  FROM posts
  WHERE id = NEW.parent;
	IF NEW.parent != 0 
  THEN	
    NEW.branch = parent_branch;
ELSE
    NEW.branch = NEW.id;
END
IF;
  IF parent_path is null 
  THEN
    NEW.path =  NEW.path || cast
(0 as BIGINT) || cast
(NEW.id as BIGINT);
  ELSE
    NEW.path = parent_path || cast
(NEW.id as BIGINT);
END
IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insertPost
BEFORE
INSERT ON
Posts
FOR
EACH
ROW
EXECUTE PROCEDURE insertPost
();


CREATE OR REPLACE FUNCTION insertThread
()
  RETURNS TRIGGER AS $$
BEGIN

	update forums f
	SET threads
	= threads + 1
  WHERE f.slug = NEW.forum;

	RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insertThread
AFTER
INSERT ON
threads
FOR
EACH
ROW
EXECUTE PROCEDURE insertThread
();