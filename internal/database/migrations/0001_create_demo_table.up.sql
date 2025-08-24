BEGIN TRANSACTION;
CREATE TABLE demo (
    key text,
    val jsonb
);
COMMIT;