# Tech stack choice for Mocks

## Status

Accepted

## Context

We want to select which languages / platforms we will be using for the development of the mocked components within this project. The project contains multiple components and the tech choices can vary for each component.
The original purpose of this project was a timeboxed PoC. But after the PoC, when the project was picked up for further development, time constraints remain a large factor. Other factors to take into consideration are:
* Speed of development
* Availability of developers

The mocked components considered for this ADR are:
* CS/EV backend
* CS/EV frontend
* eMSP

## Decision

### CS/EV backend
We will use Python for backend development of this component, because the original code was forked from mobilityhouse, which used Python.

### CS/EV frontend
For the front-end of this component we will use Typescript on Node.js. The main reason being that the PoC was started in Typescript/Node.js

### eMSP
We will use Go for the development of this mocked component. The main reason being to avoid adding yet another language. Go is well suited for web service development, and we can utilize CSMS for a reference on patterns / approaches.

## Consequences

### CS/EV backend

*Advantages*:

* Base code already developped in Go

*Drawbacks*:

* Python is not well suited for multithreaded web service development.
* Deployment to AWS is complicated.

### CS/EV frontend

*Advantages*:

* Base code already developped in Go
* Typescript offers static typing

*Drawbacks*:

* Dependency management could become complicated.

### eMSP

*Advantages*:

* Easy for local development; no need for complex build scripts; new developers can get started very quickly.

* Sinergy with CSMS component: we can use the CSMS component as a reference, which increases consistency.

*Drawbacks*:

*  Availability of experienced Go developers. 
