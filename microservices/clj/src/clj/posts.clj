(ns clj.posts
  (:require [hugsql.core :as hugsql]))

(def db
  {:classname "org.postgresql.Driver"
   :subprotocol "postgresql"
   :subname "//localhost/social"
   })

(hugsql/def-db-fns "sql/posts.sql")
