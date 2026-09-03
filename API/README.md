# Documentacion de la api

## cosas a tener en cuenta

- todas las fechas de la api van a estar en formato `Año-Mes-Dia Hora:Minutos`

## notas importantes

> recrear la tabla de los cheques de rifas con la llave de las rifas creadas

## Ends points

28/37

### Generales

#### cuotas
+ [X] `/cuotas/{id}` GET -> da los datos del usuario para la secion cuotas
+ [X] `/cuotas/{id}` PUT `json: {"cuotas": int >= 0, "multas": int >= 0, "metodo": "EFECTIVO" | "TRANSFERENCIA"}` -> realiza el pago de las cuotas y multas del usuario `json: {"codigo_cheque": int}`
+ [ ] `/cuotas/bloc/{id}` PUT `json: {semana: int}` -> bloquea o desbloquea la semana
+ [X] `/cuotas/cheque/{codigo}` GET -> retorna el pdf del cheque solicitado (si no existe el pdf lo crea y lo retorna)

#### prestamos
+ [X] `/prestamos/{id}` GET -> da los datos del usuario para la secion prestamos
+ [X] `/prestamos/{id}` PUT `json: {"monto": int > 0, "tipo": string ["DEUDA", "INTERESES"], "codigo": int > 0}` -> realiza el pago del monto del prestamo en cuestion
+ [X] `/prestamos/{id}` POST `json: {"monto": int, "motivo": string, "fiadores": []{"id_user":int, "deuda":int} }` -> crea un prestamo para la persona en cuestion
+ [X] `/prestamos/cheque/{codigo}` GET -> retorna el pdf del cheque solicitado (si no existe el pdf lo crea y lo retorna)

#### analisis
+ [ ] `/analizar` GET -> devuelve el analisis general del programa
+ [ ] `/analizar/{id}` GET -> devuelve el analisis del usuario en cuestion
+ [ ] `/analizar/boletas` GET -> devuelve el analisis general de las boletas

#### transferencias 
+ [ ] `/transferencias` GET -> retorna todas las transferencias de todos los usuarios
+ [ ] `/transferencias/{id}` GET -> retorna todas las transferencias de la persona en cuestion

#### rifas
+ [X] `/rifas/{id}` GET -> retorna todos los datos de las rifas
+ [X] `/rifas/{id}` POST `json: {boletas: [int, ...]}` -> da todas las boletas especifcadas al usuario
+ [X] `/rifas/{id}` PUT `json: {monto: int}` -> paga la deuda del usuario
+ [ ] `/rifas/cheque/{id}` GET -> retorna el cheque pedido

#### deudas
+ [X] `/deudas/{id}` GET -> retorna todas las deudas del usuario en cuestion
+ [X] `/deudas/{id}?codigo=int` PUT -> paga la deuda del usuario
+ [X] `/deudas/{id}` POST `json: {motivo: int, monto: int, tipo: "MULTA POR CABEZA" | "MULTA EXTRA"}` -> crea una deuda para el usuario

<!-- ### anotaciones
+ [X] `/anotaciones/{id}` GET -> retora todos los datos de usuario
+ [X] `/anotaciones/{id}` POST `json: {nota: string, monto: int, tipo: string}` -> crea una anotacion para el usuario -->

#### ver usuarios
+ [X] `/ver/usuarios?nombre=str` GET -> retorna todos los usuarios que coincidan
+ [X] `/ver/acuerdo` GET -> retorna todos los usuarios que necesiten hacer acuerdo
+ [X] `/ver/boletas/{rifa}?boleta=str` GET -> retorna todo lo necesario para identificar una boleta en una rifa
+ [X] `/ver/prestamo/{codigo}` GET -> retorna toda la informacion de un prestamo

<!-- ### registros
+ [X] `/registros/{mes(int)}` GET -> retorna todos los registros del mes en cuestion -->

### Administrativas

#### insertar boletas
+ [X] `/insertar` GET -> da todas las boletas que se pueden eliminar
+ [X] `/insertar` POST `json: [{codigo: int, boleta: str}, ...]` -> inserta la boleta especificada en la base de datos
+ [X] `/insertar?codigo=int` DELETE  -> eliminina la boleta especificada, si esta no existe retorna un error

#### modificar socios
+ [X] `/modificar` GET -> retorna todos los usuarios del programa
+ [X] `/modificar` POST `json: {nombre: string, telefono string, puestos: int}` -> crea un usuario nuevo
+ [X] `/modificar/fechas` PUT `json: {codigo: int, fecha: string}` -> modifica la fecha de un prestamos
+ [ ] `/modificar/semana` PUT `json: {semana: int, tipo: string, valor: int}` -> modifica un valor de una semana especificada

#### ajustes

<!-- + [ ] `/ajuste?ajuste=str` GET -> retorna un ajuste en concreto, el solicitado
+ [ ] `/ajuste` PUT `json: {ajuste: string, nvalor: int, tvalor: string}` -> modifica un ajuste
+ [ ] `/ajuste` POST `json: {ajuste: string, nvalor: int, tvalor: string}` -> crea un ajuste -->

+ [X] `/ajuste/calendario` GET -> retorna el calendario del programa
+ [X] `/ajuste/calendario` POST `json: {fecha_inicial: string, fechas_dobles: []string, hora_de_cierre: string}` -> crea un calendario, si habia uno antes lo borra y reescribe 
+ [X] `/ajustes` GET -> retorna todos los ajustes del programa
+ [X] `/ajustes/boletas` GET -> retorna todas las rifas creadas
+ [X] `/ajustes/boletas` PUT `json: {numero_de_boletas: int, premios: [int, ], costo_de_boletas: int, costos_de_administracion: int}` -> guarda todo los valores para una rifa
+ [ ] `/ajustes/boletas` POST -> Cierra la rifa y carga todas las deudas de boletas como prestamos


> crear una fucion de logs la cual registre si se creo una multa y los ineteses de un prestamo

## Ideas del proyecto

> Esta api se va a desarroyar como un backend for frontend, es decir no va a haber multiples endpoints los cuales van a tener multiples usos y conultas, todo se va a hacer pensando en que el backend esta hecho para el frontend

