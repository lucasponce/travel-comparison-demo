# travel-comparison-demo
A Microservices demo based on Istio Service Mesh tool. 

This goal of this demo is to demostrate [Istio](https:/istio.io/) cappabilities observed and managed by [Kiali](https://kiali.io) tool. 

## Travel Portal and Travel Agency

This demo creates two groups of services to simulate a travel portal scenario.

In a first namespace called **travel-portal** there will be deployed a set of services that will represent the portals 
where users access to search and book travels.

One of the characteristics of **travel-portal** is that there will be different portals depending of which city users are 
interested on. In this example, there will be three portals for *Rome*, *Paris* and *London*.

Also, every portal will have two versions, one version to handle regular web traffic and a second version that will 
manage vip users with special offers and discounts.

For our demo, we don't care about the details of the application but the relationships between services as
it is shown in the picture:

[[images/travel-portal.png|alt=travel-portal]]   



