# LLD vs Machine coding round
Both rounds generally use the same exact problem statements, an LLD rounds tests architectural breadth and interactive design thinking in 45-60 minutes, wheras a machine coding round tests end-to-end execution speed, production code hygiene and live extensibility under 90-120 minutes.

## LLD grading
1. Encapsulation and clear responsibilities.
2. The SOLID principles
3. Composition over inheritance.
4. design patterns
5. Testability

## LLD Steps
1. Clarify requirements and scope (3–5 min)
2. Identify the core entities and relationships (5 min)
3. Define interfaces and apply SOLID + patterns (5 min)
4. Write the code (15–20 min): mplement the core classes and the one or two methods that carry the system's logic
5. Walk through scenarios and discuss extensions (5 min)

## MC grading
1. Working, bug-free code
2. Clean, modular structure
3. State & concurrency handling (when relevant): Proper handling of shared state, locks/thread safety, or data consistency if multiple actors/requests are simulated.
4. Live extensibility
5. Time & scope management: Getting an MVP working first before polishing or adding optional features.

## MC Steps
1. Read requirements carefully & identify must-have vs nice-to-have (5–10 min)
2. Sketch entities, interfaces & folder structure (10 min)
3. Build core MVP end-to-end with a runnable driver/demo (45–60 min)
4. Handle edge cases, validations & concurrency if applicable (15–20 min)
5. Refactor for cleanliness & prepare for live follow-up extensions (10 min)