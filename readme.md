# ![RealWorld Example App](logo.png)

> ### Golang + Revel codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.


### [Demo](https://github.com/gothinkster/realworld)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)


This codebase was created to demonstrate a fully fledged fullstack application built with **Golang + Revel** including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the **Golang + Revel** community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.


# How it works

Below is the recommended layout of a Revel application, supplemented with domain entities and services.

    my_gocode/ - GOPATH root
        src/ - GOPATH src/ directory
            github.com/revel/revel/ - Revel source code
            bitbucket.org/me/sample/ - Sample app root
                entities/ - domain entities
                app/ - app sources
                    controllers/ - app controllers
                        init.go - interceptor registration
                    models/ - app domain models
                    jobs/ - app domain jobs
                    services/ - app domain services
                    routes/ - reverse routes (generated code)
                    views/ - templates
                    tmp/ - app main file, generated code
                tests/ - test suites
                conf/ - configuration files
                    app.conf - main configuration file
                    routes - routes definition file
                messages/ - i18n message files
                public/ - static/public assets
                    css/ - stylesheet files
                    js/ - javascript files
                    images/ - image files
                    
# Getting started

1. Install _docker_ and _docker-compose_ to your system
2. Add `127.0.0.1 api.realworld.wip` to your `/etc/hosts` file
3. Copy `.env.docker` to `.env` in the project root
4. Start [nginx-proxy](https://github.com/jwilder/nginx-proxy)
5. Generate a Docker bundle from the Compose file `docker-compose build`
6. Create and start containers `docker-compose up`

* To start the **nginx-proxy**, type the following command:

`docker run -d -p 80:80 -p 443:443 --name nginx-proxy --net reverse-proxy -v $HOME/certs:/etc/nginx/certs:ro -v /etc/nginx/vhost.d -v /usr/share/nginx/html -v /var/run/docker.sock:/tmp/docker.sock:ro --label com.github.jrcs.letsencrypt_nginx_proxy_companion.nginx_proxy=true jwilder/nginx-proxy`
