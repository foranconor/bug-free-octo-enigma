function Fact(num)
	if num == 0 then
		return 1
	else
		return num * Fact(num - 1)
	end
end

Lib = require("test2")

return Lib.neg(Fact(6))
