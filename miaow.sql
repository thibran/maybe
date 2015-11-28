PRAGMA foreign_keys = ON;

CREATE TABLE folder (
    folderid INTEGER PRIMARY KEY,
    path TEXT UNIQUE NOT NULL,
    c INTEGER DEFAULT 1 NOT NULL
);

CREATE TABLE event(
    eventid INTEGER PRIMARY KEY,
    time DATETIME DEFAULT CURRENT_TIMESTAMP,
    folderref INTEGER NOT NULL,
    FOREIGN KEY(folderref) REFERENCES folder(folderid) ON DELETE CASCADE
);

CREATE TRIGGER insert_folder AFTER INSERT
ON folder BEGIN
    INSERT INTO event (folderref) VALUES (new.folderid);
END;

CREATE TRIGGER update_folder AFTER UPDATE
ON folder BEGIN
    INSERT INTO event (folderref) VALUES (new.folderid);
END;


INSERT INTO folder (path) VALUES ("/tmp/test");
UPDATE folder SET
    c = (SELECT c+1 FROM folder WHERE path = "/tmp/test")
WHERE path = "/tmp/test";

UPDATE folder SET
    c = (SELECT c+1 FROM folder WHERE path = "/tmp/test")
WHERE path = "/tmp/test";

SELECT * FROM folder;
SELECT * FROM event;


/*
-- deletes old value before replacing -> bad
INSERT OR REPLACE INTO folder (folderid, path, c)
SELECT old.folderid, new.path, old.c
FROM ( SELECT "/tmp/test" AS path
) AS new
LEFT JOIN ( SELECT
    folderid,
    path,
    (c+1) AS c
    FROM folder WHERE path = "/tmp/test"
) AS old ON new.path = old.path;
*/