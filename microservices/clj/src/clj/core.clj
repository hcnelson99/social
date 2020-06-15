(ns clj.core
  (:require [compojure.core :refer :all]
            [compojure.route :as route]
            [selmer.parser :as s]
            [ring.util.response :refer [resource-response redirect]]
            [ring.util.anti-forgery :refer [anti-forgery-field]]
            [ring.middleware.defaults :refer [wrap-defaults site-defaults]]
            [clojure.string :as string]
            [clj.db :refer [db] :as db]
            [buddy.sign.jwt :as jwt]
            [buddy.hashers :as hashers]
            ))
(comment 
  (use 'clojure.repl)
  (db/register-user db {:username "Shoop" :password_hash "foo"})
  (db/get-user db {:username "bar"})
  )



; TODO: @security check that this can only serve things from within the
; resources directory (relative path names shouldn't allow one to request the
; database file, for example. Best thing is probably to also limit this to only
; serving from some static directory.
(defn serve-static [file]
  (resource-response file))

(defn index [] 
  (s/render-file "templates/index.html" {:links ["posts" "login" "logout" "register"]}))

(def mock-post-data 
  {:username "Henry Nelson"
   :comments [{:author "Bob" :date (new java.util.Date)
               :content ["This is my post"]}
              {:author "Alice" :date (new java.util.Date)
               :content ["This is my" "multi-paragraph" "post"]}]
   })

(defn posts [] 
  (s/render-file "templates/posts.html" 
                 {:username "Shoop"
                  :comments
                  (map 
                    #(let [text (:text %)]
                       (assoc % :content (string/split text #"\n\n"))
                       )
                    (db/get-all-posts db))
                  }))

(defn submit-comment [req] 
  (let [comment-text (get-in req [:form-params "comment"])]
    (db/insert-post db {:author "Default author" :text comment-text})
    (redirect "/posts" :see-other)))

(defn login-page [req]
  (let [error-message (get-in req [:flash :error-message])] 
    (s/render-file "templates/login.html" {:error error-message})))

(defn register-page [req]
  (let [error-message (get-in req [:flash :error-message])] 
    (s/render-file "templates/register.html" {:error error-message})))

(defn login [req]
  (let [params (:form-params req)
        username (params "username")
        password (params "password")
        user (db/get-user db {:username username})]
    (if (and user (hashers/check password (:password_hash user)))
      (redirect "/" :see-other)
      (assoc  (redirect "/login" :see-other) :flash {:error-message "Invalid username and/or password. Try again."}))
    ))

(comment
  (try
    (db/register-user db {:username "Shoop" :password_hash "test"})
    (catch org.postgresql.util.PSQLException _
      nil)
    ) )

(defn register-user [username password]
  (try
    (db/register-user db {:username username 
                       :password_hash (hashers/derive password {:alg :scrypt})})
    (catch org.postgresql.util.PSQLException _ nil)))

(defn register [req]
  (let [params (:form-params req)
        username (params "username")
        password (params "password")]
    (if (register-user username password)
        (redirect "/" :see-other)
        (assoc (redirect "/register" :see-other) :flash {:error-message "Username already taken. Try again."}))))

(defroutes handler
  (GET "/" [] (index))
  (GET "/posts" [] (posts))
  (POST "/posts" req (submit-comment req))
  (GET "/login" req (login-page req))
  (POST "/login" req (login req))
  (GET "/logout" [] "This is where you logout")
  (GET "/register" req (register-page req))
  (POST "/register" req (register req))
  (GET "/static/:file" [file] (serve-static file))
  (route/not-found "404 page not found"))

(def site
  (wrap-defaults handler 
    (assoc site-defaults :static false)))

(defn server-init []
  (s/add-tag! :csrf-field (fn [_ _] (anti-forgery-field))))
