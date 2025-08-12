-- Add missing tables and columns

-- Create task_assignments table if not exists
CREATE TABLE IF NOT EXISTS task_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    full_name TEXT NOT NULL,
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_task_assignments_ticket_id ON task_assignments(ticket_id);
CREATE INDEX IF NOT EXISTS idx_task_assignments_user_id ON task_assignments(user_id);

-- Create checklists table if not exists
CREATE TABLE IF NOT EXISTS checklists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    completed BOOLEAN DEFAULT false,
    position INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);

-- Create index for better performance
CREATE INDEX IF NOT EXISTS idx_checklists_ticket_id ON checklists(ticket_id);
CREATE INDEX IF NOT EXISTS idx_checklists_position ON checklists(ticket_id, position);

-- Add sequence for ticket numbering if not exists
CREATE SEQUENCE IF NOT EXISTS ticket_number_seq START 1;

-- Update existing tickets to have proper ticket numbers if they don't have them
UPDATE tickets
SET ticket_no = 'TASK-' || LPAD(nextval('ticket_number_seq')::text, 4, '0')
WHERE ticket_no IS NULL OR ticket_no = '';
