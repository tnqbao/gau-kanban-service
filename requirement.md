## Backend Requirements (Checklist)

- **Boards**
    - [x] `POST /boards` → create board (title, description, owner, team)
    - [x] `GET /boards` → list boards for user
    - [x] `GET /boards/:id` → get board details with columns, tickets, members, and labels
    - [x] `PATCH /boards/:id` → update title, description
    - [x] `DELETE /boards/:id` → delete board (owner only)
    - [x] `PATCH /boards/:id/archive` → archive board (owner only)
    - [ ] persist last opened board in `user_preferences`

- **Columns**
    - [x] `POST /columns` → create column (board_id, title, wipLimit) - auto append to end
    - [x] `GET /columns/:id` → get column details with tickets
    - [x] `PATCH /columns/:id` → update name, wipLimit
    - [x] `DELETE /columns/:id` → delete column
    - [x] `PATCH /columns/:id/reorder` → reorder columns (fractional ordering, no board_id needed)
    - [ ] enforce WIP limit (backend validation + return count/limit)

- **Cards/Tickets**
    - [x] `POST /tickets` → create ticket (column_id, title) with auto-generated ticket number
    - [x] `GET /tickets/:id` → get ticket details (markdown description, checklist, labels, assignees)
    - [x] `PATCH /tickets/:id` → update ticket (title, description, priority, due_date, status)
    - [x] `DELETE /tickets/:id` → delete ticket
    - [x] `PATCH /tickets/:id/move` → move ticket to new column/order
    - [ ] maintain orderIndex in DB for optimistic reordering

- **Labels**
    - [x] `POST /labels` → create label for board (board_id, title)
    - [x] `POST /labels/assign` → assign label to ticket (label_id, ticket_id)
    - [ ] `GET /labels/:id` → get label details
    - [ ] `PATCH /labels/:id` → update label
    - [ ] `DELETE /labels/:id` → delete label
    - [ ] `DELETE /tickets/:id/labels/:labelId` → remove label from ticket

- **Assignees**
    - [x] `POST /assignees` → assign member to ticket (ticket_id, member_id)
    - [ ] `GET /tickets/:id/assignees` → get ticket assignees
    - [ ] `DELETE /tickets/:id/assignees/:memberId` → remove assignee from ticket

- **Checklists**
    - [x] `POST /checklists` → create checklist item (ticket_id, title)
    - [x] `PATCH /checklists/:id` → update checklist item
    - [x] `DELETE /checklists/:id` → delete checklist item
    - [ ] `GET /tickets/:id/checklists` → get ticket checklists
    - [ ] `PATCH /checklists/:id/toggle` → toggle checklist completion

- **Members**
    - [x] `POST /members` → add member to board (user_id, board_id)
    - [x] `GET /boards/:id/members` → get board members (included in board details)
    - [ ] `DELETE /members/:id` → remove member from board

- **Users**
    - [x] User creation/retrieval integrated with board creation
    - [x] `POST /users` → create user with custom user_id and fullname
    - [ ] `GET /users` → list users for member assignment
    - [ ] `GET /users/:id` → get user details

- **Filtering & Search**
    - [x] `GET /boards/:id/search?search=` → search tickets by title and ticket number
    - [x] `GET /boards/:id/filter?assignee=&label=&status=` → filter tickets by assignee, label, or status

- **Checklist & Progress**
    - [ ] `PATCH /tickets/:id/checklist/:itemId` → toggle checklist item
    - [ ] backend calculates checklist progress
    - [ ] column stats: (#tickets / wipLimit), % done for Build/QA

- **Dashboard**
    - [ ] `GET /boards/:id/dashboard` → return project title, dueDate, prototype links (desktop, mobile, tablet), team list with avatars

- **Offline / Sync (stretch)**
    - [ ] `GET /boards/:id/sync?since=timestamp` → fetch delta changes
    - [ ] resolve conflicts on push (last-write-wins or merge)

## API Implementation Notes

- **Authentication**: All APIs require JWT token in Authorization header
- **Access Control**: Users must be board members or owners to access board resources
- **Ticket Numbering**: Auto-generated with format #XXXXXX (6-digit padded)
- **Fractional Ordering**: Used for columns and tickets with 1000-unit spacing
- **Response Format**: Consistent JSON structure with status codes and messages
- **Entity Relationships**: 
  - Boards have owners and members
  - Columns belong to boards and contain tickets
  - Tickets have labels, assignees, and checklists
  - Labels belong to boards and can be assigned to multiple tickets

## Database Schema

- **users**: id, full_name, created_at, updated_at
- **boards**: id, title, description, owner_id, archived, created_at, updated_at
- **members**: id, user_id, board_id, role, created_at, updated_at
- **columns**: id, board_id, title, position, wip_limit, created_at, updated_at
- **tickets**: id, ticket_number, title, description, column_id, position, priority, due_date, status, completed, created_at, updated_at
- **labels**: id, board_id, name, color, description, created_at, updated_at
- **ticket_labels**: id, ticket_id, label_id, created_at, updated_at
- **ticket_assignees**: id, ticket_id, member_id, created_at, updated_at
- **checklists**: id, ticket_id, title, order, status, created_at, updated_at

## Next Implementation Priorities

1. **Complete CRUD operations** for labels, assignees, and checklists
2. **Ticket movement API** with proper position management
3. **Search and filtering** functionality
4. **WIP limit enforcement** in columns
5. **Soft delete for tickets** with undo capability
6. **Progress tracking** for checklists and columns
