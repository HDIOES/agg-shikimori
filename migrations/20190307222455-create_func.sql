
-- +migrate Up
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION MIN3(m1 INT, m2 INT, m3 INT) RETURNS INT AS $$
 BEGIN
  IF m1 > m2 THEN
   IF m3 > m2 THEN
    RETURN m2;
   ELSE
    RETURN m3;
   END IF;
  ELSE
   IF m1 > m3 THEN
    RETURN m3;
   ELSE
    RETURN m1;
   END IF;
  END IF;
 END;
$$
LANGUAGE PLPGSQL;
-- +migrate StatementEnd
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION LOWENSTAIN_DISTANCE(S1 VARCHAR(255), S2 VARCHAR(255)) RETURNS INT AS $$
DECLARE 
 countOfLetters INT;
 mas INT[][];
 length1 INT;
 length2 INT;
 chararray1 VARCHAR(1)[];
 chararray2 VARCHAR(1)[];
 i INT; 
 j INT;
BEGIN
 length1 := length(S1);
 length2 := length(S2);
 mas := array_fill(0, ARRAY[length1, length2]);
 FOR i IN 1 .. length1
 LOOP
    chararray1 := array_append(chararray1, SUBSTRING(S1, i, 1)::varchar(1));
 END LOOP;
 FOR i IN 1 .. length2
 LOOP
    chararray2 := array_append(chararray2, SUBSTRING(S2, i, 1)::varchar(1));
 END LOOP;
 FOR i IN 1 .. length1
 LOOP
  FOR j IN 1 .. length2
  LOOP
   IF i = 1 AND j = 1 THEN
    mas[i][j] := 0;
   END IF;
   IF i > 1 AND j = 1 THEN
    mas[i][j] := i - 1;
   END IF;
   IF j > 1 AND i = 1 THEN
    mas[i][j] := j - 1;
   END IF;
   IF i > 1 AND j > 1 THEN 
    IF chararray1[i] = chararray2[j] THEN
     mas[i][j] := mas[i - 1][j - 1];
    ELSE
     mas[i][j] := MIN3(mas[i][j - 1], mas[i - 1][j - 1], mas[i - 1][j]) + 1;
    END IF;
   END IF;
  END LOOP;
 END LOOP;
 RETURN mas[length1][length2];
END;
$$
LANGUAGE PLPGSQL;
-- +migrate StatementEnd
-- +migrate Down
