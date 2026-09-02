import src.functions.ajustes as fa

from nicegui import ui


def render_ajutes():
    ajustes, err = fa.obtener_todos_los_ajustes()
    if err != "":
        ui.notify(err, type="negative")
        return

    with ui.tabs().classes("w-full") as tabs:
        calendario_tab = ui.tab("Calendario")
        ajustes_tab = ui.tab("Ajustes")
        rifas_tab = ui.tab("Rifas")
    with ui.tab_panels(tabs, value=calendario_tab).classes("w-full"):
        with ui.tab_panel(calendario_tab):
            estado, calendario = fa.obtener_calendario()

            if estado:
                if len(calendario) != 50:
                    with ui.card(align_items="center").classes("w-full bg-red-200"):
                        ui.label("🚨 No existe un calendario por favor creelo").classes(
                            "text-2xl"
                        )
                else:
                    with ui.row():
                        for i in range(5):
                            ui.table(rows=calendario[10 * i : 10 * (i + 1)])

                with ui.card().classes("w-full"):
                    ui.label("Crear calendario:").classes("text-2xl")

                    with ui.row():
                        primera_fecha = ui.date_input("Primera fecha:").classes("w-50")

                        fechas_dobles_colum = ui.column()
                        fechas_dobles = []

                        def add_fecha_doble():
                            with fechas_dobles_colum:
                                with ui.row() as row:
                                    fecha_doble = ui.date_input("Fecha doble").classes(
                                        "w-50"
                                    )

                                    entry = {"row": row, "fecha": fecha_doble}
                                    fechas_dobles.append(entry)

                                    ui.button(
                                        icon="delete",
                                        on_click=lambda: remove_pair(entry),
                                    ).props("flat color=red")

                        def remove_pair(entry):
                            fechas_dobles_colum.remove(entry["row"])
                            fechas_dobles.remove(entry)

                        ui.button(icon="add", on_click=add_fecha_doble, color="purple")

                        def crear_fechas():
                            fecha_to_func = primera_fecha.value
                            dobles_to_func = [
                                val["fecha"].value
                                for val in fechas_dobles
                                if val["fecha"].value != ""
                            ]

                            err = fa.crear_listado_de_fechas(
                                fecha_to_func, dobles_to_func
                            )

                            if err != "":
                                ui.notify(err, type="negative")
                            else:
                                ui.navigate.to("/ajustes")

                        ui.button("Crear", on_click=crear_fechas)

                    # ui.label(
                    #     "Todas las fechas deben estar en el formato 'año-mes-dia' use el " \
                    #     "selector de fechas para evitar errores NO ESCRIBA USTED MISMO " \
                    #     "LAS FECHAS."
                    # ).classes("text-xl text-grey")
            else:
                ui.notify(calendario, type="negative")

        with ui.tab_panel(rifas_tab):
            with ui.grid(columns=2).classes("w-full"):
                with ui.card(): #.classes("h-130"):
                    ui.label("Crear una Rifa:").classes("text-2xl font-bold")

                    cantidad_de_boletas = ui.number("Cantidad de boletas:").classes("w-full")

                    with ui.row():
                        premio_1 = ui.number("Primer premio:").classes("w-2/5")
                        premio_2 = ui.number("Segundo premio:").classes("w-2/5")
                        premio_3 = ui.number("Tercer premio:").classes("w-2/5")
                        premio_4 = ui.number("Cuarto premio:").classes("w-2/5")

                    costo_de_boleta = ui.number("Costo de una boleta:").classes("w-full")

                    costos_de_administracion = ui.number("Costos de administracion:").classes("w-full")

                    def crear_rifa() -> None:
                        datos = {
                            "numero_de_boletas": int(cantidad_de_boletas.value), 
                            "premios": [
                                {
                                    "numero": 1,
                                    "premio": int(premio_1.value),
                                },
                                {
                                    "numero": 2,
                                    "premio": int(premio_2.value),
                                },
                                {
                                    "numero": 3,
                                    "premio": int(premio_3.value),
                                },
                                {
                                    "numero": 4,
                                    "premio": int(premio_4.value),
                                },
                            ], 
                            "costo_de_boletas": int(costo_de_boleta.value), 
                            "costos_de_administracion": int(costos_de_administracion.value),
                        }

                        err = fa.enviar_rifa_al_servidor(datos)

                        if err != "":
                            ui.notification(err, type="negative")
                            return

                        ui.navigate.to("/ajustes")

                        

                    ui.button("Crear", on_click=crear_rifa)

                rifas_creadas, err = fa.obtener_rifas_creadas()
                if err != "":
                    ui.notification(err, type="negative")
                    return

                if rifas_creadas:
                    for rifa in rifas_creadas:
                        with ui.card():
                            ui.label(f"Rifa №: {rifa["codigo"]}").classes("text-2xl font-bold")
                            ui.label(f"Boletas en la rifa: {rifa["numero_de_boletas"]:,}").classes("text-1xl")
                            ui.label(f"Costo de la boleta: {rifa["costo_de_boletas"]:,}").classes("text-1xl")
                            ui.label(f"Costos de administracion: {rifa["costos_de_administracion"]:,}").classes("text-1xl")
                            ui.label(f"Ganancia por boleta: {rifa["ganancia_por_boleta"]:,}").classes("text-1xl")
                            ui.label("Premios:").classes("text-1xl font-bold")
                            ui.table(rows=rifa["premios"])





        # with ui.tab_panel(ajustes_tab):
        #     with ui.row().classes("w-full"):
        #         with ui.column().classes("flex-1"):
        #             ui.label("Ajustes por defecto:").classes("text-3xl")
        #             ui.table(rows=fa.AJUSTES_POR_DEFECTO)
        #             ui.button("Crear Ajustes", on_click=fa.crear_ajustes_por_defecto)

        #             ui.separator()
        #             ui.button("Restaurar Ajustes", on_click=fa.restaurar_ajustes_por_defecto)

        #         #ui.space()

        #         with ui.column().classes("flex-1"):
        #             ui.label("Ajustes actuales:").classes("text-3xl")
        #             ui.table(rows=ajustes["ajustes"])

        #             with ui.card():
        #                 ui.label(
        #                     "El ajuste calendario es ignorado por defecto para " \
        #                     "esta tabla pero lo puede consultar y modificar en " \
        #                     "la seccion 'calendario'"
        #                 ).classes("text-xl text-grey")

        #             ui.separator()

        #             ui.label("Crear ajustes:").classes("text-2xl")

        #             ajuste_crear = ui.input("Nombre del ajuste:").classes("w-80")
        #             nvalor_crear = ui.number("Valor numerico del ajuste:").classes("w-80")
        #             tvalor_crear = ui.input("Valor textual del ajuste:").classes("w-80")

        #             def crear_ajuste():
        #                 err = fa.crear_ajuste(
        #                     ajuste_crear.value, nvalor_crear.value, tvalor_crear.value
        #                 )

        #                 if err != "":
        #                     ui.notify(err, type="negative")
        #                 else:
        #                     ui.notify("Ajuste creado correctamente", type="positive")

        #                 ui.navigate.to("/ajustes")

        #             ui.button("Crear Ajuste", on_click=crear_ajuste)

        #             ui.separator()

        #             ui.label("Modificar ajuste").classes("text-2xl")

        #             ajuste_modif = ui.select(
        #                 ajustes["nombres_ajustes"],
        #                 label="Nombre del ajuste:"
        #             ).classes("w-80")
        #             nvalor_modif = ui.number("Valor numerico del ajuste:").classes("w-80")
        #             tvalor_modif = ui.input("Valor textual del ajuste:").classes("w-80")

        #             def modificar_ajuste():
        #                 err = fa.modificar_ajuste(
        #                     ajuste_modif.value, nvalor_modif.value, tvalor_modif.value #type:ignore
        #                 )

        #                 if err != "":
        #                     ui.notify(err, type="negative")
        #                 else:
        #                     ui.notify("Ajuste modificado correctamente", type="positive")

        #                 ui.navigate.to("/ajustes")

        #             ui.button("Modificar Ajuste", on_click=modificar_ajuste)

        # with ui.tab_panel(rifas_tab):
        #     pass
