# Meeting room

Design and implement a Meeting Room Scheduler.

## Functional requirements

### FR-1: Meeting room details
 - There can be multiple meeting rooms; the exact number should be configurable.

> **Interview prompts**
>
> - **How many meeting rooms can be there?**

### FR-2: Book a meeting
 - A booking can be made up to 30 days in advance.
 - A room can be booked for a minimum of 30 minutes and a maximum of 4 hours.

> **Interview prompts**
>
> - **How far ahead booking can we do?**
> - **what is the min and max time range for which we can book a room??**

### FR-3: Meeting details
 - A meeting can have up to 20 participants.
 - Meeting details can be limited to:
    - Title
    - Start time
    - End time
    - Participants

> **Interview prompts**
>
> - **How many participants can be there in on meeting?**

### FR-4: Editing a meeting
 - Meeting time: Can be edited.
 - Participants: Can be added or removed.
 - A meeting can be cancelled while it is in progress. It cannot be cancelled after its scheduled end time.
 - Edits are allowed until the meeting's scheduled end time.

> **Interview prompts**
>
> - **Can a meeting time and participants edited? and till when are the edits allowed?**

### FR-5: Meeting communications
- Send meeting notifications/details to participants. Notifications are sent when a meeting is edited, including changes to:
    - Meeting time
    - Participant list
    - Cancellation

> **Interview prompts**
>
> - **Are we also building notification system to send meeting details to participants?**

## Out of scope (this extension)

> **Interview prompts**
>
> - **Deferred idea?** Point to a child extension or “not in this scope.”

## Non-functional requirements

### NFR-1: Concurrency

Concurrent booking attempts are expected, including multiple users attempting to book the same room for overlapping time slots.

> **Interview prompts**
>
> - **Concurrent booking will happen?**

## Extensions from here

Optional follow-ups (sibling extensions). Add a row when you create `extensions/<child_id>/` and register it in `extensions.json`.

- _(none yet)_ — e.g. networked play, configurable rules
