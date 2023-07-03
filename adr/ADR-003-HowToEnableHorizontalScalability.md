# ADR-003: How to enable horizontal scalability

## Status

Accepted

## Context
The CSMS needs to be able to manage a large (50k+) pool of charge stations. Besides overcoming the limited amount of
possible websocket connections, we need a solution which allows for high performance (99% calls completed within 1s) 
regardless of volume of requests. 

## Decision
We will decouple the websocket handling to a separate process (gateway) and use a MQTT message broker to pass messages 
between the gateway and the CSMS. While there are a few drawbacks from the approach, we concluded that the advantages 
well outweigh those.

## Consequences

There are a number of notable advantages to use this approach:

*Advantages*:
* Independent scaling of the CSMS and the websocket server
* Decoupled architecture
* Regarding MQTT:
    * low footprint
    * lightweight topic management: creating a large amount of topics on the fly is easy to achieve.
    * ability to subscribe to single topics
    * shared subscriptions
* Regarding gateway:
    * RPC logic is completely decoupled from CSMS, which will make adding OCPP1.6 support easier
* Multiple broker implementations (managed, such as AWS IoT, 3rd party such as HiveMQ and open source such as Mosquitto) 

*Drawbacks*:
* MQTT doesn't offer persistence. This isn't something we require as fast throughput is fundamental, the ability to 
replay past messages has no value in OCPP.
* Complexity: Adding a double abstraction layer adds a lot of complexity to the architecture. We consider this to be a
minor issue compared to the value it brings.
* Cognitive load: adding MQTT increases the tech stack the team needs to understand. Also, MQTT is probably less known 
to most developers compared to other alternatives
* More complex to launch: Launching the application now requires a lot of different dependencies:
    * MQTT broker
    * Websocket gateway
    * CSMS
  This can be mitigated by using docker-compose

### Alternatives considered

**SNS**
*Advantages*: 
* robust
* well-known
* persistence

*Disadvantages*: 
* AWS only
* heavy-weight 
* upfront stream creation
* persistence is different

**Kafka**
*Advantages*: 
* robust
* well-known
* persistence

*Disadvantages*: 
* topics are heavyweight
* every subscriber would have to consume all the topics