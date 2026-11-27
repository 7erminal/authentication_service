-- PROPOSAL ONLY - NOT APPLIED. Do not copy into migrations/ without approval.
--
-- Adds nullable tracing columns to the 4 token tables so a token row can be
-- tied back to the request that created it. All columns are NULL-able with no
-- default, so this is backwards compatible: existing rows stay NULL and current
-- INSERTs (which do not name these columns) keep working.
--
-- No primary key changes, no unique constraints, no indexes (add an index on
-- request_id later only if lookups by request ID become a real need).
-- Values are the sanitized IDs ([A-Za-z0-9_-]{8,64}), so VARCHAR(64) is enough.
--
-- To adopt with Atlas (after approval):
--   1. add the fields to the beego models (orm:"null;size(64)") so ./loader
--      reflects them,
--   2. `atlas migrate diff add_trace_columns --env beego` (regenerates this
--      and updates migrations/atlas.sum),
--   3. review, then `atlas migrate apply` in a maintenance window.
-- On MySQL 8 these are ALGORITHM=INSTANT column additions (no table rebuild).

ALTER TABLE `access_tokens`
  ADD COLUMN `trace_id`   varchar(64) NULL,
  ADD COLUMN `request_id` varchar(64) NULL,
  ADD COLUMN `session_id` varchar(64) NULL;

ALTER TABLE `refresh_tokens`
  ADD COLUMN `trace_id`   varchar(64) NULL,
  ADD COLUMN `request_id` varchar(64) NULL,
  ADD COLUMN `session_id` varchar(64) NULL;

ALTER TABLE `customer_access_tokens`
  ADD COLUMN `trace_id`   varchar(64) NULL,
  ADD COLUMN `request_id` varchar(64) NULL,
  ADD COLUMN `session_id` varchar(64) NULL;

ALTER TABLE `customer_refresh_tokens`
  ADD COLUMN `trace_id`   varchar(64) NULL,
  ADD COLUMN `request_id` varchar(64) NULL,
  ADD COLUMN `session_id` varchar(64) NULL;

-- Rollback:
-- ALTER TABLE `access_tokens`           DROP COLUMN `trace_id`, DROP COLUMN `request_id`, DROP COLUMN `session_id`;
-- ALTER TABLE `refresh_tokens`          DROP COLUMN `trace_id`, DROP COLUMN `request_id`, DROP COLUMN `session_id`;
-- ALTER TABLE `customer_access_tokens`  DROP COLUMN `trace_id`, DROP COLUMN `request_id`, DROP COLUMN `session_id`;
-- ALTER TABLE `customer_refresh_tokens` DROP COLUMN `trace_id`, DROP COLUMN `request_id`, DROP COLUMN `session_id`;
