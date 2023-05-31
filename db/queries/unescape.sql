--- Undo some cases of overzealous escaping.
UPDATE interpretations SET why = replace(why, "&#39;", "'") where why like '%&#39;%';
