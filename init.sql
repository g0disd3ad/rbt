CREATE TABLE IF NOT EXISTS translations (
    word_eng VARCHAR(255) NOT NULL,
    word_rus VARCHAR(255) NOT NULL,
    PRIMARY KEY (word_eng, word_rus)
);