from nicegui import ui

import src.functions.insertar_boletas as fib


class Boletas:
    def __init__(self):
        self.values = []

    def append(self, entry):
        self.values.append(entry)

    def remove(self, entry):
        self.values.remove(entry)

    def get_data_to_send(self):
        try:
            datos = [
                {
                    "codigo": int(item["n1"].value), 
                    "boleta": "-".join(
                        [
                            str(abs(int(num))).zfill(3)
                            for num in [
                                item["n1"].value,
                                item["n2"].value,
                                item["n3"].value,
                                item["n4"].value,
                            ]
                        ]                  
                    )
                }
                for item in self.values
            ]
        except:
            ui.notification("Todos los numeros de todas las boletas deben estar completos", type="warning")
        
        return datos


def render_insertar_boletas() -> None:
    boletas_disponibles, err = fib.obtener_boletas_disponibles()
    if err != "":
        ui.notify(err, type="warning")
        return

    if not boletas_disponibles:
        boletas_disponibles = []

    idx_boletas: list[int, ] = [boleta["codigo"] for boleta in boletas_disponibles]

    with ui.row().classes("w-full"):
        with ui.column().classes("w-1/2"):
            with ui.card().classes("w-full"):
                ui.label("Insertar boletas:").classes("text-3xl font-blod")

                with ui.row():
                    container = ui.column()
                    boletas = Boletas()

                    def add_entry_row():
                        with container, ui.row() as row:
                            numero_1 = ui.number().classes("w-1/6")
                            numero_2 = ui.number().classes("w-1/6")
                            numero_3 = ui.number().classes("w-1/6")
                            numero_4 = ui.number().classes("w-1/6")

                            # Guardamos el par y la fila (para poder borrarla luego)
                            entry = {
                                "row": row,
                                "n1": numero_1,
                                "n2": numero_2,
                                "n3": numero_3,
                                "n4": numero_4,
                            }
                            boletas.append(entry)

                            # Botón individual para eliminar esta fila específica
                            ui.button(
                                icon="delete",
                                on_click=lambda: remove_pair(entry),
                            ).props("flat color=red")

                    def remove_pair(entry):
                        container.remove(entry["row"])
                        boletas.remove(entry)
 
                    def ingresar_boleta() -> None:
                        boletas_a_enviar = boletas.get_data_to_send()
                        
                        if not boletas_a_enviar:
                            ui.notification("Se tiene que agregar al menos una boleta", type="warning")
                            return

                        fib.enviar_boletas_ingresadas(boletas_a_enviar)

                        ui.navigate.to("/insertarboletas")

                with ui.row():
                    ui.button("Ingresar", on_click=ingresar_boleta, color="green")
                    ui.button(icon="add", on_click=add_entry_row, color="purple")


            with ui.card().classes("w-full"):
                ui.label("Eliminar boletas:").classes("text-2xl font-blod")
                boleta_a_eliminar = ui.select(idx_boletas, label="Boleta a eliminar:").classes("w-full")

                def eliminar_boleta() -> None:
                    if boleta_a_eliminar.value is None:
                        ui.notification("Seleccione una boleta", type="warning")
                        return

                    fib.eliminar_boleta(boleta_a_eliminar.value)
                    ui.navigate.to("/insertarboletas")

                ui.button("Eliminar", on_click= eliminar_boleta, color="red")

        with ui.column():
            if boletas_disponibles:
                ui.table(rows=boletas_disponibles).classes("w-full")
            else:
                ui.label("No hay boletas disponibles").classes("text-lg text-grey")