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


--INSERT INTO folder (path) VALUES ("/tmp/test");

INSERT OR REPLACE INTO folder (folderid, path, c)
SELECT old.folderid, new.path, old.c
FROM ( SELECT "/tmp/test" AS path
) AS new
LEFT JOIN ( SELECT
    folderid,
    path,
    (c+1) AS c
    FROM folder WHERE path = "/tmp/test"
) AS old ON new.path = old.path
;

SELECT * FROM folder;
SELECT * FROM event;





/*
SELECT new.folderid, path p, count c FROM (
    SELECT * FROM folder WHERE path = "/tmp/foo"
) AS new
WHERE path = "/tmp/foo";
*/

/*
INSERT OR REPLACE INTO folder (path) VALUES (
    
);
*/

/*
INSERT INTO folder (path) VALUES ("/tmp/test");
INSERT INTO folder (path) VALUES ("/home/tux");

SELECT * FROM folder;
SELECT * FROM event WHERE folderref = 1;
*/

/*
--INSERT INTO event (folderref) VALUES (1);
--INSERT INTO event (folderref) VALUES (1);

DELETE FROM folder WHERE folderid = 1;

SELECT * FROM folder;
SELECT * FROM event;
*/

/*
INSERT OR REPLACE INTO folder (folderid, path, c)
SELECT old.folderid, old.path, new.c
FROM ( SELECT
    "/tmp/test" AS path,
    10          AS c
) AS new
LEFT JOIN (
    SELECT folderid, path FROM folder
) AS old ON new.path = old.path
;
*/