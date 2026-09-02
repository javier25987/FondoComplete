# Ends points

2/34

## Generales

### cuotas
+ [] `/cuotas/{id}` GET -> da los datos del usuario para la secion cuotas
+ [] `/cuotas/{id}` PUT `json: {"cuotas": int >= 0, "multas": int >= 0, "metodo": "EFECTIVO" | "TRANSFERENCIA"}` -> realiza el pago de las cuotas y multas del usuario `json: {"codigo_cheque": int}`
+ [ ] `/cuotas/bloc/{id}` PUT `json: {semana: int}` -> bloquea o desbloquea la semana
+ [] `/cuotas/cheque/{codigo}` POST -> crea un cheque de cuotas con los datos indicados

### prestamos
+ [] `/prestamos/{id}` GET -> da los datos del usuario para la secion prestamos
+ [] `/prestamos/{id}` PUT `json: {"monto": int > 0, "codigo": int > 0}` -> realiza el pago para el prestamo en cuestion
+ [] `/prestamos/{id}` POST `json: {"monto": int, "motivo": string, "fiadores": []{"id_user":int, "deuda":int} }` -> crea un prestamo para la persona en cuestion
+ [] `/prestamos/cheque/{codigo}` POST -> crea un cheque de prestamo con los datos indicados

### analisis
+ [ ] `/analizar` GET -> devuelve el analisis general del programa
+ [ ] `/analizar/{id}` GET -> devuelve el analisis del usuario en cuestion

### transferencias 
+ [ ] `/transferencias/{id}` GET -> retorna todas las transferencias de la persona en cuestion si el id es cero retorna todas las transferencias

### rifas
+ [ ] `/rifas/{id}` GET -> retorna todos los datos de las rifas
+ [ ] `/rifas/{id}` POST `json: {boletas: [int, ...]}` -> da todas las boletas especifcadas al usuario
+ [ ] `/rifas/{id}` PUT `json: {monto: int}` -> paga la deuda del usuario

### deudas
+ [] `/deudas/{id}` GET -> retorna todas las deudas del usuario en cuestion
+ [] `/deudas/{id}?codigo=int` PUT -> paga la deuda del usuario
+ [] `/deudas/{id}` POST `json: {motivo: int, monto: int, tipo: "MULTA POR CABEZA" | "MULTA EXTRA"}` -> crea una deuda para el usuario

<!-- ### anotaciones
+ [X] `/anotaciones/{id}` GET -> retora todos los datos de usuario
+ [X] `/anotaciones/{id}` POST `json: {nota: string, monto: int, tipo: string}` -> crea una anotacion para el usuario -->

### ver usuarios
+ [] `/ver/usuarios?nombre=str` GET -> retorna todos los usuarios que coincidan
+ [] `/ver/acuerdo` GET -> retorna todos los usuarios que necesiten hacer acuerdo
+ [] `/ver/boletas/{rifa}?boleta=str` GET -> retorna todo lo necesario para identificar una boleta en una rifa
+ [] `/ver/prestamo/{codigo}` GET -> retorna toda la informacion de un prestamo

<!-- ### registros
+ [X] `/registros/{mes(int)}` GET -> retorna todos los registros del mes en cuestion -->

## Administrativas

### insertar boletas
+ [] `/insertar` GET -> da todas las boletas que se pueden eliminar
+ [] `/insertar` POST `json: [{codigo: int, boleta: str}, ...]` -> inserta la boleta especificada en la base de datos
+ [] `/insertar?codigo=int` DELETE  -> eliminina la boleta especificada, si esta no existe retorna un error

### modificar socios
+ [X] `/modificar` GET -> retorna todos los usuarios del programa
+ [X] `/modificar` POST `json: {nombre: string, telefono string, puestos: int}` -> crea un usuario nuevo
+ [ ] `/modificar/fechas` PUT `json: {codigo: int, fecha: string}` -> modifica la fecha de un prestamos
+ [ ] `/modificar/semana` PUT `json: {semana: int, tipo: string, valor: int}` -> modifica un valor de una semana especificada

### ajustes
+ [] `/ajuste?ajuste=str` GET -> retorna un ajuste en concreto, el solicitado
+ [] `/ajuste` PUT `json: {ajuste: string, nvalor: int, tvalor: string}` -> modifica un ajuste
+ [] `/ajuste` POST `json: {ajuste: string, nvalor: int, tvalor: string}` -> crea un ajuste
+ [] `/ajustes` GET -> retorna todos los ajustes del programa
+ [] `/ajustes/boletas` PUT `json: {numero_de_boletas: int, premios: [int, ], costo_de_boletas: int, costos_de_administracion: int}` -> guarda todo los valores para una rifa
+ [ ] `/ajustes/boletas` POST -> Cierra la rifa y carga todas las deudas de boletas como prestamos


> crear una fucion de logs la cual registre si se creo una multa y los ineteses de un prestamo