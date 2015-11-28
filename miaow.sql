/*
PLAN:
    1. get folder by path with count
        - if not nil, update and increase count (check for max INT too)
        - if nil, create new folder entry
    2. test how big the db gets with 30.000 folder entries
*/

PRAGMA foreign_keys = ON;
--PRAGMA recursive_triggers = ON;

CREATE TABLE folder (
    folderid INTEGER PRIMARY KEY,
    path TEXT UNIQUE NOT NULL,
    c INTEGER DEFAULT 1 NOT NULL --count
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

-- keep 30.000 folders
CREATE TRIGGER folder_limit AFTER INSERT
ON folder BEGIN
    DELETE FROM folder WHERE folderid IN (
        SELECT folderid FROM folder
        ORDER BY c DESC LIMIT -1 OFFSET 30000
    );
END;

-- keep 10 events per folder entry
CREATE TRIGGER event_limit_per_folder AFTER INSERT
ON event BEGIN
    DELETE FROM event WHERE eventid IN (
        SELECT eventid FROM event
        WHERE folderref = new.folderref
        ORDER BY time DESC LIMIT -1 OFFSET 10
    );
END;

/*
UPDATE folder SET
    c = (SELECT c+1 FROM folder WHERE path = "/tmp/test")
WHERE path = "/tmp/test";
*/

INSERT INTO folder (path, c) VALUES ("/tmp/test", 2);
INSERT INTO folder (path, c) VALUES ("/etc/apt", 1);
INSERT INTO event (time, folderref) VALUES (DATETIME("now", "-10 Minute"), 2);

INSERT INTO folder (path, c) VALUES ("/home/tux", 3);

INSERT INTO event (time, folderref) VALUES ("2014-09-09 14:00:00", 1);
INSERT INTO event (time, folderref) VALUES (DATETIME("now", "-30 Minute"), 1);
    
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