-- Drop tail_number column from accounts table
DROP INDEX IF EXISTS idx_accounts_tail_number;
ALTER TABLE accounts DROP COLUMN IF EXISTS tail_number;
