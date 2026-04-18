CREATE TABLE IF NOT EXISTS icd_codes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT NOT NULL,
    description TEXT NOT NULL,
    version TEXT NOT NULL
);

-- Indexing for performance
CREATE INDEX idx_icd_codes_code ON icd_codes(code);
CREATE INDEX idx_icd_codes_description ON icd_codes(description);
