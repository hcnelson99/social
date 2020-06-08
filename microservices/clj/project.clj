(defproject clj "0.1.0-SNAPSHOT"
  :description "FIXME: write description"
  :url "http://example.com/FIXME"
  :license {:name "EPL-2.0 OR GPL-2.0-or-later WITH Classpath-exception-2.0"
            :url "https://www.eclipse.org/legal/epl-2.0/"}
  :dependencies [[org.clojure/clojure "1.10.1"]
                 [ring/ring-defaults "0.3.2"]
                 [compojure "1.6.1"]
                 [selmer "1.12.27"]
                 [com.layerware/hugsql "0.5.1"]
                 [org.postgresql/postgresql "42.2.13"]
                 [org.clojure/tools.logging "1.1.0"]]
  :main ^:skip-aot clj.core
  :target-path "target/%s"
  :plugins [[lein-ring "0.12.5"]]
  :ring {:handler clj.core/site
         :init clj.core/server-init
         :nrepl {:start? true}}
  :profiles {:uberjar {:aot :all}
             :dev {:dependencies 
                   [[javax.servlet/servlet-api "2.5"]]}})
