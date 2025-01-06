local materials = require("lua/materials")
local cfg = require("lua/essential")

cfg.style = "Cut"

cfg.risers = materials.thin_triboard
cfg.treads = materials.thick_triboard
cfg.stringers = materials.radiata300x40

-- clear stringer rebates, not strictly needed as they'll be ignored
cfg.rebates.stringer = nil
cfg.rebates.wedge_angle = nil

return cfg
