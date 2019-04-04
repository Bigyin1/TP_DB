INSERT INTO users(nickname, fullname, about, email)
	VALUES ( '1', '2', '4', '4' );

INSERT INTO users(nickname, fullname, about, email)
	VALUES ( '2', '3', '4', '5' );

INSERT INTO users(nickname, fullname, about, email)
	VALUES ( '3', '4', '4', '6' );

INSERT INTO users(nickname, fullname, about, email)
	VALUES ( '4', '5', '4', '7' );

INSERT INTO users(nickname, fullname, about, email)
	VALUES ( '5', '6', '4', '8' );

INSERT INTO users(nickname, fullname, about, email)
	VALUES ( '6', '7', '4', '9' );

INSERT INTO users(nickname, fullname, about, email)
	VALUES ( '7', '8', '4', '10' );

INSERT INTO users(nickname, fullname, about, email)
	VALUES ( '8', '9', '4', '11' );


----FORUMS------

INSERT INTO forums(owner, title, slug)
	VALUES ('1', 'cat', 'q');

INSERT INTO forums(owner, title, slug)
	VALUES ('2', 'cat', 'w');

INSERT INTO forums(owner, title, slug)
	VALUES ('3', 'cat', 'e');

INSERT INTO forums(owner, title, slug)
	VALUES ('4', 'cat', 'r');


----Threads-----

INSERT INTO threads(author, forum, message, slug, title, created)
	VALUES ('5', 'q', 'aq', 'q', 'aqq', '2009-06-04 19:25:21');

INSERT INTO threads(author, forum, message, slug, title, created)
	VALUES ('6', 'e', 'aq', 'w', 'aqq', '2009-06-04 19:25:21');

INSERT INTO threads(author, forum, message, slug, title, created)
	VALUES ('7', 'r', 'aq', 'e', 'aqq', '2009-06-04 19:25:21');

INSERT INTO threads(author, forum, message, slug, title, created)
	VALUES ('8', 'w', 'aq', 'r', 'aqq', '2009-06-04 19:25:21');

INSERT INTO threads(author, forum, message, slug, title, created)
	VALUES ('2', 'q', 'aq', 't', 'aqq', '2009-06-04 19:25:21');

INSERT INTO threads(author, forum, message, slug, title, created)
	VALUES ('1', 'w', 'aq', 'y', 'aqq', '2009-06-04 19:25:21');


----Posts------

INSERT INTO posts(author, forum, message, parent, thread)
	VALUES ('1', 'r', 'hi fillip', '0', '1');

INSERT INTO posts(author, forum, message, parent, thread)
	VALUES ('2', 'e', 'hi fillip', '0', '3');

INSERT INTO posts(author, forum, message, parent, thread)
	VALUES ('1', 'r', 'hi fillip', '0', '1');

INSERT INTO posts(author, forum, message, parent, thread)
	VALUES ('4', 'q', 'hi fillip', '0', '2');

INSERT INTO posts(author, forum, message, parent, thread)
	VALUES ('3', 'w', 'hi fillip', '0', '4');

INSERT INTO posts(author, forum, message, parent, thread)
	VALUES ('8', 'q', 'hi fillip', '0', '2');

INSERT INTO posts(author, forum, message, parent, thread)
	VALUES ('7', 'w', 'hi fillip', '0', '5');