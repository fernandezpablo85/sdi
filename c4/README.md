# Chapter 4: Design a Rate Limiter

## Overview

Code for rate limiting system stuff, based on Chapter 4 of "System Design Interview" by Alex Xu.

## Problem Statement

Design a rate limiter to control the rate of traffic sent by clients. This is a critical component for:
- Preventing resource starvation from DoS attacks
- Reducing cost by limiting excess requests
- Preventing server overload

## Implementations

- [Go implementation](golang/README.md) - Token bucket algorithm with Redis backend
