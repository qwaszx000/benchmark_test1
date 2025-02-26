The idea of this project is to compare performance of languages and their frameworks

###### Goals of this project:
1. To get handy comparison results, performed by myself
2. To satisfy my curiosity about languages and performance
3. To practice writing performant and high-load software, measuring it, designing high-load systems

###### There will be 1 test(maybe i'll add second later for chosen technologies):
1. Plain text constant response without any logic behind it
2. More real usage: db requests, maybe caching with memcached/redis

---

###### Metrics:
1. CPU usage
2. Memory usage
3. RPM(successful requests per minute)
4. My subjective opinion

###### Testing methodology:
I will be using one or more of these ways to perform tests
1. Docker env in parallel
2. Docker env in serial order
3. Plain running on my pc in serial order
4. VMs in serial order
5. VMs in parallel
6. Some other way???

###### Measuring methodology:
Again, i'll choice one or more of these
1. [Prometheus](https://prometheus.io) + [Grafana](https://grafana.com/)
2. My own tool written on one lang or on different languages + plotting with python matplotlib

---

###### Versions:
 - `rustc --version` -> `rustc 1.85.0 (4d91de4e4 2025-02-17)`
 - `go version` -> `go version go1.24.0 linux/amd64`

---

### Results:

#### Subjective
Go //Looks like best language now for performance/ease of use
1. Gin -- fast and simple to use(at least for simple things) (12M release file)
2. Gnet -- feels like writing using simple tcp sockets, takes time and adds complexity (6.4M)
3. Stdlib net/http -- fast to use, but i like Gin/Fiber more (7.9M)
4. Fiber -- simple to use, just like Gin (9.3M)
5. FastHttp -- something between Gnet and Gin, i think (7.6M)

Rust //Complex thing, but i would use it if Go didn't exist, probably 
1. Actix -- i would give it second place after rocket (release file if 5.8M)
2. Rocket -- nice and simple, i like it(based on this simple project, unsure about something bigger) (4.9M)
3. Faf -- like gnet on Go with rust "things" -- "newborn" thread 0 or thread 1 overflows it's stack even when handler just returns 0 or when using example code
4. Ntex -- same as actix(but release file is 4.3M)

JS
1. Express.js -- fast to make simple project
2. Hyperexpress -- also easy to use

Python
1. Flask -- quite easy
2. FastAPI -- feels on par with flask, maybe docs are not so good

PHP
1. Symfony -- feels a little lighter than Laravel, but mostly same verdict
2. Laravel -- big and complex framework, i wouldn't consider it for microservice - but i would use it to create website

---

Zig //I really like the idea of zig, i want to like zig; but i don't like it's syntax sadly; i'll play with zig later
1. ZZZ -- 
2. ZAP -- 
3. Stdlib -- 

Nim //Seems quite raw to me, i'll test it later; but i wouldn't consider it for production anyways
1. Httpbeast -- 
2. Jester -- 
3. Stdlib -- 

---

#### Objective
In progress

