---
title: Interview approach
sidebar: Interview approach
description: Shared state, invariants, then pick a synchronization primitive—before reaching for a mutex.
---

# The mental model for interviews

```mermaid
graph LR
    %% Setup the entry point
    Start([Interviewer Introduces Concurrency])
    
    %% The Trap / Wrong Path
    Start -->|Instinct | Trap[/"❌ 'I will use a Mutex!'"/]
    style Trap fill:#fecaca,stroke:#dc2626,stroke-width:2px,stroke-dasharray: 5 5
    
    %% The Right Path
    Start -->|Better Approach| Step1
    
    Step1["<strong>Step 1: What is shared? (State)</strong><br>• Inventory, parking spots, balance, queue"]
    --> Step2["<strong>Step 2: What happens concurrently? (Actions)</strong><br>• reserve(), cancel(), update(), read()"]
    --> Step3["<strong>Step 3: What must remain true? (Invariants)</strong><br>• inventory >= 0<br>• 1 seat = max 1 reservation"]
    --> Step4["<strong>Step 4: Pick the Primitive</strong><br>• Choose simplest mechanism<br>• (Atomic, Mutex, RW Lock, Semaphore)"]
    --> Takeaway("<strong>Key Takeaway:</strong><br>Concurrency = unpredictable interleaving of operations.<br>Identify state + invariant first, then pick the simplest tool.")

    %% Styling for paths
    style Start fill:#fff7ed,stroke:#ea580c,stroke-width:2px
    style Step1 fill:#f8fafc,stroke:#94a3b8,stroke-width:1px
    style Step2 fill:#f8fafc,stroke:#94a3b8,stroke-width:1px
    style Step3 fill:#f8fafc,stroke:#94a3b8,stroke-width:1px
    style Step4 fill:#f8fafc,stroke:#94a3b8,stroke-width:1px
    style Takeaway fill:#d4edda,stroke:#28a745,stroke-width:2px

```

### The key takeaway
For LLD, don't think:
> Concurrency = locks.

Think:
> Concurrency = unpredictable interleaving of operations on shared state.

Your job is to identify the shared state + invariant, then choose the simplest mechanism that preserves that invariant.
