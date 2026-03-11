-- Add tail_number column to accounts table
ALTER TABLE accounts ADD COLUMN tail_number VARCHAR(10);

-- Update existing NULL values to empty string
UPDATE accounts SET tail_number = '' WHERE tail_number IS NULL;

-- Create index for tail_number (optional but helpful for filtering)
CREATE INDEX idx_accounts_tail_number ON accounts(tail_number);
