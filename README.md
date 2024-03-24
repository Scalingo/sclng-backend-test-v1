# Backend Technical Test at Scalingo

## Preamble

* This project aim to deliver some Endpoints to treat github repositories data.

## Architecture

* The project code was restructed using well known REST API golang archetype. This may led to a clean and best-practice-compliant approach when the project size will grow. The code is organized as follow:

1. cmd/: This directory contains the main applications for the project. The sub-directory name match the desired executable name i.e cmd/sclng-backend-test-v1.

2. internal/: This folder holds private application code. This is where we place our domain types, such as AgrRepo struct, handlers, router, ...etc. This folder contain two packages: one for models and another for the API package (server package). For instance, there is no unit code testing. Whereas, with each file we must add a unit testing file for all function exp. 'handlers_test.go' must holds all unit tests for the 'handlers.go' file.

3. config/: This directory is where we retrieve the configuration of our application like  environment variables (exp. Listening  Port), perform dependency injections. In a realword  application, the configuration must be read from a configuration file (YAML config file) to allow configuration injection.

4. tests/: This folder must contains integrations test files like Venom tests. For instance it is an empty folder. 

# Software design

* Our API listen on the port 5000 of the localhost host,

* Other the Ping Endpoint, it implemente two Endpoints which are: /listRepo and /listAgragateRepo

- /listRepo: list all data of the last 100 public repositories,

- /listAgragateRepo: List agregated data, which are in the AgrRepo struct,  of the last 100 public repositories. Data concerned are: Repository FullName, Owner, Repository Name and the most imporatante info is the languages part. 
In the languages part we agregate statistiques about all the languages of the repository. For performance enhancing, we use 'goroutine' to parallelize languages statisqtique  checkout. We use  Mutex for synchronization (Protect access to a shared slice) and WaitGroup to wait for all the goroutines launched to finish.


* The two Endpoint accept any number of parameters with their values regarding the github REST API search Endpoints like language, licence, ...etc. For an exaustive list of search critera see https://docs.github.com/fr/rest/search/search?apiVersion=2022-11-28#search-repositories  and https://docs.github.com/fr/rest/search/search?apiVersion=2022-11-28#constructing-a-search-query

Exemple:  to limit the results of our Endpoints only to Repository of 'Java' language and licence 'MIT' we must add 'language=Java&license=mit' as URL parameters.

* Our API use goroutine 

## Execution

```bash
docker compose up
```

Application will be then running on port `5000`

## Tests

1. To get the last 100 repositories data do:

```bash
$ curl  "http://localhost:5000/listRepo"
```

2.  To limit our last 100 repositories to those with Java language and MIT License :
 
```bash
$ curl  "http://localhost:5000/listRepo?language=Java&license=mit"
```

3. To get agregated data for repositories languages do:

```bash
$ curl  "http://localhost:5000/listAgragateRepo"
```
4. To limit our agregated data, for the last 100 repositories, to those with Java language and MIT License :

```bash
$ curl  "http://localhost:5000/listAgragateRepo?language=Java&license=mit"
```
5. To test the healthy of the application we can regiter somme health check routes exp. 'health/alive' and '/health/ready'. These kind of routes are importante in some complex deployment environement in which we integrate several components. Here, we let only the canevas Endpoint ping check to check the readness of the server: 
```bash
$ curl localhost:5000/ping
```
