# Tech stack choice for CSMS

## Status

Accepted

## Context

We want to select which languages / platforms we will be using for the development of the CSMS (Charge Station Management Service) component of this project. The project contains multiple components and the tech choices can vary for each component.
The original purpose of this project was a timeboxed PoC. But after the PoC, when the project was picked up for further development, time constraints remain a large factor. Other factors to take into consideration are:
* Available connectivity
* Performance
* Ease of development
* Availability of developers

## Decision

We will use Go (golang) for development of this component. The main reason being that the PoC was started in Go due to historical / circumstantial reasons. 

## Consequences

*Advantages*:

* Base code already developed in Go
* Go scores high in performance and memory footprint
* Easy for local development; no need for complex build scripts; new developers can get started very quickly.

*Drawbacks*:

* Availability of experienced Go developers. 