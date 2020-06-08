(ns clj.core
  (:require [compojure.core :refer :all]
            [compojure.route :as route]
            [selmer.parser :as s]
            [ring.util.response :refer [resource-response redirect]]
            [ring.util.anti-forgery :refer [anti-forgery-field]]
            [ring.middleware.defaults :refer [wrap-defaults site-defaults]]
            [clojure.string :as string]
            [clj.posts :refer [db] :as posts]
            ))

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
                    (posts/get-all-posts db))
                  }))

(defn submit-comment [req] 
  (let [comment-text (get-in req [:form-params "comment"])]
    (posts/insert-post db {:author "Default author" :text comment-text})
    (posts)))

(defroutes handler
  (GET "/" [] (index))
  (GET "/posts" [] (posts))
  (POST "/posts" req (submit-comment req))
  (GET "/login" [] "This is where you login")
  (GET "/logout" [] "This is where you logout")
  (GET "/register" [] "This is where you register")
  (GET "/static/:file" [file] (serve-static file))
  (route/not-found "404 page not found"))

(def site
  (wrap-defaults handler 
    (assoc site-defaults :static false)))

(defn server-init []
  (s/add-tag! :csrf-field (fn [_ _] (anti-forgery-field))))

(comment
  (use 'clojure.repl)
  (require '[ring.util.request :as request]
           '[ring.util.response :as response]
           '[ring.util.codec :as codec]
           )
  (response/resource-response "style.css" {:root "public"})
  (doc response/resource-response)
  (get-in site-defaults [:static :resources] false)
  (site {:request-method :get :uri "/static/style.css"})
  )

