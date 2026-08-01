ALTER TABLE wishes
    ADD COLUMN IF NOT EXISTS group_name VARCHAR(100) NOT NULL DEFAULT 'Общее',
    ADD COLUMN IF NOT EXISTS group_color VARCHAR(20) NOT NULL DEFAULT 'slate';

UPDATE wishes
SET group_name = 'Общее'
WHERE group_name IS NULL OR group_name = '';

UPDATE wishes
SET group_color = 'slate'
WHERE group_color IS NULL OR group_color = '';
