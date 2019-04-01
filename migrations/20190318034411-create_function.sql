
-- +migrate Up
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION countOfNgramm(s1 VARCHAR(255), n INT) RETURNS INT AS $$
DECLARE 
 words VARCHAR(100)[];
 count INT := 0;
 word_length INT;
BEGIN
words := regexp_split_to_array(s1, ' ');
FOR i IN 1 .. array_length(words, 1)
LOOP
 word_length := length(words[i]);
 IF word_length >= 3 THEN
  count := count + word_length - n + 1;
 END IF; 
END LOOP;
 RETURN count;
END;
$$
LANGUAGE PLPGSQL;
-- +migrate StatementEnd
-- +migrate Down
