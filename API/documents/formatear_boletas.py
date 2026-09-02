import sqlite3


def main() -> None:
    conexion = sqlite3.connect("FondoTest.db")
    cursor = conexion.cursor()

    cursor.execute(
        """
            SELECT 
                br.codigo_rifa,
                br.codigo_boleta,
                br.boleta 
            FROM boletas_rifas br;
        """
    )

    datos = cursor.fetchall()

    for rifa, codigo_boleta, boleta in datos:
        cursor.execute(
            """
                UPDATE boletas_rifas
                SET boleta = ?
                WHERE codigo_rifa = ? AND codigo_boleta = ?; 
            """,
            (
                "-".join([numero.zfill(3) for numero in boleta.split("-")]),
                rifa,
                codigo_boleta,
            ),
        )

    conexion.commit()

    print("Todas las boletas ya fueron formateadas.")


if __name__ == "__main__":
    main()
