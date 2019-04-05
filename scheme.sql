CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS forumusers CASCADE;


----------------USERS----------------
CREATE TABLE IF NOT EXISTS users (
  nickname CITEXT  COLLATE ucs_basic  NOT NULL  PRIMARY KEY,
  fullname CITEXT                     NOT NULL,
  email    CITEXT                     NOT NULL UNIQUE,
  about    TEXT
);


----------------FORUMS----------------
CREATE TABLE IF NOT EXISTS forums (
  id      BIGSERIAL  PRIMARY KEY,
  posts   INT  NOT NULL  DEFAULT 0,
  slug    CITEXT  NOT NULL  UNIQUE,
  threads INT   NOT NULL  DEFAULT 0,
  title   CITEXT  NOT NULL,
  owner CITEXT  NOT NULL  REFERENCES users(nickname)
);

----------------THREADS----------------
CREATE TABLE IF NOT EXISTS threads (
  id        BIGSERIAL  PRIMARY KEY,
  author    CITEXT  NOT NULL  REFERENCES users(nickname),
  created   TIMESTAMPTZ  NOT NULL DEFAULT transaction_timestamp(),
  forum     CITEXT  NOT NULL  REFERENCES forums(slug),
  message   TEXT  NOT NULL,
  slug      CITEXT  DEFAULT NULL UNIQUE,
  title     CITEXT  NOT NULL,
  votes     INT  NOT NULL  DEFAULT 0
);


----------------POSTS----------------
CREATE TABLE IF NOT EXISTS posts (
  id        BIGSERIAL                             PRIMARY KEY,
  author    CITEXT                      NOT NULL  REFERENCES users(nickname),
  created   TIMESTAMPTZ                 NOT NULL  DEFAULT transaction_timestamp(),
  forum     CITEXT                      NOT NULL  REFERENCES forums(slug),
  is_edited BOOLEAN                     NOT NULL  DEFAULT FALSE,
  message   CITEXT                      NOT NULL,
  parent    BIGINT DEFAULT 0            NOT NULL ,
  thread    BIGINT                      NOT NULL REFERENCES threads(id),
  branch    BIGINT                      NOT NULL DEFAULT 0,
  path      BIGINT[]
);

----------------VOTES----------------
CREATE TABLE IF NOT EXISTS votes (
  nickname  CITEXT                          NOT NULL          REFERENCES users(nickname),
  thread    BIGINT                          NOT NULL          REFERENCES threads (id),
  voice     SMALLINT                        NOT NULL,
  PRIMARY KEY (nickname, thread)
);

CREATE OR REPLACE FUNCTION insertPost()
  RETURNS TRIGGER AS $$
DECLARE 
    parent_branch INT;
    parent_path BIGINT[];
BEGIN
  
  update forums f
  SET posts = posts + 1
  WHERE f.slug = NEW.forum;

  SELECT path, branch into parent_path, parent_branch
  FROM posts
  WHERE id = NEW.parent;
  IF NEW.parent != 0 
  THEN	
    NEW.branch = parent_branch;
  ELSE
    NEW.branch = NEW.id;
  END IF;
  IF parent_path is null 
  THEN
    NEW.path =  NEW.path || cast(0 as BIGINT) || cast(NEW.id as BIGINT);
  ELSE
    NEW.path = parent_path || cast(NEW.id as BIGINT);
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insertPost
BEFORE INSERT ON Posts
FOR EACH ROW EXECUTE PROCEDURE insertPost();


CREATE OR REPLACE FUNCTION insertThread()
  RETURNS TRIGGER AS $$
BEGIN
  
  update forums f
  SET threads = threads + 1
  WHERE f.slug = NEW.forum;

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER insertThread
AFTER INSERT ON threads
FOR EACH ROW EXECUTE PROCEDURE insertThread();