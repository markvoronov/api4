-- В down-миграции возвращаем обратно в TIME
ALTER TABLE url5
ALTER COLUMN dateAdd TYPE time
  USING to_char(dateAdd, 'HH24:MI:SS')::time;