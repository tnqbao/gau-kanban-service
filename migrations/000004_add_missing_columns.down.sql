-- Rollback missing columns

-- Remove user_full_name from task_assignments table
ALTER TABLE task_assignments DROP COLUMN IF EXISTS user_full_name;

-- Remove position from tickets table
ALTER TABLE tickets DROP COLUMN IF EXISTS position;

-- Remove ticket_no from tickets table
ALTER TABLE tickets DROP COLUMN IF EXISTS ticket_no;
