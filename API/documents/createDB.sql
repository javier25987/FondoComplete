CREATE TABLE ajustes (
    ajuste TEXT UNIQUE NOT NULL,
    valor_n INTEGER NOT NULL,
    valor_t TEXT NOT NULL
);

CREATE TABLE anotaciones (
    codigo INTEGER PRIMARY KEY AUTOINCREMENT,
    id_user INTEGER NOT NULL,
    anotacion TEXT NOT NULL,
    monto INTEGER NOT NULL,
    tipo TEXT NOT NULL,
    fecha TEXT NOT NULL
);

CREATE TABLE boletas_rifa_1 (
    codigo INTEGER NOT NULL,
    boleta TEXT UNIQUE NOT NULL,
    dada_a INTEGER NOT NULL
);

CREATE TABLE boletas_rifa_2 (
    codigo INTEGER NOT NULL,
    boleta TEXT UNIQUE NOT NULL,
    dada_a INTEGER NOT NULL
);

CREATE TABLE boletas_rifa_3 (
    codigo INTEGER NOT NULL,
    boleta TEXT UNIQUE NOT NULL,
    dada_a INTEGER NOT NULL
);

CREATE TABLE boletas_rifa_4 (
    codigo INTEGER NOT NULL,
    boleta TEXT UNIQUE NOT NULL,
    dada_a INTEGER NOT NULL
);

CREATE TABLE cuotas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pagas INTEGER NOT NULL,
    adeudas INTEGER NOT NULL,
    multas_activas TEXT NOT NULL,
    multas_pagas TEXT NOT NULL,
    revisiones INTEGER NOT NULL,
    bloqueos TEXT NOT NULL
);

CREATE TABLE datos_de_rifas (
    codigo INTEGER PRIMARY KEY AUTOINCREMENT,
    numero_de_boletas INTEGER NOT NULL,
    premios TEXT NOT NULL,
    costo_de_boletas INTEGER NOT NULL,
    costos_de_administracion INTEGER NOT NULL,
    ganancias_por_boleta INTEGER NOT NULL
);

CREATE TABLE deudas_rifas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    deuda INTEGER NOT NULL
);

CREATE TABLE usuarios (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nombre TEXT NOT NULL,
    puestos INTEGER NOT NULL,
    telefono TEXT NOT NULL,
    estado BOOLEAN NOT NULL
);

CREATE TABLE prestamos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    deudas_por_fiador INTEGER NOT NULL,
    fiador_de TEXT NOT NULL
);

CREATE TABLE registros (
    fecha TEXT NOT NULL,
    ingreso INTEGER NOT NULL,
    egreso INTEGER NOT NULL,
    mes INTEGER NOT NULL
);

CREATE TABLE transferencias (
    id_user INTEGER NOT NULL,
    fecha TEXT NOT NULL,
    monto INTEGER NOT NULL
);

CREATE TABLE multas (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    semanales INTEGER NOT NULL,
    extras INTEGER NOT NULL,
    por_cabeza INTEGER NOT NULL
);

CREATE TABLE capitales (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pago INTEGER NOT NULL,
    retirado INTEGER NOT NULL,
    proyectado INTEGER NOT NULL
);

CREATE TABLE prestamos_hechos (
    codigo INTEGER PRIMARY KEY AUTOINCREMENT,
    id_user INTEGER NOT NULL,
    estado BOOLEAN NOT NULL,
    interes INTEGER NOT NULL,
    intereses_vencidos INTEGER NOT NULL,
    intereses_generados INTEGER NOT NULL,
    deuda INTEGER NOT NULL,
    monto INTEGER NOT NULL,
    fechas_de_pago TEXT NOT NULL,
    revisiones INTEGER NOT NULL,
    fiadores TEXT NOT NULL,
    deudas_con_fiadores TEXT NOT NULL,
    motivo TEXT NOT NULL
);

CREATE TABLE deudas (
    codigo INTEGER PRIMARY KEY AUTOINCREMENT,
    id_user INTEGER NOT NULL,
    monto INTEGER NOT NULL,
    fecha_de_creacion TEXT NOT NULL,
    fecha_de_pago TEXT NOT NULL,
    tipo TEXT NOT NULL,
    motivo TEXT NOT NULL,
    estado BOOLEAN NOT NULL
);

CREATE TABLE cheques_de_cuotas (
    codigo INTEGER PRIMARY KEY AUTOINCREMENT,
    id_user INTEGER NOT NULL,
    puestos INTEGER NOT NULL,
    cuotas_pagas INTEGER NOT NULL,
    multas_pagas INTEGER NOT NULL,
    metodo_de_pago TEXT NOT NULL,
    fecha TEXT NOT NULL,
    hora TEXT NOT NULL
);

CREATE TABLE cheques_de_prestamos (
    codigo INTEGER PRIMARY KEY AUTOINCREMENT,
    id_user INTEGER NOT NULL,
    codigo_prestamo INTEGER NOT NULL,
    monto_pago INTEGER NOT NULL,
    deuda_antes INTEGER NOT NULL,
    deuda_despues INTEGER NOT NULL,
    fecha TEXT NOT NULL,
    hora TEXT NOT NULL
);

CREATE TABLE notas_de_prestamos (
    codigo INTEGER NOT NULL,
    anotacion TEXT NOT NULL,
    fecha TEXT NOT NULL
);
