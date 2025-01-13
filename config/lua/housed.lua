local materials = require("lua/materials")
local cfg = require("lua/stair")

cfg.style = "Housed"

cfg.risers = materials.riser_triboard
cfg.treads = materials.tread_triboard
cfg.stringers = materials.radiata250x40

-- This is the width of the stringer that protrudes above the pitch line
cfg.skirting = 40
-- "Horns" are the extra material at the bottom and top of housed stringers
cfg.bottom_horn_length = 100
cfg.top_horn_height = 80
cfg.top_horn_length = 80

return cfg
