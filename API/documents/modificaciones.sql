ALTER TABLE cheques_de_prestamos ADD COLUMN tipo TEXT;
ALTER TABLE cheques_de_prestamos ADD COLUMN intereses_antes INTEGER;
ALTER TABLE cheques_de_prestamos ADD COLUMN intereses_despues INTEGER;

UPDATE cheques_de_prestamos
SET 
	tipo = "DEUDA",
	intereses_antes = 0,
	intereses_despues = 0

  

