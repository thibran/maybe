CREATE TABLE d (_id INTEGER PRIMARY KEY, str TEXT);

CREATE TRIGGER t_d_max_rows BEFORE INSERT
ON d
BEGIN
DELETE FROM d WHERE (SELECT COUNT(_id) FROM d) = 2
       AND _id = (SELECT _id FROM d LIMIT 1);
END;


INSERT INTO d (str) VALUES ("a");
INSERT INTO d (str) VALUES ("b");
INSERT INTO d (str) VALUES ("c");


/*
SELECT * FROM d;
SELECT MIN(_id) FROM d;
SELECT _id FROM d LIMIT 1;
DELETE FROM d WHERE _id = (SELECT _id FROM d LIMIT 1);
*/


SELECT * FROM d;

/*
DELETE FROM d WHERE (SELECT COUNT(_id) FROM d) > 2
AND _id = (SELECT _id FROM d LIMIT 1);

SELECT * FROM d WHERE (SELECT COUNT(_id) FROM d) > 2 LIMIT 1;
SELECT _id FROM d LIMIT 1;
DELETE FROM d WHERE _id = (SELECT MIN(_id) FROM d);
SELECT * FROM d;
*/