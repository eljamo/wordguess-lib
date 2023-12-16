-- query: FindAll
SELECT "date", "word" FROM recent_word_log;

-- query: FindByDate
SELECT "date", "word" FROM recent_word_log WHERE "date" = ?;

-- query: InsertWord
INSERT INTO recent_word_log ("date", "word") VALUES (?, ?);

-- query: DeleteByDate
DELETE FROM recent_word_log WHERE "date" = ?;
