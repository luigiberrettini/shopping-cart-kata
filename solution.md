# Shopping cart kata solution



## Introduction

Due to the scale needed an e-commerce system, usually relies, on multiple distributed systems each one involving many components.

This means that the design and implementation of even only one element has an high level of complexity also to support the versatility required by the business domain.

It is pretty common to use [Domain-Driven Design](https://wikipedia.org/wiki/Domain-driven_design) to analyze an area of the business, a DDD subdomain, maybe with some [EventStorming](https://www.eventstorming.com) sessions, with the goal of getting a better awareness of concepts, processes and rules that belong to it and start talking all the same ubiquitous language.

Multiple subdomains cooperate by means of integration events and web APIs and even within a single subdomain it is possible to find one or more [microservices](https://microservices.io) that, whenever possible, use a [choreography rather than orchestration](https://www.thoughtworks.com/insights/blog/scaling-microservices-event-stream) integration model.

The implementation of a single microservice is often based on the [ports and adapters](https://web.archive.org/web/20060711221010/http://alistair.cockburn.us:80/index.php/Hexagonal_architecture) pattern also known as hexagonal architecture.

Moreover, when different use cases need different levels of scale and flexibility, the [CQRS](https://www.martinfowler.com/bliki/CQRS.html) pattern is used, resulting in different write and read models often needing polyglot persistence, sometimes in combination with [Event Sourcing](https://www.martinfowler.com/eaaDev/EventSourcing.html).

This is what goes under the term [Event-Driven Architecture](https://martinfowler.com/articles/201701-event-driven.html).

The goal is usually taming complexity by decomposing a huge system in multiple simpler parts with only one responsibility.
Unfortunately this brings along a big increase in the overall system complexity and the challenges related to its implementation, deploy (e.g. SDLC, release management, containerization and container orchestrators like [Kubernetes](https://kubernetes.io)), handling ([observability](https://distributed-systems-observability-ebook.humio.com/)), maintenance, evolution.



## Overview

One of the key element for any retailer is the **shopping cart**, essential both for customers to place order and for the company to make profit.

Considering the small number of articles in the catalog, the not so complex promotion rules and the quite simple cart features to be supported, I believe that the solution should be as simple as possible.

It would be overkill to use DDD, CQRS, ES, a full-fledged hexagonal architecture and code with full instrumentation for observability (logging, application and business metrics monitoring, distributed tracing at least).

The client application should communicate with an API, acting in the same way as usually a backend server does behind a frontend component, be it a [Backend For Frontend](https://samnewman.io/patterns/architectural/bff) or the backend part of a [microfrontend](https://micro-frontends.org).

A cart microservice exposing the functionalities required through a web API without authentication or authorization is the right choice in this case, but there would be no reason to deviate from the goal and also implement separate microservices for the catalog and promotion features, especially considering that it would require a huge effort to mimic only a little part of the features of a modern persistence store and that actually catalog and promotion are two other subdomains.

Tests should be used to guarantee that main use cases are covered and the software is maintainable and evolvable without introducing misbehaviors or regressions.

As per the build and test phase few or no scripts should be used and it would be premature to introduce Kubernetes and [Helm](https://helm.sh) at this stage.



## Proof of concept
The PoC contains a possible [Go](https://golang.org/) implementation of the solution developed using [Visual Studio Code](https://code.visualstudio.com):
 - a command line client
 - a web API leveraging on the [Gorilla mux](http://www.gorillatoolkit.org/pkg/mux) package

Code has been tested for the main use cases and behaviors, but no benchmark or test coverage has been used.

The build and test phase supports [Docker](https://www.docker.com) and [Docker Compose](https://docs.docker.com/compose)


### Main features of the web API
  - Level 2 of the [Richardson maturity model](https://www.martinfowler.com/articles/richardsonMaturityModel.html)
  - Media type `application/json`: no hypermedia controls even if links are embedded in responses and location headers are used
  - Support for conditional requests: `ETag`, `If-None-Match`, `If-Match` headers
  - In-memory storage, implemented with simple data structures, to handle articles, carts (and their ETags) and promotion rules
  - Add article with quantity (`POST`) and set article quantity (`PUT`) routes to implement the desired add article capability
  - Catalog route only to support client (improperly put in the cart service to avoid creating an API only for it)
  - The promotion engine, based on rules related to an item or the cart, determine percentage/value discounts or new values that:
     - are applied to part of the quantity of a cart item
     - are applied to the cart subtotal
     - determine a present with a quantity to add to the cart (usually a gift or a sample)
     - are applied to shipping costs


### Environment setup

#### Development environment setup
 1. Install [Go SDK](https://golang.org/dl)
 2. Execute the following commands from the terminal:
    ```shell
    cd $GOPATH
    cd src
    git clone http://github.com/luigiberrettini/shopping-cart-kata
    cd shopping-cart-kata
    go get -d -v ./...
    ```
 3. Run tests by executing the following command from the terminal (`-count=1` prevents caching):
    ```shell
    go test ./... -timeout=60s -parallel=4 -count=1
    ```
 4. Run the applications by executing the following commands from the terminal (check that the `bin` subdirectory of `GOPATH` is in the `PATH`):
    ```shell
    # First terminal
    go install -v ./...
    cartsvc
    cartcli
    ```
    ```shell
    # Second terminal
    cartcli
    ```
 
#### Plain docker Linux local environment setup
 1. Install [Docker](https://docs.docker.com/install)
 2. Execute the following commands
    ```shell
    git clone http://github.com/luigiberrettini/shopping-cart-kata
    cd shopping-cart-kata
    docker build . -t cartools

    docker network create -d bridge bridgenet
    docker run --name cartsvc --hostname cartsvc --net bridgenet -p 8000:8000 -d --rm cartools /cartsvc -listen=cartsvc:8000 -authority=cartsvc:8000
    docker run --name cartcli --hostname cartcli --net bridgenet -ti --rm cartools /cartcli -baseUrl=http://cartsvc:8000
    
    docker network rm bridgenet
    ```

#### Docker Compose Linux local environment setup
 1. Install [Docker](https://docs.docker.com/install)
 2. Install [Docker compose](https://docs.docker.com/compose/install)
 3. Execute the following commands
    ```shell
    git clone http://github.com/luigiberrettini/shopping-cart-kata
    cd shopping-cart-kata
    docker-compose build

    docker-compose run cartcli

    docker-compose down
    ```

#### Docker Linux virtual machine environment setup
 1. Install [VirtualBox and its Extension Pack](http://www.virtualbox.org/wiki/Downloads)
 2. Install [Vagrant](https://www.vagrantup.com/downloads.html)
 3. Download the [Vagrantfile](Vagrantfile) provided with this repo
 4. Execute the command `vagrant up`
 5. Connect via SSH to localhost on port 2200 with username `vagrant` and password `vagrant`
 6. Execute step 3 of the previous section



## Evolutions
  1. Use a static analysis tool to keep code quality high
  2. Use code coverage tool
  3. Add a few integration tests (evaluate use of Docker)
  4. Instrument the application for distributed tracing
  5. Track and publish application metrics (add also benchmarks) and business KPIs
  6. Add logging in a way that it is possible to easily switch the logging target (e.g. terminal, file, DB)
  7. Add authentication and authorization to convert anonymous carts into the user cart
  8. Use Kubernetes and Helm to support advanced deployment and scalability scenarios
  9. Use DDD, CQRS, hexagonal architecture, domain and integration events and evaluate using ES
 10. Support article removal explicitly and by setting the quantity to zero
 11. Extend support for conditional HTTP requests:
      - strong ETag validation
      - `Last Modified`
      - `If-Modified-Since`
      - `If-Unmodified-Since`
      - `If-Range`
 13. Move to a higher [Richardson Maturity Model](https://www.martinfowler.com/articles/richardsonMaturityModel.html) and [Amundsen Maturity Model](http://www.amundsen.com/talks/2016-11-apistrat-wadm/2016-11-apistrat-wadm.pdf) also using a proper [API design methodology](https://www.infoq.com/articles/web-api-design-methodology/) and a high [H factor](http://amundsen.com/hypermedia/hfactor) media type ([comparison chart](http://gtramontina.com/h-factors)) like [Mason](https://github.com/JornWildt/Mason), [Hyper](http://hyperjson.io/spec.html) or [UBER](https://rawgit.com/uber-hypermedia/specification/master/uber-hypermedia.html)
 12. Implement catalog service and subdomain (evaluate using GraphQL)
 13. Implement promotion service and subdomain (evaluate using GraphQL)
 14. Use a distributed cache for carts persisting logged user carts also on a NoSQL store
 15. Use a [Lucene](http://lucene.apache.org)-based store to search the catalog
 16. Use a promotion engine relying on a rule engine to back the promotion service handling complex business rules support