
-- +migrate Up
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION get_rank(rus_tsvector TSVECTOR, eng_tsvector TSVECTOR, ts_query TSQUERY) RETURNS REAL AS $$
 DECLARE
  russian_rank REAL;
  eng_rank REAL;
 BEGIN
  russian_rank := ts_rank(rus_tsvector, ts_query);
  eng_rank := ts_rank(eng_tsvector, ts_query);
  IF russian_rank > eng_rank THEN
   RETURN russian_rank;
  ELSE
   RETURN eng_rank;
  END IF;
 END;
$$
LANGUAGE PLPGSQL;
-- +migrate StatementEnd
-- +migrate Down
