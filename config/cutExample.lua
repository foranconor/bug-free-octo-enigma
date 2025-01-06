local cfg = require("lua/cut")

cfg.details = {
	quote = "T843",
	stairway = "Upper",
	address = "40 Raiha St",
	site_contact = {
		name = "John Smith",
		phone = "021 123 4567",
	},
}

cfg.sections = {
	{
		kind = "Straight",
		steps = 15,
		width = 998,
	},
}

cfg.testing = true
cfg.other = function(v)
	return v + 1
end

return cfg
