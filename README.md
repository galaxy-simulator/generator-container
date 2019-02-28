[![Go Report Card](https://goreportcard.com/badge/git.darknebu.la/GalaxySimulator/generator-container)](https://goreportcard.com/report/git.darknebu.la/GalaxySimulator/generator-container)
# generator-container

The generator container generates stars by [random-sampling](https://en.wikipedia.org/wiki/Simple_random_sample) random values on the [Navarro-Frenk-White profile](https://en.wikipedia.org/wiki/Navarro%E2%80%93Frenk%E2%80%93White_profile).

In order to do so, it functions as a manager-container creating random stars.
In the end, the generator container doesn't calculate the NFW-values by itself, but accesses the [NFW-containers](https://git.darknebu.la/GalaxySimulator/NFW-container) that
is specialized on this job.

This has the advantage that a lot of manager containers can be started an the generator container
interacts with a reverse proxy such as [traefik](https://traefik.io) that load-balances the request.
The complete process of generating stars can though be speed up a lot. 

Instructing the generator container about where to look for the nfw-containers works in the
following way: simply set the environment variable `nfwurl` to the base-url of the nfw-container or proxy.

## Examples

- single nfw-container

    Let's suppose we've got a generator-container running on `localhost:8080` and we want it
    to access the nfw-container running on `localhost:8081`. The solution is to simply set the
    `nfwurl` to `localhost:8081`:
    ```bash
    $ go build . -o generator
    $ nfwurl=localhost:8081 ./generator 
    ```
   As soon as we access the /gen endpoint of the generator endpoint, the generator generates
   a random star and makes a request to the NFW-container with the coordinates of the star.
   The nfw-container returns the NFW-value of the given star and allows the generator container
   to continue testing if the star exists or not.
    
- multiple nfw-containers

    Let's suppose the amount of stars we want to generate is really big and generating them on
    a single host is not fast enough: the solution is to use more hardware!
    
    So let's assume we've got multiple nodes running in a docker-swarm environment that uses
    traefik as a reverse-proxy. We instruct the generator-container to use the nfw.docker.localhost
    endpoint that is balanced through traefik.
    ```bash
    $ go build . -o generator
    $ nfwurl=nfw.docker.localhost ./generator 
    ```
