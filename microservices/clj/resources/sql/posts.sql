-- :name get-all-posts :? :*
select author, date, text from comments;

-- :name insert-post :!
insert into comments (author, text) values (:author, :text);
