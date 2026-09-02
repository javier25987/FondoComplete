--- 1. TABLAS INDEPENDIENTES ---

-- Tablas del sistema --

CREATE TABLE ajustes (
    ajuste TEXT UNIQUE NOT NULL,
    valor_n INTEGER NOT NULL,
    valor_t TEXT NOT NULL
);

CREATE TABLE registros (
    fecha DATE NOT NULL,
    ingreso INTEGER NOT NULL,
    egreso INTEGER NOT NULL,
    mes INTEGER NOT NULL
);

CREATE TABLE datos_de_rifas (
    codigo SERIAL PRIMARY KEY,
    numero_de_boletas INTEGER NOT NULL,
    costo_de_boletas INTEGER NOT NULL,
    costos_de_administracion INTEGER NOT NULL,
    ganancias_por_boleta INTEGER NOT NULL
);

-- Tablas de usuarios --

CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    nombre TEXT NOT NULL,
    puestos INTEGER NOT NULL,
    telefono TEXT NOT NULL,
    estado BOOLEAN NOT NULL DEFAULT TRUE,
    cuotas_pagas INTEGER NOT NULL DEFAULT 0,
    cuotas_deudas INTEGER NOT NULL DEFAULT 0,
    semanas_revisadas INTEGER NOT NULL DEFAULT 0,
    deuda_rifas INTEGER NOT NULL DEFAULT 0,
    multas_semanales INTEGER NOT NULL DEFAULT 0,
    multas_extras INTEGER NOT NULL DEFAULT 0,
    multas_por_cabeza INTEGER NOT NULL DEFAULT 0,
    capital_pago INTEGER NOT NULL DEFAULT 0,
    capital_retirado INTEGER NOT NULL DEFAULT 0,
    capital_proyectado INTEGER NOT NULL DEFAULT 0
);

--- 2. TABLAS DEPENDIENTES (Con Llaves Foráneas) ---

CREATE TABLE anotaciones (
    codigo SERIAL PRIMARY KEY,
    id_user INTEGER NOT NULL,
    anotacion TEXT NOT NULL,
    monto INTEGER NOT NULL,
    tipo TEXT NOT NULL,
    fecha DATE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_usuario
    	FOREIGN KEY (id_user)
    	REFERENCES usuarios(id) -- ON DELETE CASCADE
);

CREATE TABLE prestamos_hechos (
    codigo SERIAL PRIMARY KEY,
    id_user INTEGER NOT NULL,
    estado BOOLEAN NOT NULL DEFAULT TRUE,
    interes INTEGER NOT NULL,
    intereses_vencidos INTEGER NOT NULL DEFAULT 0,
    intereses_generados INTEGER NOT NULL DEFAULT 0,
    deuda INTEGER NOT NULL,
    monto INTEGER NOT NULL,
    fechas_revisadas INTEGER NOT NULL DEFAULT 0,
    fecha DATE DEFAULT CURRENT_TIMESTAMP,
    motivo TEXT NOT NULL,

    CONSTRAINT fk_usuario_prestamo
    	FOREIGN KEY (id_user)
    	REFERENCES usuarios(id)
);

CREATE TABLE deudas (
    codigo SERIAL PRIMARY KEY,
    id_user INTEGER NOT NULL,
    monto INTEGER NOT NULL,
    fecha_de_creacion TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    fecha_de_pago TIMESTAMP(0) NOT NULL,
    tipo TEXT NOT NULL,
    motivo TEXT NOT NULL,
    estado BOOLEAN NOT NULL DEFAULT TRUE,

    CONSTRAINT fk_usuario_deuda
    	FOREIGN KEY (id_user)
    	REFERENCES usuarios(id)
);

CREATE TABLE transferencias (
    codigo SERIAL PRIMARY KEY,
    id_user INTEGER NOT NULL,
    fecha TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    monto INTEGER NOT NULL,

    CONSTRAINT fk_usuario_transf
    	FOREIGN KEY (id_user)
    	REFERENCES usuarios(id)
);

CREATE TABLE cheques_de_cuotas (
    codigo SERIAL PRIMARY KEY,
    id_user INTEGER NOT NULL,
    cuotas_pagas INTEGER NOT NULL,
    multas_pagas INTEGER NOT NULL,
    metodo_de_pago TEXT NOT NULL,
    fecha TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_usuario_cheque
    	FOREIGN KEY (id_user)
    	REFERENCES usuarios(id)
);

