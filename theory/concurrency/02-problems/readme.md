---
title: Three concurrency problem types
sidebar: Overview
description: Correctness, coordination, and scarcity—what breaks and typical tools.
---

## The 3 concurrency problems you should recognize
Most LLD concurrency questions fall into three buckets:

| Problem | Core question | Typical solution |
| :--- | :--- | :--- |
| **Correctness** | Can concurrent operations corrupt state? | Locks, atomics |
| **Coordination** | How do threads communicate/wait? | Queues, conditions, channels |
| **Scarcity** | How do we limit access to a finite resource? | Semaphores, pools |



#### Correctness
```text
Two threads → same seat
```
Need to protect shared state.

#### Coordination
```text
Producer → Queue → Consumer
```
Consumer may need to wait until work exists.

#### Scarcity
```text
100 requests
     ↓
10 DB connections
```
Only 10 operations can use the resource concurrently. These categories are more useful in interviews than memorizing individual concurrency APIs.
