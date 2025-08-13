-- Add missing columns to existing tables

-- Add user_full_name to task_assignments table
ALTER TABLE task_assignments ADD COLUMN user_full_name TEXT NOT NULL DEFAULT '';

-- Add position and ticket_no columns to tickets table (if not already added)
DO $$
BEGIN
    -- Add position column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'tickets' AND column_name = 'position') THEN
        ALTER TABLE tickets ADD COLUMN position INTEGER DEFAULT 0;
    END IF;

    -- Add ticket_no column if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'tickets' AND column_name = 'ticket_no') THEN
        ALTER TABLE tickets ADD COLUMN ticket_no TEXT UNIQUE NOT NULL DEFAULT '';

        -- Update existing tickets with generated ticket numbers
        DECLARE
            ticket_record RECORD;
            counter INTEGER := 1;
        BEGIN
            FOR ticket_record IN
                SELECT id FROM tickets ORDER BY created_at
            LOOP
                UPDATE tickets
                SET ticket_no = 'TASK-' || LPAD(counter::TEXT, 4, '0')
                WHERE id = ticket_record.id;
                counter := counter + 1;
            END LOOP;
        END;
    END IF;
END $$;
