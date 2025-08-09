ALTER TABLE tickets ADD COLUMN ticket_no TEXT UNIQUE NOT NULL DEFAULT '';

-- Update existing tickets with generated ticket numbers
DO $$
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
END $$;