CREATE TABLE cheques_de_prestamos (
    codigo SERIAL PRIMARY KEY,
    id_user INTEGER NOT NULL,
    codigo_prestamo INTEGER NOT NULL,
    monto_pago INTEGER NOT NULL,
    deuda_antes INTEGER NOT NULL,
    deuda_despues INTEGER NOT NULL,
    fecha TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_usuario_pago
    	FOREIGN KEY (id_user)
    	REFERENCES usuarios(id),

    CONSTRAINT fk_prestamo_pago
    	FOREIGN KEY (codigo_prestamo)
    	REFERENCES prestamos_hechos(codigo)
);

CREATE TABLE notas_de_prestamos (
    codigo_prestamo INTEGER NOT NULL,
    anotacion TEXT NOT NULL,
    fecha DATE DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_nota_prestamo
    	FOREIGN KEY (codigo_prestamo)
    	REFERENCES prestamos_hechos(codigo)
);

CREATE TABLE boletas_rifas (
    codigo_rifa INTEGER NOT NULL,
    codigo_boleta INTEGER NOT NULL,
    boleta TEXT NOT NULL,
    dada_a INTEGER,

    CONSTRAINT fk_usuario_rifas
    	FOREIGN KEY (dada_a)
    	REFERENCES usuarios(id),

    CONSTRAINT fk_codigo_rifa
        FOREIGN KEY (codigo_rifa)
        REFERENCES datos_de_rifas(codigo),

    CONSTRAINT unique_codigo_rifa_boleta
        UNIQUE(codigo_rifa, boleta)
);

--- 3. TABLAS PARA NORMALIZAR ---

CREATE TABLE cuotas_bloqueos (
    id_user INTEGER NOT NULL,
    semana INTEGER NOT NULL,

    CONSTRAINT fk_bloqueos
    	FOREIGN KEY (id_user)
    	REFERENCES usuarios(id),

    CONSTRAINT unique_id_user_semana
        UNIQUE(id_user, semana)
);

CREATE TABLE cuotas_multas_activas (
    id_user INTEGER NOT NULL,
    semana INTEGER NOT NULL,
    valor INTEGER NOT NULL,

    CONSTRAINT fk_multas_activas
    	FOREIGN KEY (id_user)
    	REFERENCES usuarios(id),

    CONSTRAINT unique_iduser_activas
        UNIQUE(id_user, semana)
);

CREATE TABLE cuotas_multas_pagas (
    id_user INTEGER NOT NULL,
    semana INTEGER NOT NULL,
    valor INTEGER NOT NULL,

    CONSTRAINT fk_multas_pagas
    	FOREIGN KEY (id_user)
    	REFERENCES usuarios(id),

    CONSTRAINT unique_iduser_pagas
        UNIQUE(id_user, semana)
);

CREATE TABLE premios_de_rifas (
    codigo_rifa INTEGER NOT NULL,
    posicion INTEGER NOT NULL,
    premio INTEGER NOT NULL,

    CONSTRAINT fk_premios_de_rifas
    	FOREIGN KEY (codigo_rifa)
    	REFERENCES datos_de_rifas(codigo),

    CONSTRAINT unique_codigo_rifa_posicion
        UNIQUE(codigo_rifa, posicion)
);

CREATE TABLE prestamos_fiadores (
    id_user_fiador INTEGER NOT NULL,
    codigo_prestamo INTEGER NOT NULL,
    monto_respaldado INTEGER NOT NULL,

    CONSTRAINT fk_usuario_fiador
        FOREIGN KEY (id_user_fiador)
        REFERENCES usuarios(id),

    CONSTRAINT fk_en_prestamo
        FOREIGN KEY (codigo_prestamo)
        REFERENCES prestamos_hechos(codigo)
);

CREATE TABLE prestamos_fechas_de_pago (
    codigo_prestamo INTEGER NOT NULL,
    fecha DATE NOT NULL,

    CONSTRAINT fk_fecha_de_pago
        FOREIGN KEY (codigo_prestamo)
        REFERENCES prestamos_hechos(codigo)
);

--- 4. TABLA PARA SUPER USUARIO (1976) ---

CREATE TABLE super_usuario (
    codigo_prestamo INTEGER UNIQUE NOT NULL,
    monto_respaldado INTEGER NOT NULL,

    CONSTRAINT fk_en_prestamo
        FOREIGN KEY (codigo_prestamo)
        REFERENCES prestamos_hechos(codigo)
);
