package analizar

const FORMATO_BOLETAS_GENERALES = `
	#let datos = json(%q)

	#set text(size: 13pt)
	#set page(margin: (x: 2cm, y: 2cm))

	#align(center)[
		= Boletas en la rifa #datos.at("rifa_actual")
	]

	#for usuario in datos.at("usuarios") {
		[== #usuario.at("puesto") - #usuario.at("nombre")]
		for entrega in usuario.at("entregas") [
			/ #entrega.at("fecha"): #entrega.at("boletas").join(",    ")
		]
	}

	#if datos.at("boletas_libres").len() != 0 [
		== Boletas disponibles
		#datos.at("boletas_libres").join(",  ")
	]
`
