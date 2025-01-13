local stair1 = require("lua/housed")
local stair2 = require("lua/housed")

local job = require("lua/job")

job.quote = "Q1234"
job.address = "40 Raiha St"
job.site_contact = {
	name = "Hone Smith",
	phone = "021 123 4567",
}

stair1.name = "1"
stair1.classification = "Secondary Private"

stair1.sections = {
	{
		kind = "Straight",
		steps = 10,
		width = 1000,
		split = {
			where = "middle",
			how = "short",
		},
	},
	{
		kind = "Winder",
		steps = 3,
		start_width = 1000,
		end_width = 967,
		direction = "Left",
	},
}

stair2.name = "Upper"
stair2.quantity = 4

stair2.sections = {
	{
		kind = "Straight",
		steps = 5,
		width = 1023,
		split = nil,
	},
}

job.stairs = {
	stair1,
	stair2,
}

return job
