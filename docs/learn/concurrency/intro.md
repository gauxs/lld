---
title: Introduction to concurrency
description: Shared memory, interleaving, and the three problem types you see in LLD interviews.
next: /learn/concurrency/correctness
---

# Introduction to concurrency

<p class="lead">Concurrency is what happens when multiple execution paths make progress at the same time and touch shared state. In low-level design interviews you usually stay inside one process—threads, locks, and in-memory structures—not distributed systems.</p>

Concurrency does not appear in every interview, but for senior roles it often shows up as a follow-up: two clients book the same seat, two goroutines update a counter, or ten requests hit a cache while it refreshes. The interviewer wants to see whether you can name what breaks when actions overlap and fix it without over-engineering.

## Concurrency fundamentals

A **process** is an isolated container with its own address space. Inside it, the runtime creates **threads** (or goroutines) that share the heap and globals but have their own stacks and program counters.

On multiple cores, threads may run in parallel. On one core, the scheduler **interleaves** instructions. Either way, operations from different threads can interleave in ways your source code does not make obvious—especially when one logical step is several machine instructions.

Assume concurrency whenever **shared mutable state** exists. Java, Go, C++, and Rust all face this in production code. User-facing JavaScript is mostly single-threaded with an event loop; LLD concurrency there is usually modeled as async handoff rather than shared-memory threads.

## The toolbox (quick reference)

| Primitive | Use when |
| --- | --- |
| **Atomics** | Single variable read-modify-write (counters, flags) |
| **Mutex / lock** | Critical sections; multi-field updates; check-then-act |
| **Semaphore** | Cap concurrent operations (pool size, rate budget) |
| **Condition variable** | Wait until a predicate becomes true |
| **Blocking queue / channel** | Producer–consumer handoff between threads |

Go idioms: protect shared structs with `sync.Mutex` or confine mutation to one goroutine; use **channels** for coordination. See the language-specific pages in later articles.

## Three problem types

Most interview concurrency questions fall into three buckets:

| Type | What breaks | Typical tools |
| --- | --- | --- |
| [Correctness](/learn/concurrency/correctness) | Lost updates, double booking, torn reads | Locks, atomics, confinement |
| [Coordination](/learn/concurrency/coordination) | Producers/consumers out of sync | Queues, channels, condition vars |
| [Scarcity](/learn/concurrency/scarcity) | Too many concurrent users of a limited resource | Semaphores, pools |

Real systems often mix all three. Separating them helps you pick the smallest fix for each part.

## What's next

Start with **Correctness**—how shared state gets corrupted and how to guard it—then **Coordination** and **Scarcity** for handoff and resource limits.
