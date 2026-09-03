UPDATE boletas_rifas
SET dada_a = NULL;

/*
	Exportar las boletas en un archivo csv 
*/

create table boletas_rifas (
  codigo_rifa integer,
  codigo_boleta integer,
  boleta text,
  entrega integer,

  foreign key (entrega) references entregas_de_boletas(codigo),
  foreign key (codigo_rifa) references datos_de_rifas(codigo),
  unique(codigo_rifa, boleta)
);

UPDATE boletas_rifas
SET entrega = NULL;

-- hacemos la tabla de entregas para conectar todo

CREATE TABLE entregas_de_boletas (
  codigo integer primary key autoincrement,
  entrega_a integer,
  en_rifa integer,
  fecha date default (date('now', 'localtime')),
	
  FOREIGN KEY (entrega_a) REFERENCES usuarios(id),
  FOREIGN KEY (en_rifa) REFERENCES datos_de_rifas(codigo)
);

/*
  CORREGIMOS LA NO EXISTENCIA DE UNA LLAVE FORANEA EN LA TABLA DE CHEQUES DE RIFAS
*/
