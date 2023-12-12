CREATE TABLE recent_word_of_the_day (
    'date' TEXT PRIMARY KEY,
    'word' TEXT NOT NULL UNIQUE
) STRICT;
