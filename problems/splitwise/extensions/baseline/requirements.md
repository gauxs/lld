# Splitwise

Design a expense-sharing application that helps groups of people such as roommates, friends, coworkers, or travelers split bills and keep track of who owes whom.

## Functional requirements

### FR-1: Support adding users with basic profile information (name, email, phone)

### FR-2: Support creating groups of users for shared expenses

### FR-3: Support adding expenses with three split types: equal, exact amounts, and percentage-based

### FR-4: Track pairwise net balances between users (who owes whom and how much)

### FR-5: Support partial and full settlements between any two users

### FR-6: Notify users when expenses are added or settlements occur

### FR-7: Handle rounding differences in equal splits so the total always matches

## Non-functional requirements

### NFR-1: The design should follow object-oriented principles with clear separation of concerns

### NFR-2: The system should handle concurrent expense additions without race conditions

### NFR-3: The system should be modular and extensible to support new split types

### NFR-4: The components should be testable in isolation

## Out of scope (this extension)

- Deferred capability.
