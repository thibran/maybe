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

CREATE TRIGGER event_limit_per_folder AFTER INSERT
ON event BEGIN
    DELETE FROM event WHERE eventid IN (
        SELECT eventid FROM event
        WHERE folderref = new.folderref
        ORDER BY time DESC LIMIT -1 OFFSET 2
    );
END;



INSERT INTO folder (path) VALUES ("/tmp/test");

INSERT INTO event (time, folderref) VALUES ("2014-09-09 14:00:00", 1);

UPDATE folder SET
    c = (SELECT c+1 FROM folder WHERE path = "/tmp/test")
WHERE path = "/tmp/test";

/*
eventid     time                 folderref 
----------  -------------------  ----------
1           2015-11-28 14:43:33  1         
2           2014-09-09 14:00:00  1         
3           2015-11-28 14:43:33  1  
*/


--SELECT * FROM folder;
-- SELECT * FROM event;

/*
SELECT * FROM event
WHERE folderref = 1
ORDER BY time DESC LIMIT -1 OFFSET 2;
*/

/*
DELETE FROM event WHERE eventid IN (
    SELECT eventid FROM event
    WHERE folderref = 1
    ORDER BY time DESC LIMIT -1 OFFSET 2
);
*/

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