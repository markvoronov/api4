-- В up-миграции переводим dateAdd из TIME в TIMESTAMPTZ
ALTER TABLE url5
ALTER COLUMN dateAdd TYPE timestamptz
  USING now();