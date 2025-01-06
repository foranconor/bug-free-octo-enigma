# TODO 

Basic algorithm

1. Run lua script, it will return a table that is the configuration for a stair
2. Check that all terminal elements are strictly numbers or strings
3. Run validation for the kind of stair it's claiming to be (strategy pattern)
4. Separate sections out and turn into generic structs
5. Flatten remaining tables
6. Output as stair struct

## Parsing and verification functions

We need to check that all sub tables are ok.
- rebates
- nosing
- materials (risers, treads, stringers)
- details

### Kinds which may differ (more than one may apply)

- Housed
- Cut
- Exterior
- Decking
- Sloping

###  
