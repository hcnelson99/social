(ns clj.core
  (:require [compojure.core :refer :all]
            [compojure.route :as route]
            [selmer.parser :as s]
            [ring.util.response :refer [resource-response]]

            [ring.middleware.defaults :refer [wrap-defaults site-defaults]]))


(defn serve-static [file]
  (resource-response file))

(defn index [] 
  (s/render-file "templates/index.html" {:links ["posts" "login" "logout" "register"]}))

(def mock-post-data 
  {:username "Henry"
   :comments [{:author "Bob" :date (new java.util.Date)
               :content ["This is my post"]}
              {:author "Alice" :date (new java.util.Date)
               :content ["This is my" "multi-paragraph" "post"]}]
   })

(defn posts [] 
  (s/render-file "templates/posts.html" 
                 mock-post-data))

(defroutes handler
  (GET "/" [] (index))
  (GET "/posts" [] (posts))
  (GET "/login" [] "This is where you login")
  (GET "/logout" [] "This is where you logout")
  (GET "/register" [] "This is where you register")
  (GET "/static/:file" [file] (serve-static file))
  (route/not-found "404 page not found"))

(def site
  (wrap-defaults handler 
    (assoc site-defaults :static false)))

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

