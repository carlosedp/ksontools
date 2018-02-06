local k = import "k.libsonnet";
local params = {
  option: "value",
  int: 9,
  obj: {
    a: 1,
    b: "c",
    d: {
      nested: "deep",
    },
  },
  arr: [1,2,3],
};
local a = "a";
local b = "b";
{
}