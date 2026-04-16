(ns wumpus.dodecahedron-spec
  (:require [speclj.core :refer :all]
            [wumpus.dodecahedron :refer :all]))

(describe "dodecahedron"
  (with-all cave-map (make-cave-map))

  (it "has 20 rooms"
    (should= 20 (count @cave-map)))

  (it "every room connects to exactly 3 others"
    (doseq [[room neighbors] @cave-map]
      (should= 3 (count neighbors))))

  (it "connections are symmetric"
    (doseq [[room neighbors] @cave-map
            neighbor neighbors]
      (should-contain room (get @cave-map neighbor))))

  (it "has no self-loops"
    (doseq [[room neighbors] @cave-map]
      (should-not-contain room neighbors)))

  (it "all rooms are reachable from room 1"
    (let [reachable (loop [visited #{1} frontier [1]]
                      (if (empty? frontier)
                        visited
                        (let [next-rooms (mapcat #(get @cave-map %) frontier)
                              new-rooms (remove visited next-rooms)]
                          (recur (into visited new-rooms)
                                 (vec new-rooms)))))]
      (should= 20 (count reachable))))

  (it "neighbors returns the 3 adjacent rooms"
    (should= (get @cave-map 1) (neighbors @cave-map 1)))

  (it "adjacent? returns true for connected rooms"
    (let [n (first (get @cave-map 1))]
      (should (adjacent? @cave-map 1 n))))

  (it "adjacent? returns false for unconnected rooms"
    (let [non-neighbor (first (clojure.set/difference
                                (set (range 1 21))
                                (conj (get @cave-map 1) 1)))]
      (should-not (adjacent? @cave-map 1 non-neighbor)))))
