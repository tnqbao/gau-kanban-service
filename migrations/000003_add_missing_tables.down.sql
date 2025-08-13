-- Drop the tables and sequence created in the up migration

DROP INDEX IF EXISTS idx_checklists_position;
DROP INDEX IF EXISTS idx_checklists_ticket_id;
DROP TABLE IF EXISTS checklists;

DROP INDEX IF EXISTS idx_task_assignments_user_id;
DROP INDEX IF EXISTS idx_task_assignments_ticket_id;
DROP TABLE IF EXISTS task_assignments;

DROP SEQUENCE IF EXISTS ticket_number_seq;
