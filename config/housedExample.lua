local cfg = require("lua/housed")
local materials = require("lua/materials")

cfg.details = {
	quote = "K471",
	stairway = "1",
	quantity = 3,
	address = "40 Raiha St",
	classification = "Main Private",
	site_contact = {
		name = "John Smith",
		phone = "021 123 4567",
	},
}

cfg.sections = {
	{
		kind = "Straight",
		steps = 10,
		width = 1000,
	},
	{
		kind = "Winder",
		steps = 3,
		start_width = 1000,
		end_width = 967,
		direction = "Left",
	},
}

-- override housed defaults
cfg.going = 280
cfg.nosing.overhang = 0
cfg.nosing.bottom_radius = 0
cfg.stringers = materials.radiata300x50

return cfg
