local untreated = "UT"
local triboard = "Triboard"
local radiata = "Radiata Pine"
local h3 = "H3.2"

return {
	-- sheet materials
	riser_triboard = {
		name = triboard,
		treatment = untreated,
		dims = {
			x = 2400,
			y = 1200,
			z = 15,
		},
	},
	tread_triboard = {
		name = triboard,
		treatment = untreated,
		dims = {
			x = 2400,
			y = 1200,
			z = 21,
		},
	},
	-- stringer materials
	radiata250x40 = {
		name = radiata,
		treatment = untreated,
		dims = {
			x = 4800,
			y = 235,
			z = 33,
		},
	},
	radiata300x50 = {
		name = radiata,
		treatment = untreated,
		dims = {
			x = 4800,
			y = 280,
			z = 44,
		},
	},
	exterior = {
		name = radiata,
		treatment = h3,
		dims = {
			x = 4800,
			y = 280,
			z = 44,
		},
	},
	-- decking materials
	kwila = {
		name = "Kwila",
		treatment = untreated,
		dims = {
			x = 4000,
			y = 90,
			z = 19,
		},
	},
}
