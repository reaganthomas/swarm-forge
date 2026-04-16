(ns wumpus.dodecahedron-spec
  (:require [speclj.core :refer :all]
            [wumpus.dodecahedron :refer :all]))

(describe "dodecahedron"
  (it "has 20 rooms"
    (should= 20 (count (make-cave-map)))))
