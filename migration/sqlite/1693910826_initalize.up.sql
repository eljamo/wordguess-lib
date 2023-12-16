CREATE TABLE recent_word_log (
    'date' TEXT PRIMARY KEY,
    'word' TEXT NOT NULL UNIQUE
) STRICT;
