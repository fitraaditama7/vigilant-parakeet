CREATE OR REPLACE FUNCTION update_timestamps()
  RETURNS TRIGGER
  LANGUAGE plpgsql
  AS $$
BEGIN
      new.updated_at = timezone('utc', now());
RETURN new;
END;
  $$;
