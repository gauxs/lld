# LLD vs Machine Coding Round

While both rounds evaluate the same **core problem domains** (such as a ride-sharing app, parking lot, or expense splitter), the **prompt framing and expectations** differ significantly. An **LLD round** tests architectural breadth and interactive design thinking in **45–60 minutes**, whereas a **Machine Coding (MC) round** tests end-to-end execution speed, production code hygiene, and live extensibility in **90–120 minutes**.

## How the Prompt Varies

| Dimension | LLD Round Prompt | Machine Coding (MC) Round Prompt |
| :--- | :--- | :--- |
| **Typical Archetype** | **Real-World System (e.g., Uber/Ola hailing):** Real-time driver matching, dynamic surge pricing, GPS/location updates, and full trip lifecycle. | **Rule-Driven Workflow (e.g., BlaBlaCar carpooling):** Fixed functions like `add_user`, `add_vehicle`, `offer_ride`, and `select_ride` with explicit rules. |
| **Requirement Clarity** | **Ambiguous by design:** A 2–3 line prompt (*"Design a ride-sharing app"*). You must ask clarifying questions to scope features. | **Exhaustive 1–2 page spec:** Exact method signatures, validation rules (e.g., max 1 active ride per vehicle), and mandatory bonus features. |
| **Data & Persistence** | **Database & APIs:** Expects relational/NoSQL table schemas, indexing strategies, and REST/WebSocket endpoint definitions. | **Strictly In-Memory:** Explicitly forbids databases or external frameworks; everything lives in `Maps` and `Lists`. |
| **Execution Format** | **Diagrams & Pseudocode:** Whiteboarding class relationships, sequence flows, and critical locking blocks. | **Driver Class / Test Cases:** Requires a runnable `main` method executing sample input/output flows cleanly. |

---

## LLD Grading
1. **Encapsulation & Clear Responsibilities:** Clean separation of models, services, and repositories.
2. **SOLID Principles:** Extensible abstractions without modifying existing core logic.
3. **Composition Over Inheritance:** Favoring pluggable behaviors (e.g., pricing or matching interfaces) over rigid class hierarchies.
4. **Design Patterns:** Appropriate use of Strategy, Factory, State, Observer, etc., without over-engineering.
5. **Testability:** Loose coupling and dependency injection so components can be unit-tested easily.

## LLD Steps
1. **Clarify requirements and scope (3–5 min):** Nail down core actors, scale assumptions, and out-of-scope features.
2. **Identify core entities and relationships (5 min):** Map out models, cardinalities, enums, and state transitions.
3. **Define interfaces and apply SOLID + patterns (5 min):** Establish service boundaries and extension points.
4. **Write the code (15–20 min):** Implement the core classes and the 1–2 critical methods carrying the system's logic.
5. **Walk through scenarios and discuss extensions (5 min):** Trace an end-to-end flow, handle concurrency/edge cases, and discuss follow-up features.

---

## MC Grading
1. **Working, Bug-Free Code:** Compiles and passes all provided sample test cases cleanly.
2. **Clean, Modular Structure:** Readable naming, proper package/folder organization, and single-responsibility classes.
3. **State & Concurrency Handling (when relevant):** Proper handling of shared state, locks/thread safety, or data consistency when multiple actors/requests are simulated.
4. **Live Extensibility:** Ability to plug in a new rule or strategy during the 30-minute review with minimal code changes.
5. **Time & Scope Management:** Getting a working MVP running first before polishing or adding bonus features.

## MC Steps
1. **Read requirements & prioritize (5–10 min):** Separate must-have core flows from nice-to-have bonus features.
2. **Sketch entities, interfaces & folder structure (10 min):** Set up `models`, `services`, `strategies`, and `exceptions` cleanly.
3. **Build core MVP end-to-end with a runnable driver (45–60 min):** Get sample inputs and expected outputs working early.
4. **Handle edge cases, validations & concurrency (15–20 min):** Add input checks, custom exceptions, and thread-safe locking if needed.
5. **Refactor & prepare for live extensions (10 min):** Clean up rough edges and verify where new strategies can be added quickly.
