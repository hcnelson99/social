-- :name get-all-posts :? :*
select id, author, date, text from comments;

-- :name get-post :? :1
select id, author, date, text from comments where id = :id;

-- :name insert-post :!
insert into comments (author, text) values (:author, :text);

-- :name register-user :<!
insert into users(username, password_hash) values (:username, :password_hash) returning id, session_generation;

-- :name get-user :? :1
select username, session_generation, password_hash from users where username = :username;
