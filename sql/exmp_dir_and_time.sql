CREATE TABLE d (_id INTEGER PRIMARY KEY, p TEXT UNIQUE NOT NULL);

CREATE TRIGGER t_dir_max_rows BEFORE INSERT
ON d BEGIN
	DELETE FROM d WHERE (SELECT COUNT(_id) FROM d) = 30000
	       AND _id = (SELECT _id FROM d LIMIT 1);
END;

/*
CREATE TRIGGER t_dir_insert AFTER INSERT
ON d BEGIN
	INSERT INTO t (dir_id) VALUES (new._id);
END;

CREATE TRIGGER t_dir_update AFTER UPDATE
ON d BEGIN
	INSERT INTO t (dir_id) VALUES (new._id);
END;
*/

CREATE TRIGGER t_dir_delete BEFORE DELETE
ON d BEGIN
   DELETE FROM t WHERE dir_id = old._id;
END;

CREATE TABLE t (
       _id INTEGER PRIMARY KEY,
       dir_id INTEGER NOT NULL,
       dt DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX i_time_doc_id ON t (dir_id);

/*
TODO:
- Limit t to 20 entries for every directory
- Add count column to d and increase it in t_dir_update
      UPDATE numbers SET num = num + 1;
*/

INSERT INTO d (_id, p) VALUES (1, '/home/tux');
INSERT INTO t (dir_id, dt) VALUES (1, DATETIME('now'));
INSERT INTO t (dir_id, dt) VALUES (1, '2014-11-15 18:07:00');

INSERT INTO d (_id, p) VALUES (2, '/tmp');
INSERT INTO t (dir_id, dt) VALUES (2, DATETIME('now', '-30 Minute'));


SELECT
	p,
	score,
	DATETIME('now') AS now
FROM(SELECT
	d._id,
	d.p,
	MAX((
	CASE
	WHEN dt BETWEEN DATETIME('now', '-5 Minute') AND DATETIME('now') THEN +100
	WHEN dt BETWEEN DATETIME('now', '-1 Hour') AND DATETIME('now') THEN +50
	WHEN dt BETWEEN DATETIME('now', '-3 Month') AND DATETIME('now') THEN +10
	WHEN dt BETWEEN DATETIME('now', '-1 Year') AND DATETIME('now') THEN +1
	ELSE 0
	END
	+ (SELECT COUNT(t._id) FROM t WHERE d._id = t.dir_id) * 10
	)) AS score,
	(SELECT COUNT(t._id) FROM t WHERE d._id = t.dir_id) AS c
FROM d
JOIN t ON d._id = t.dir_id
WHERE d.p LIKE '%t%'
GROUP BY p
ORDER BY score DESC, dt, d._id
)
;






/*
SELECT
	d._id,
	d.p,
	(
	CASE
	WHEN dt BETWEEN DATETIME('now', '-5 Minute') AND DATETIME('now') THEN +100
	WHEN dt BETWEEN DATETIME('now', '-1 Hour') AND DATETIME('now') THEN +50
	WHEN dt BETWEEN DATETIME('now', '-3 Month') AND DATETIME('now') THEN +10
	WHEN dt BETWEEN DATETIME('now', '-1 Year') AND DATETIME('now') THEN +1
	ELSE 0
	END
	+
	(SELECT COUNT(t._id) FROM t WHERE d._id = t.dir_id) * 10
	) AS score,
	(SELECT COUNT(t._id) FROM t WHERE d._id = t.dir_id) AS c,
	DATETIME('now') AS now
FROM d
JOIN t ON d._id = t.dir_id
WHERE d.p LIKE '%t%'

ORDER BY score DESC, dt
;
*/

/*
INSERT INTO t (dir_id, dt) VALUES (3, '2014-09-09 14:00:00');
INSERT INTO t (dir_id, dt) VALUES (3, '2014-09-09 14:00:00');

INSERT INTO d (p) VALUES ("/home/tux");
INSERT INTO d (p) VALUES ("/tmp");
INSERT INTO d (p) VALUES ("/usr");
*/


/*
SELECT * FROM d;
SELECT dir_id, dt FROM t ORDER BY dt DESC;
*/

