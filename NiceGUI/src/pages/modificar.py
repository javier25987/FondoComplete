from nicegui import ui

import src.functions.modificar as fm
import src.functions.prestamos as fp


def render_modificar():
    users, err = fm.obtener_usuarios()
    if err != "":
        ui.notify(err, type="negative")
        return

    with ui.row().classes("w-full"):
        with ui.column().classes("flex-1"):
            with ui.card().classes("w-full"):
                ui.label("Ingresar Usuarios:").classes("text-5xl font-blod")
                nombre_user = ui.input("Nombre:").classes("w-full")
                telefono = ui.input("Telefono:").classes("w-full")
                puestos = ui.number("Puestos", value=1, min=1, step=1, precision=0).classes(
                    "w-full"
                )

                dialog = ui.dialog().props("backdrop-filter='blur(8px) brightness(40%)'")

                def enviar_solicitud():
                    dialog.close()

                    err = fm.guardar_usuario(
                        {
                            "nombre": nombre_user.value,
                            "telefono": telefono.value,
                            "puestos": int(puestos.value),
                        }
                    )

                    if err != "":
                        ui.notify(err, type="negative")
                    else:
                        ui.navigate.to("/modificar")

                def cancelar_solicitud():
                    dialog.close()
                    ui.notify("Solicitud cancelada", type="warning")

                def preparar_dialog():
                    dialog.clear()

                    with dialog, ui.column():
                        ui.label("Confirmar Usuario:").classes(
                            "text-5xl text-white font-bold"
                        )

                        ui.separator().classes("my-4")

                        ui.label(f"Nombre: {nombre_user.value}").classes(
                            "text-2xl text-white"
                        )
                        ui.label(f"Telefono: {telefono.value}").classes(
                            "text-2xl text-white"
                        )
                        ui.label(f"Puestos: {int(puestos.value)}").classes(
                            "text-2xl text-white"
                        )

                        ui.separator().classes("my-4")

                        with ui.row().classes("w-full justify-end my-4"):
                            ui.button(
                                "Cancelar", on_click=cancelar_solicitud, color="red"
                            ).props("flat")
                            ui.button(
                                "Confirmar y Enviar",
                                color="green",
                                on_click=enviar_solicitud,
                            )

                    dialog.open()

                ui.button("Agregar Usuario", on_click=preparar_dialog)

            with ui.card().classes("w-full"):
                ui.label("Modificar fechas de un prestamo:").classes("text-3xl font-blod")

                prestamo_a_modificar = ui.number("Prestamo a modificar:", min=1).classes("w-full")
                nueva_fecha = ui.date_input("Nueva fecha inicial:").classes("w-full")

                def modificar_fecha_inicial():
                    if prestamo_a_modificar.value is None or nueva_fecha.value is None:
                        ui.notify("Todos los valores deben estar completos", type="warning")
                        return
                    
                    err = fm.modificar_fechas_de_un_prestamo(
                        {
                            "codigo": int(prestamo_a_modificar.value),
                            "fecha": nueva_fecha.value,
                        }
                    )

                    if err != "":
                        ui.notify(err, type="negative")
                    else:
                        ui.navigate.to("/modificar")
                    
                ui.button("Modificar", on_click=modificar_fecha_inicial)

            with ui.card().classes("w-full"):
                ui.label("Abonar al prestamo:").classes("text-3xl font-blod")

                prestamo_a_abonarle = ui.number("Prestamo a abonarle:", value=1, min=0).classes("w-full")

                tipo_de_abono = ui.radio(["DEUDA", "INTERESES"], value="DEUDA")

                monto = ui.number("Monto a abonar:", value=1, min=0).classes("w-full")

                propietario_del_prestamo = ui.number("Propietario del prestamo", value=1, min=0).classes("w-full")

                fecha_del_pago = ui.date_input("Fecha del abono").classes("w-full")

                dialog = ui.dialog().props(
                    "backdrop-filter='blur(8px) brightness(40%)'"
                )

                def abonar_a_prestamo_con_fecha_modificada():
                    dialog.close()

                    err = fp.abonar_a_prestamo(
                        int(propietario_del_prestamo.value),
                        {
                            "monto": int(monto.value),
                            "codigo": int(prestamo_a_abonarle.value),
                            "tipo": tipo_de_abono.value,
                            "fecha": fecha_del_pago.value,
                        },
                    )

                    if err != "":
                        ui.notify(err, type="negative")
                    else:
                        ui.navigate.to("/modificar")
                        ui.notify(
                            "Abono relizado correctamente", type="positive"
                        )

                def cancelar_solicitud():
                    ui.notify("Solicitud cancelada", type="warning")
                    dialog.close()

                def preparar_dialog():
                    dialog.clear()

                    with dialog, ui.column():
                        ui.label("Confirmar Abono:").classes(
                            "text-4xl font-bold text-white"
                        )

                        ui.label(
                            f"Abonar a prestamo: {int(prestamo_a_abonarle.value)}"
                        ).classes("text-xl font-bold text-white")
                        ui.label(
                            f"Monto a abonar: {int(monto.value):,}"
                        ).classes("text-xl font-bold text-white")
                        ui.label(
                            f"Abono a: {tipo_de_abono.value}"
                        ).classes("text-xl font-bold text-white")
                        ui.label(
                            f"Nueva fecha: {fecha_del_pago.value}"
                        ).classes("text-xl font-bold text-white")

                        with ui.row().classes("w-full justify-end mt-4"):
                            ui.button(
                                "Cancelar",
                                on_click=cancelar_solicitud,
                                color="red",
                            ).props("flat")
                            ui.button(
                                "Confirmar y Enviar",
                                color="green",
                                on_click=abonar_a_prestamo_con_fecha_modificada,
                            )

                    dialog.open()

                ui.button(
                    "Iniciar Abono", on_click=preparar_dialog, color="green"
                ).classes("w-full")

        with ui.column().classes("flex-1"):
            if users:
                ui.table(rows=users).classes("w-full")
            else:
                ui.label("No hay usuarios en el sistema").classes("text-lg text-grey")
