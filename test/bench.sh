ab -n 50000 -kc 100 -T 'application/json' -H 'Accept-Encoding: gzip, deflate, br' -H 'Accept: application/json' -H 'DNT: 1' -H 'Origin: http://localhost:8080' -p post1.json http://localhost:8080/query

ab -n 50000 -kc 100 -T 'application/json' -H 'Accept-Encoding: gzip, deflate, br' -H 'Accept: application/json' -H 'DNT: 1' -H 'Origin: http://localhost:8080' -p post2.json http://localhost:8080/query
