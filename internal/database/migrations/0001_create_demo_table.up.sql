BEGIN TRANSACTION;
CREATE TABLE demo (
    key text NOT NULL ,
    val jsonb NOT NULL
);
COMMIT;