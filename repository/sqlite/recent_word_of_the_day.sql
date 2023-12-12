-- query: FindAll
SELECT "date", "word" FROM recent_word_of_the_day;

-- query: FindByDate
SELECT "date", "word" FROM recent_word_of_the_day WHERE "date" = ?;

-- query: InsertWord
INSERT INTO recent_word_of_the_day ("date", "word") VALUES (?, ?);

-- query: DeleteByDate
DELETE FROM recent_word_of_the_day WHERE "date" = ?;
