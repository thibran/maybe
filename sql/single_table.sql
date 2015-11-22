CREATE TABLE d (
       _id INTEGER PRIMARY KEY,        -- biggest ID = last used directory
       p TEXT UNIQUE NOT NULL,         -- path
       c INTEGER DEFAULT 1 NOT NULL,   -- counter
       dt DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX index_path ON d (p);

CREATE TRIGGER t_max_rows AFTER INSERT
ON d BEGIN
	DELETE FROM d WHERE (SELECT COUNT(_id) FROM d) = 30000
	       AND _id = (SELECT _id FROM d LIMIT 1);
END;


INSERT OR REPLACE INTO d (p, dt) VALUES ('/home/tux', DATETIME('now', '-10 Hour', '-30 Minute'));
INSERT OR REPLACE INTO d (p, dt) VALUES ('/tmp', DATETIME('now', '-3 Minute', '-26 Second'));
INSERT OR REPLACE INTO d (p, dt) VALUES ('/usr', DATETIME('now', '-1 Day', '-6 Hour', '-13 Second'));


INSERT OR REPLACE INTO d (p, dt, c) VALUES ('/etc', DATETIME('now', '-12 Second'),
       (SELECT c+1 FROM d WHERE p = '/etc'));

INSERT OR REPLACE INTO d (p, dt, c) VALUES ('/home/tux', DATETIME('now'),
       (SELECT c+1 FROM d WHERE p = '/home/tux'));

--SELECT *, (SELECT _id FROM d ORDER BY _id DESC LIMIT 1) AS last_dir FROM d;

SELECT p FROM(
SELECT p,
(
	-- weighing last usage time
	CASE
	WHEN dt BETWEEN DATETIME('now', '-5 Minute') AND DATETIME('now') THEN +100
	WHEN dt BETWEEN DATETIME('now', '-1 Hour') AND DATETIME('now') THEN +50
	WHEN dt BETWEEN DATETIME('now', '-3 Month') AND DATETIME('now') THEN +10
	WHEN dt BETWEEN DATETIME('now', '-1 Year') AND DATETIME('now') THEN +1
	ELSE 0
	END
	+
	-- last used dir
	CASE WHEN _id = (SELECT _id FROM d ORDER BY _id DESC LIMIT 1) THEN +50 ELSE 0 END
	+
	-- incorprate how often dir was opened
	c * 10
	) AS score
FROM d
WHERE p LIKE '%t%'
ORDER BY score DESC, dt, _id
LIMIT 1
);


-- SELECT
-- *
-- FROM(
-- SELECT
-- *,
-- (
-- 	-- weighing last usage time
-- 	CASE
-- 	WHEN dt BETWEEN DATETIME('now', '-5 Minute') AND DATETIME('now') THEN +100
-- 	WHEN dt BETWEEN DATETIME('now', '-1 Hour') AND DATETIME('now') THEN +50
-- 	WHEN dt BETWEEN DATETIME('now', '-3 Month') AND DATETIME('now') THEN +10
-- 	WHEN dt BETWEEN DATETIME('now', '-1 Year') AND DATETIME('now') THEN +1
-- 	ELSE 0
-- 	END
-- 	+
-- 	-- last used dir
-- 	CASE WHEN _id = (SELECT _id FROM d ORDER BY _id DESC LIMIT 1) THEN +50 ELSE 0 END
-- 	+
-- 	-- incorprate how often dir was opened
-- 	c * 10
-- 	) AS score
-- FROM d
-- ORDER BY score DESC, dt, _id
-- )
-- ;
