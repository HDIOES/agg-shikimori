
-- +migrate Up
-- +migrate StatementBegin
DO $$
DECLARE 
 anime_cursor CURSOR FOR SELECT id, russian FROM anime;
 rec RECORD;
 words VARCHAR(50)[];
 ngramm VARCHAR(3);
BEGIN  
 OPEN anime_cursor;
  LOOP
    FETCH FROM anime_cursor INTO rec;
    EXIT WHEN NOT FOUND;
    IF rec.russian IS NULL THEN
	 CONTINUE;
	END IF;
	words := regexp_split_to_array(rec.russian, ' ');
	FOR w_i IN 1 .. array_length(words, 1)
	 LOOP
      FOR i in 1 .. length(words[w_i]) - 2
       LOOP
        ngramm := substring(words[w_i], i, 3)::varchar(3);
        BEGIN
         INSERT INTO ngramm (ngramm_value, anime_id) VALUES(lower(ngramm), rec.id);
         EXCEPTION WHEN OTHERS THEN 
          BEGIN
           raise notice 'error occurs with row with id = %. anime.name = %. ngramm = %', rec.id, rec.russian, ngramm;
          END;
        END; 
       END LOOP;
	 END LOOP;
  END LOOP;
 CLOSE anime_cursor; 
END $$;
-- +migrate StatementEnd
-- +migrate Down
