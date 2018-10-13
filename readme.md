# ![RealWorld Example App](logo.png)

> ### Golang + Revel codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.


### [Demo](https://github.com/gothinkster/realworld)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)


This codebase was created to demonstrate a fully fledged fullstack application built with **[YOUR_FRAMEWORK]** including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the **Golang + Revel** community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.


# How it works

> Describe the general architecture of your app here

# Getting started

1. Install _docker_ and _docker-compose_ to your system
2. Add `127.0.0.1 api.realworld.wip` to your `/etc/hosts` file
3. Copy `.env.docker` to `.env` in the project root
4. Start [nginx-proxy](https://github.com/jwilder/nginx-proxy)
5. Generate a Docker bundle from the Compose file `docker-compose build`
6. Create and start containers with demon mode `docker-compose up`

* To start the **nginx-proxy**, type the following command:

`docker run -d -p 80:80 -p 443:443 --name nginx-proxy --net reverse-proxy -v $HOME/certs:/etc/nginx/certs:ro -v /etc/nginx/vhost.d -v /usr/share/nginx/html -v /var/run/docker.sock:/tmp/docker.sock:ro --label com.github.jrcs.letsencrypt_nginx_proxy_companion.nginx_proxy=true jwilder/nginx-proxy`
