/*
    esto es una enmienda de todos mis errores hechos en la base de datos anterior y
    corregidos para que el programa funcione de mejor manera, tambien es una
    rendicion por querer hacer todo mas avanzado cuando la solicitudu del programa
    es muy esspecifica entorno a que se requiere, hice una sobre ingenieria con 
    una base de datos posgres cuando una sql ya funciona muy bien
*/

/*
    tablas de almacenamiento para el programa
*/

create table registros (
  id_registro integer primary key autoincrement,
  fecha datetime default (datetime('now', 'localtime')),
  registro text
);

create table datos_de_rifas (
  codigo integer primary key autoincrement,
  numero_de_boletas integer,
  costo_de_boletas integer,
  costos_de_administracion integer,
  ganancias_por_boleta integer
);

create table calendario(
	semana primary key,
	fecha datetime
);

/*
 * tabla primaria de usuarios
 */

create table usuarios (
  id integer primary key autoincrement,
  nombre text,
  puestos integer,
  telefono text,
  estado boolean default true,
  cuotas_pagas integer default 0,
  cuotas_deudas integer default 0,
  semanas_revisadas integer default 0,
  deuda_rifas integer default 0,
  multas_semanales integer default 0,
  multas_extras integer default 0,
  multas_por_cabeza integer default 0,
  capital_pago integer default 0,
  capital_retirado integer default 0,
  capital_proyectado integer default 0
);

/*
 * tablas dependientes (con llaves foráneas)
 */

create table anotaciones (
  codigo integer primary key autoincrement,
  id_user integer,
  anotacion text,
  monto integer,
  tipo text,
  fecha date default (date('now', 'localtime')),

  foreign key (id_user) references usuarios(id)
);

create table prestamos_hechos (
  codigo integer primary key autoincrement,
  id_user integer,
  estado boolean  default true,
  interes integer,
  intereses_vencidos integer  default 0,
  intereses_generados integer  default 0,
  deuda integer,
  monto integer,
  fechas_revisadas integer  default 0,
  fecha date default (date('now', 'localtime')),
  motivo text,
  
  foreign key (id_user) references usuarios(id)
);

create table deudas (
  codigo integer primary key autoincrement,
  id_user integer,
  monto integer,
  fecha_de_creacion datetime default (datetime('now', 'localtime')),
  fecha_de_pago datetime,
  tipo text,
  motivo text,
  estado boolean default true,

  foreign key (id_user) references usuarios(id)
);

create table transferencias (
  codigo integer primary key autoincrement,
  id_user integer,
  fecha datetime default (datetime('now', 'localtime')),
  monto integer,

  foreign key (id_user) references usuarios(id)
);

create table cheques_de_cuotas (
  codigo integer primary key autoincrement,
  id_user integer,
  cuotas_pagas integer,
  multas_pagas integer,
  metodo_de_pago text,
  fecha datetime default (datetime('now', 'localtime')),

  foreign key (id_user) references usuarios(id)
);

create table cheques_de_rifas (
  codigo integer primary key autoincrement,
  id_user integer,
  monto_pago integer,
  deuda_antes integer,
  deuda_despues integer,
  en_rifa integer,
  fecha datetime default (datetime('now', 'localtime')),

  foreign key (id_user) references usuarios(id),
  foreign key (en rifa) references datos_de_rifas(codigo)
);

create table cheques_de_prestamos (
  codigo integer primary key autoincrement,
  id_user integer,
  codigo_prestamo integer,
  monto_pago integer,
  deuda_antes integer,
  deuda_despues integer,
  fecha datetime default (datetime('now', 'localtime')),
  tipo text,
  intereses_antes integer,
  intereses_despues integer,

  foreign key (id_user) references usuarios(id),
  foreign key (codigo_prestamo) references prestamos_hechos(codigo)
);

create table notas_de_prestamos (
  codigo_prestamo integer,
  anotacion text,
  fecha date default (date('now', 'localtime')),

  foreign key (codigo_prestamo) references prestamos_hechos(codigo)
);

-- create table boletas_rifas (
--   codigo_rifa integer,
--   codigo_boleta integer,
--   boleta text,
--   dada_a integer,

--   foreign key (dada_a) references usuarios(id),
--   foreign key (codigo_rifa) references datos_de_rifas(codigo),
--   unique(codigo_rifa, boleta)
-- );

create table boletas_rifas (
  codigo_rifa integer,
  codigo_boleta integer,
  boleta text,
  entrega integer,

  foreign key (entrega) references entregas_de_boletas(codigo),
  foreign key (codigo_rifa) references datos_de_rifas(codigo),
  unique(codigo_rifa, boleta)
);

CREATE TABLE entregas_de_boletas (
  codigo integer primary key autoincrement,
  entrega_a integer,
  en_rifa integer,
  fecha date default (date('now', 'localtime')),
	
  FOREIGN KEY (entrega_a) REFERENCES usuarios(id),
  FOREIGN KEY (en_rifa) REFERENCES datos_de_rifas(codigo)
);

/*
 * tablas para normalizar
 */

create table cuotas_bloqueos (
  id_user integer,
  semana integer,

	foreign key (id_user) references usuarios(id),
	unique(id_user, semana)
);

create table cuotas_multas_activas (
  id_user integer,
  semana integer,
  valor integer,

	foreign key (id_user) references usuarios(id),
  unique(id_user, semana)
);

create table cuotas_multas_pagas (
  id_user integer,
  semana integer,
  valor integer,

	foreign key (id_user) references usuarios(id),
  unique(id_user, semana)
);

create table premios_de_rifas (
  codigo_rifa integer,
  posicion integer,
  premio integer,

	foreign key (codigo_rifa) references datos_de_rifas(codigo),
  unique(codigo_rifa, posicion)
);

create table prestamos_fiadores (
  id_user_fiador integer,
  codigo_prestamo integer,
  monto_respaldado integer,

  foreign key (id_user_fiador) references usuarios(id),
  foreign key (codigo_prestamo) references prestamos_hechos(codigo)
);

create table prestamos_fechas_de_pago (
  codigo_prestamo integer,
  fecha date,
  
  foreign key (codigo_prestamo) references prestamos_hechos(codigo)
);

/*
 * tabla para super usuario (1976)
 */

create table super_usuario (
  codigo_prestamo integer unique,
  monto_respaldado integer,

  foreign key (codigo_prestamo) references prestamos_hechos(codigo)
);
