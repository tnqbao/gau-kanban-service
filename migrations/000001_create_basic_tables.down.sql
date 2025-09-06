-- Drop triggers
DROP TRIGGER IF EXISTS update_checklists_updated_at ON checklists;
DROP TRIGGER IF EXISTS update_ticket_assignees_updated_at ON ticket_assignees;
DROP TRIGGER IF EXISTS update_ticket_labels_updated_at ON ticket_labels;
DROP TRIGGER IF EXISTS update_labels_updated_at ON labels;
DROP TRIGGER IF EXISTS update_tickets_updated_at ON tickets;
DROP TRIGGER IF EXISTS update_columns_updated_at ON columns;
DROP TRIGGER IF EXISTS update_members_updated_at ON members;
DROP TRIGGER IF EXISTS update_boards_updated_at ON boards;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_checklists_order;
DROP INDEX IF EXISTS idx_checklists_ticket_id;
DROP INDEX IF EXISTS idx_ticket_assignees_member_id;
DROP INDEX IF EXISTS idx_ticket_assignees_ticket_id;
DROP INDEX IF EXISTS idx_ticket_labels_label_id;
DROP INDEX IF EXISTS idx_ticket_labels_ticket_id;
DROP INDEX IF EXISTS idx_labels_board_id;
DROP INDEX IF EXISTS idx_tickets_due_date;
DROP INDEX IF EXISTS idx_tickets_priority;
DROP INDEX IF EXISTS idx_tickets_status;
DROP INDEX IF EXISTS idx_tickets_position;
DROP INDEX IF EXISTS idx_tickets_column_id;
DROP INDEX IF EXISTS idx_columns_position;
DROP INDEX IF EXISTS idx_columns_board_id;
DROP INDEX IF EXISTS idx_members_user_board;
DROP INDEX IF EXISTS idx_members_board_id;
DROP INDEX IF EXISTS idx_members_user_id;
DROP INDEX IF EXISTS idx_boards_archived;
DROP INDEX IF EXISTS idx_boards_owner_id;

-- Drop tables in reverse order of creation (due to foreign key constraints)
DROP TABLE IF EXISTS checklists;
DROP TABLE IF EXISTS ticket_assignees;
DROP TABLE IF EXISTS ticket_labels;
DROP TABLE IF EXISTS labels;
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS columns;
DROP TABLE IF EXISTS members;
DROP TABLE IF EXISTS boards;
DROP TABLE IF EXISTS users;
